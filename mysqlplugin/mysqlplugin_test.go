// +build medium

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
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/intelsdi-x/snap-plugin-collector-mysql/stats"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
)

type nullSqlsource struct {
}

func (self *nullSqlsource) GetStatus(parseInnodb bool) (stats.Stats, error) {
	return nil, nil
}
func (self *nullSqlsource) GetInnodb() (stats.Stats, error) {
	return nil, nil
}
func (self *nullSqlsource) GetMasterStatus() (stats.Stats, error) {
	return nil, nil
}
func (self *nullSqlsource) GetSlaveStatus() (stats.Stats, error) {
	return nil, nil
}
func (self *nullSqlsource) Close() error {
	return nil
}

type collectorMock struct {
	mock.Mock
}

func (self *collectorMock) Discover() ([]metric, error) {
	args := self.Called()
	var r0 []metric = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).([]metric)
	}
	return r0, args.Error(1)
}

func (self *collectorMock) Collect(metrics map[int]bool) (map[string]interface{}, error) {
	args := self.Called(metrics)
	var r0 map[string]interface{} = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).(map[string]interface{})
	}
	return r0, args.Error(1)
}

func testingConfig() (cfg1 plugin.ConfigType, cfg2 *cdata.ConfigDataNode) {
	cfg1 = plugin.NewPluginConfigType()
	cfg2 = cdata.NewNode()
	cfg1.AddItem("mysql_connection_string", ctypes.ConfigValueStr{Value: "root:r00tme@tcp(localhost:3306)/"})
	cfg2.AddItem("mysql_connection_string", ctypes.ConfigValueStr{Value: "root:r00tme@tcp(localhost:3306)/"})

	cfg1.AddItem("mysql_use_innodb", ctypes.ConfigValueBool{Value: true})
	cfg2.AddItem("mysql_use_innodb", ctypes.ConfigValueBool{Value: true})

	return
}

func TestGetMetricTypes(t *testing.T) {
	Convey("GetMetricTypes", t, func() {

		orgMakeStats := makeStats
		orgMakeCollector := makeCollector

		Reset(func() {

			makeCollector = orgMakeCollector
			makeStats = orgMakeStats

		})

		mock := &collectorMock{}

		makeStats = func(connectionString string) (mysqlSource, error) { return new(nullSqlsource), nil }
		makeCollector = func(statsSource mysqlSource, useInnodb bool) collector { return mock }

		cfg1, _ := testingConfig()

		sut := New()

		Convey("if initalization succeeds", func() {

			mock.On("Discover").Return([]metric{metric{Name: "aaa/bbb", Call: 1}, metric{Name: "x/y/z", Call: 2}}, nil)

			dut, dut_err := sut.GetMetricTypes(cfg1)

			Convey("returns complete list of metrics", func() {

				content := map[string]bool{}

				for _, v := range dut {
					content[v.Namespace().String()] = true
				}

				So(content["/intel/mysql/aaa/bbb"], ShouldBeTrue)
				So(content["/intel/mysql/x/y/z"], ShouldBeTrue)

			})

			Convey("and no error", func() {

				So(dut_err, ShouldBeNil)

			})

		})

		Convey("if initialization fails", func() {

			Convey("on stats construction", func() {

				makeStats = func(connectionString string) (mysqlSource, error) { return nil, errors.New("x") }

				_, dut_err := sut.GetMetricTypes(cfg1)

				Convey("error is returned", func() {

					So(dut_err, ShouldNotBeNil)

				})

			})

			Convey("on discovery", func() {

				mock.On("Discover").Return(nil, errors.New("x"))

				_, dut_err := sut.GetMetricTypes(cfg1)

				Convey("error is returned", func() {
					So(dut_err, ShouldNotBeNil)
				})

			})

		})

	})
}

