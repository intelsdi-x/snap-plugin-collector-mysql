/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysqlplugin

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

	"github.com/intelsdi-x/snap-plugin-utilities/config"

	"github.com/intelsdi-x/snap-plugin-collector-mysql/stats"
	"github.com/intelsdi-x/snap/core"
)

const (
	// Name of plugin
	Name = "mysql"
	// Version of plugin
	Version = 3
	// Type of plugin
	Type = plugin.CollectorPluginType
)

// MySQLPlugin is implementation of plugin.Plugin interface.
type MySQLPlugin struct {
	initialized      bool
	initializedMutex *sync.Mutex

	callDiscovery map[string]int

	mysql collector
}

// New returns initialized instance of MySQL Plugin collector
func New() *MySQLPlugin {
	self := new(MySQLPlugin)
	self.initializedMutex = new(sync.Mutex)
	self.callDiscovery = map[string]int{}

	return self
}

// CollectMetrics finds required request ids required to collect given metrics,
// asks collector service for metrics associated with these calls and returns
// requested metrics. Error is returned when metric collection failed or plugin
// initialization was unsuccessful.
func (p *MySQLPlugin) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {

	if len(mts) > 0 {
		err := p.init(mts[0])

		if err != nil {
			return nil, err
		}

	} else {
		return mts, nil
	}

	t := time.Now()

	results := make([]plugin.MetricType, len(mts))

	calls := map[int]bool{}

	for _, mt := range mts {
		name := parseName(mt.Namespace().Strings())
		calls[p.callDiscovery[name]] = true
	}

	metrics, err := p.mysql.Collect(calls)

	if err != nil {
		return nil, err
	}

	for i, mt := range mts {
		results[i] = plugin.MetricType{
			Namespace_: mt.Namespace(),
			Data_:      metrics[parseName(mt.Namespace().Strings())],
			Timestamp_: t,
		}
	}

	return results, nil
}

// GetMetricTypes returns list of available metrics. If initialization failed
// error is returned.
func (p *MySQLPlugin) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	err := p.init(cfg)

	if err != nil {
		return nil, err
	}

	mts := []plugin.MetricType{}

	for k := range p.callDiscovery {
		mts = append(mts, plugin.MetricType{Namespace_: core.NewNamespace(makeName(k)...)})
	}

	return mts, nil
}

// GetConfigPolicy returns plugin config policy
func (p *MySQLPlugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// Meta returns plugin's metadata
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(Name, Version, Type, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// init performs one time initialization of plugin. Reads configuration from cfg
// and constructs all service objets that will be used during plugin's lifetime.
// returns error if initialization failed.
func (p *MySQLPlugin) init(cfg interface{}) error {
	p.initializedMutex.Lock()
	defer p.initializedMutex.Unlock()

	if p.initialized {
		return nil
	}

	cfgItems, err := config.GetConfigItems(cfg, "mysql_connection_string", "mysql_use_innodb")

	if err != nil {
		return fmt.Errorf("plugin initalization failed : [%v]", err)
	}

	sqlStats, err := makeStats(cfgItems["mysql_connection_string"].(string))

	if err != nil {
		return err
	}

	p.mysql = makeCollector(sqlStats, cfgItems["mysql_use_innodb"].(bool))

	metrics, err := p.mysql.Discover()
	if err != nil {

		// for easier mocking
		if sqlStats != nil {
			sqlStats.Close()
		}

		return err
	}

	for _, m := range metrics {
		p.callDiscovery[m.Name] = m.Call
	}

	p.initialized = true
	return nil
}

// for mocking
var makeStats = func(connectionString string) (mysqlSource, error) { return stats.New(connectionString) }
var makeCollector = func(statsSource mysqlSource, useInnodb bool) collector { return NewCollector(statsSource, useInnodb) }

// prefix of all namespaces
var namespacePrefix = []string{"intel", "mysql"}

type collector interface {
	Discover() ([]metric, error)
	Collect(metrics map[int]bool) (map[string]interface{}, error)
}

// makeName makes namespace from metric path (with segments separated by '/' ) and
// namespace prefix
func makeName(m string) []string {
	res := []string{}
	res = append(res, namespacePrefix...)
	res = append(res, strings.Split(m, "/")...)

	return res
}

// parseName extracts metric path from namespace by trimming prefix and concatenating
// remaining segments with '/'.
func parseName(ns []string) string {
	return strings.Join(ns[len(namespacePrefix):], "/")
}
