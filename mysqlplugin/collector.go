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
	"math"
	"os"
	"time"

	"github.com/intelsdi-x/snap-plugin-collector-mysql/stats"
)

const (
	callGlobal = iota
	callInnoDB
	callMaster
	callSlave
)

var width32bit = math.Pow(2, 32.0)
var width64bit = math.Pow(2, 64.0)

// NewCollector constructs new metricCollector that will query given statsSource.
// useInnodb indicates if innodb statistics are gathered (and gathering will
// fail if they are unavailable).
func NewCollector(statsSource mysqlSource, useInnodb bool) *metricCollector {
	self := new(metricCollector)
	self.counters = map[string]metricValue{}
	self.UseInnodb = useInnodb
	self.StatsSource = statsSource
	return self
}

// Collect performs given set of calls (indicated by true value in metrics map).
// returns map of metric values (accessible by metric name). If any of requested
// calls fail error is returned.
func (mc *metricCollector) Collect(metrics map[int]bool) (map[string]interface{}, error) {

	res := map[string]interface{}{}

	if metrics[callGlobal] {
		st, err := mc.StatsSource.GetStatus(mc.UseInnodb)
		if err != nil {
			return nil, err
		}

		mc.updateStats(res, st)
	}

	if metrics[callInnoDB] {
		st, err := mc.StatsSource.GetInnodb()
		if err != nil {
			return nil, err
		}

		mc.updateStats(res, st)
	}

	if metrics[callMaster] {
		st, err := mc.StatsSource.GetMasterStatus()
		if err != nil {
			return nil, err
		}

		mc.updateStats(res, st)
	}

	if metrics[callSlave] {
		st, err := mc.StatsSource.GetSlaveStatus()
		if err != nil {
			return nil, err
		}

		mc.updateStats(res, st)
	}

	return res, nil
}

// Discover performs metric discovery. Returns valid metric names and associated
// Call id's. If mandatory request fails error is returned. No error is returned
// when master or slave stats can't be read because server may not be configured
// to work in master-slave mode.
func (mc *metricCollector) Discover() ([]metric, error) {
	res := []metric{}

	st, err := mc.StatsSource.GetStatus(mc.UseInnodb)
	if err != nil {
		return nil, err
	}
	addMetrics(&res, st, callGlobal)

	if mc.UseInnodb {
		st, err = mc.StatsSource.GetInnodb()
		if err != nil {
			return nil, err
		}
		addMetrics(&res, st, callInnoDB)
	}

	// server may not have master or slave stats

	st, err = mc.StatsSource.GetMasterStatus()
	if err == nil {
		addMetrics(&res, st, callMaster)
	}

	st, err = mc.StatsSource.GetSlaveStatus()
	if err == nil {
		addMetrics(&res, st, callSlave)
	}

	return res, nil

}

// for unit testing
var timeNow = func() time.Time { return time.Now() }

type mysqlSource interface {
	GetStatus(parseInnodb bool) (stats.Stats, error)
	GetInnodb() (stats.Stats, error)
	GetMasterStatus() (stats.Stats, error)
	GetSlaveStatus() (stats.Stats, error)
	Close() error
}

// metric contains name of metric and id of call that collects particular
// metric.
type metric struct {
	Name string
	Call int
}

// metricValue holds value of metric and time of last collection.
type metricValue struct {
	Value          int64
	CollectionTime time.Time
}

// metricCollector implements logic for discovering available metrics
// and associated queries, performing given set of queries and performing
// rate calculation.
type metricCollector struct {
	StatsSource mysqlSource
	UseInnodb   bool

	counters map[string]metricValue
}

// addMetrics appends metric names from st to dst array setting Call
// field to given value.
func addMetrics(dst *[]metric, st stats.Stats, call int) {
	for k := range st {
		*dst = append(*dst, metric{Name: k, Call: call})
	}
}

// helper func that converts Stat to nullable value.
// Returns Stat.Value or nil.
func val(s stats.Stat) interface{} {
	if s.IsNull {
		return nil
	}
	return s.Value
}

// updateStats adds metrics from st to res. While gauges are copied as they are, values for
// counters and derives are differentiated and represents rate of change in time.
// and for them send a null on the first measurement (or if the last time was null too)
func (mc *metricCollector) updateStats(res map[string]interface{}, st stats.Stats) {

	for k, v := range st {
		switch v.Type {
		case stats.Gauge:
			res[k] = val(v)

		case stats.Derive, stats.Counter:
			if v.IsNull {
				res[k] = nil
				delete(mc.counters, k)
				continue
			}

			mv := metricValue{Value: v.Value, CollectionTime: timeNow()}

			old, ok := mc.counters[k]

			if !ok {
				// for metrics representing a rate of change
				// send a null on the first measurement
				res[k] = nil
				mc.counters[k] = mv
				continue
			}

			delta := float64(mv.Value - old.Value)

			if v.Type == stats.Counter {
				// for counters the behaviour differs when value_new < value_old
				// and wrapping around should be taken into account

				if delta < 0 {
					//set wrapper width to 64 bits width as default
					wrapperWidth := width64bit
					if float64(old.Value) < width32bit {
						wrapperWidth = width32bit
					}

					delta = wrapperWidth + delta
				}

			}

			res[k] = delta / mv.CollectionTime.Sub(old.CollectionTime).Seconds()

			mc.counters[k] = mv

		default:
			fmt.Fprintln(os.Stderr, "Metric `", k, "` cannot be classified as a one of the supported data type (gauge, derive or counter)")
		}
	}
}