func TestCollectMetrics(t *testing.T) {
	Convey("CollectMetrics", t, func() {

		orgMakeStats := makeStats
		orgMakeCollector := makeCollector

		Reset(func() {

			makeCollector = orgMakeCollector
			makeStats = orgMakeStats

		})

		mocked := &collectorMock{}

		makeStats = func(connectionString string) (mysqlSource, error) { return new(nullSqlsource), nil }
		makeCollector = func(statsSource mysqlSource, useInnodb bool) collector { return mocked }

		_, cfg2 := testingConfig()

		sut := New()

		mts10 := make([]plugin.MetricType, 10)
		metrics10 := make([]metric, 10)

		for i, _ := range mts10 {
			mts10[i] = plugin.MetricType{Namespace_: core.NewNamespace("intel", "mysql", fmt.Sprintf("stat%d", i)), Config_: cfg2}
			metrics10[i] = metric{Name: fmt.Sprintf("stat%d", i), Call: i}

		}

		Convey("performs init even if GetMetricTypes was not called", func() {

			mts := []plugin.MetricType{
				plugin.MetricType{Namespace_: core.NewNamespace("intel", "mysql", "aaa", "bbb"), Config_: cfg2},
			}

			mocked.On("Discover").Return([]metric{metric{Name: "aaa/bbb", Call: 1}, metric{Name: "x/y/z", Call: 2}}, nil)
			mocked.On("Collect", mock.Anything).Return(nil, errors.New("x"))

			sut.CollectMetrics(mts)

			mocked.AssertCalled(t, "Discover")

		})

		Convey("requests all required calls", func() {

			var dut interface{}
			mocked.On("Discover").Return(metrics10, nil)
			mocked.On("Collect", mock.Anything).Return(nil, errors.New("x")).Run(func(args mock.Arguments) {
				dut = args.Get(0)
			})

			sut.CollectMetrics(mts10)

			for i, _ := range mts10 {
				So(dut.(map[int]bool)[i], ShouldBeTrue)
			}

		})

		Convey("omits uneccessary calls", func() {

			for i, _ := range mts10 {
				newMts := []plugin.MetricType{}
				newMts = append(newMts, mts10[0:i]...)
				newMts = append(newMts, mts10[i+1:]...)

				var dut interface{}
				*mocked = collectorMock{}
				mocked.On("Discover").Return(metrics10, nil)
				mocked.On("Collect", mock.Anything).Return(nil, errors.New("x")).Run(func(args mock.Arguments) {
					dut = args.Get(0)
				})

				sut.CollectMetrics(newMts)
				mocked.AssertCalled(t, "Collect", mock.Anything)
				So(dut.(map[int]bool)[i], ShouldBeFalse)

			}

		})

		Convey("exposes correct data", func() {

			mocked.On("Discover").Return(metrics10, nil)
			result := map[string]interface{}{}

			for _, v := range metrics10 {
				result[v.Name] = 100 + v.Call
			}

			mocked.On("Collect", mock.Anything).Return(result, nil)

			dut, _ := sut.CollectMetrics(mts10)

			vals := map[string]interface{}{}
			for _, v := range dut {
				vals[v.Namespace().String()] = v.Data()
			}

			for _, v := range metrics10 {
				So(vals["/intel/mysql/"+v.Name], ShouldEqual, 100+v.Call)
				result[v.Name] = 100 + v.Call
			}

		})

		Convey("returns error if collection failed", func() {

			mocked.On("Discover").Return(metrics10, nil)
			mocked.On("Collect", mock.Anything).Return(nil, errors.New("x"))

			_, dut_err := sut.CollectMetrics(mts10)

			So(dut_err, ShouldNotBeNil)

		})

	})
}

func TestGetConfigPolicy(t *testing.T) {
	Convey("GetConfigPolicy", t, func() {
		sut := New()
		dut, dutErr := sut.GetConfigPolicy()
		Convey("Returns correct type", func() {
			So(dut, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
		})
		Convey("Is not nil", func() {
			So(dut, ShouldNotBeNil)
		})

		Convey("Returns no error", func() {
			So(dutErr, ShouldBeNil)
		})
	})
}
