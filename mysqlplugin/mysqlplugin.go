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
	"os"
	"strings"
	"sync"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

	"github.com/intelsdi-x/snap-plugin-utilities/config"

	"github.com/intelsdi-x/snap-plugin-collector-mysql/stats"
)

const (
	// Name of plugin
	Name = "mysql"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

// prefix of all namespaces
var namespacePrefix = []string{"intel", "mysql"}

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

type collector interface {
	Discover() ([]metric, error)
	Collect(metrics map[int]bool) (map[string]interface{}, error)
}

// MySQLPlugin is implementation of plugin.Plugin interface.
type MySQLPlugin struct {
	initialized      bool
	initializedMutex *sync.Mutex

	callDiscovery map[string]int

	mysql collector
}

// for mocking
var makeStats = func(connectionString string) (mysqlSource, error) { return stats.New(connectionString) }
var makeCollector = func(statsSource mysqlSource, useInnodb bool) collector { return NewCollector(statsSource, useInnodb) }

// init performs one time initialization of plugin. Reads configuration from cfg
// and constructs all service objets that will be used during plugin's lifetime.
// returns error if initialization failed.
func (p *MySQLPlugin) init(cfg interface{}) error {
	p.initializedMutex.Lock()
	defer p.initializedMutex.Unlock()

	if p.initialized {
		return nil
	}

	cfgItems, err := config.GetConfigItems(cfg, []string{"mysql_connection_string", "mysql_use_innodb"})

	if err != nil {
		panic(fmt.Errorf("plugin initalization failed : [%v]", err))
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

// CollectMetrics finds required request ids required to collect given metrics,
// asks collector service for metrics associated with these calls and returns
// requested metrics. Error is returned when metric collection failed or plugin
// initialization was unsuccessful.
func (p *MySQLPlugin) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {

	if len(mts) > 0 {
		err := p.init(mts[0])

		if err != nil {
			return nil, err
		}

	} else {
		return mts, nil
	}

	// it's not worth to abort collection
	// when only os.Hostname() raised error
	host, _ := os.Hostname()
	t := time.Now()

	results := make([]plugin.PluginMetricType, len(mts))

	calls := map[int]bool{}

	for _, mt := range mts {
		name := parseName(mt.Namespace())
		calls[p.callDiscovery[name]] = true
	}

	metrics, err := p.mysql.Collect(calls)

	if err != nil {
		return nil, err
	}

	for i, mt := range mts {
		results[i] = plugin.PluginMetricType{
			Namespace_: mt.Namespace(),
			Data_:      metrics[parseName(mt.Namespace())],
			Source_:    host,
			Timestamp_: t,
		}
	}

	return results, nil
}

// GetMetricTypes returns list of available metrics. If initialization failed
// error is returned.
func (p *MySQLPlugin) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	err := p.init(cfg)

	if err != nil {
		return nil, err
	}

	mts := []plugin.PluginMetricType{}

	for k, _ := range p.callDiscovery {
		mts = append(mts, plugin.PluginMetricType{Namespace_: makeName(k)})
	}

	return mts, nil
}

// GetConfigPolicy
func (p *MySQLPlugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

func New() *MySQLPlugin {
	self := new(MySQLPlugin)
	self.initializedMutex = new(sync.Mutex)
	self.callDiscovery = map[string]int{}

	return self
}

// Returns plugin's metadata
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(Name, Version, Type, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}
