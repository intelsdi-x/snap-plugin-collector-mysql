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
	"github.com/intelsdi-x/snap-plugin-collector-mysql/stats"
	"time"
)

const (
	CALL_GLOBAL = iota
	CALL_INNODB
	CALL_MASTER
	CALL_SLAVE
)

type mysqlSource interface {
	GetStatus(parseInnodb bool) (stats.Stats, error)
	GetInnodb() (stats.Stats, error)
	GetMasterStatus() (stats.Stats, error)
	GetSlaveStatus() (stats.Stats, error)
	Close() error
}

type metric struct {
	Name string
	Call int
}

type metricValue struct {
	Value          int64
	CollectionTime time.Time
}

type metricCollector struct {
	StatsSource mysqlSource
	UseInnodb   bool

	counters map[string]metricValue
}

func addMetrics(dst *[]metric, st stats.Stats, call int) {
	for k, _ := range st {
		*dst = append(*dst, metric{Name: k, Call: call})
	}
}

func (self *metricCollector) Discover() ([]metric, error) {
	res := []metric{}

	st, err := self.StatsSource.GetStatus(self.UseInnodb)
	if err != nil {
		return nil, err
	}
	addMetrics(&res, st, CALL_GLOBAL)

	if self.UseInnodb {
		st, err = self.StatsSource.GetInnodb()
		if err != nil {
			return nil, err
		}
		addMetrics(&res, st, CALL_INNODB)
	}

	// server may not have master or slave stats

	st, err = self.StatsSource.GetMasterStatus()
	if err == nil {
		addMetrics(&res, st, CALL_MASTER)
	}

	st, err = self.StatsSource.GetSlaveStatus()
	if err == nil {
		addMetrics(&res, st, CALL_SLAVE)
	}

	return res, nil

}

func val(s stats.Stat) interface{} {
	if s.IsNull {
		return nil
	} else {
		return s.Value
	}
}

// for unit testing
var timeNow = func() time.Time { return time.Now() }

func (self *metricCollector) updateStats(res map[string]interface{}, st stats.Stats) {
	t := timeNow()
	for k, v := range st {
		if v.Type == stats.TYPE_GAUGE {
			res[k] = val(v)
		} else {
			//DOES'NT MATTER IF IT'S COUNTER OR DERIVE

			mv := metricValue{Value: v.Value, CollectionTime: t}

			if v.IsNull {
				delete(self.counters, k)
			}

			old, ok := self.counters[k]

			if ok {
				delta := mv.Value - old.Value
				res[k] = float64(delta) / mv.CollectionTime.Sub(old.CollectionTime).Seconds()
				self.counters[k] = mv
			} else {
				res[k] = val(v)
				if !v.IsNull {
					self.counters[k] = mv
				}
			}

		}
	}
}

func (self *metricCollector) Collect(metrics map[int]bool) (map[string]interface{}, error) {

	res := map[string]interface{}{}

	if metrics[CALL_GLOBAL] {
		st, err := self.StatsSource.GetStatus(self.UseInnodb)
		if err != nil {
			return nil, err
		}

		self.updateStats(res, st)
	}

	if metrics[CALL_INNODB] {
		st, err := self.StatsSource.GetInnodb()
		if err != nil {
			return nil, err
		}

		self.updateStats(res, st)
	}

	if metrics[CALL_MASTER] {
		st, err := self.StatsSource.GetMasterStatus()
		if err != nil {
			return nil, err
		}

		self.updateStats(res, st)
	}

	if metrics[CALL_SLAVE] {
		st, err := self.StatsSource.GetSlaveStatus()
		if err != nil {
			return nil, err
		}

		self.updateStats(res, st)
	}

	return res, nil
}

func NewCollector(statsSource mysqlSource, useInnodb bool) *metricCollector {
	self := new(metricCollector)
	self.counters = map[string]metricValue{}
	self.UseInnodb = useInnodb
	self.StatsSource = statsSource
	return self
}
