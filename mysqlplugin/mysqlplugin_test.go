// +build unit

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
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MySQLPluginSuite struct {
	suite.Suite
}

func (suite *MySQLPluginSuite) TearDownSuite() {}

func (suite *MySQLPluginSuite) SetupSuite() {}

func (suite *MySQLPluginSuite) TestGetMetricTypes() {
	Convey("Given MySQL plugin initialized", mps.T(), func() {
		mySQLPlugin := New()

		Convey("When one wants to get list of available meterics", func() {
			mts, err := mySQLPlugin.GetMetricTypes(plugin.PluginConfigType{})

			Convey("Then error should not be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then list of metrics is returned", func() {
				So(len(mts), ShouldEqual, 12)

				namespaces := []string{}
				for _, m := range mts {
					namespaces = append(namespaces, strings.Join(m.Namespace(), "/"))
				}

				So(namespaces, ShouldContain, "intel")

			})
		})
	})
}

func (suite *MySQLPluginSuite) TestCollectMetrics() {
	Convey("Given MySQL plugin initlialized", suite.T(), func() {
		mySQLPlugin := New()

		Convey("When one wants to get values for given metric types", func() {
			mTypes := []plugin.PluginMetricType{
				plugin.PluginMetricType{Namespace_: []string{"intel"}},
			}

			metrics, err := mySQLPlugin.CollectMetrics(mTypes)

			Convey("Then no erros should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then proper metrics values are returned", func() {
				So(len(metrics), ShouldEqual, 4)

				stats := map[string]uint64{}
				for _, m := range metrics {
					n := strings.Join(m.Namespace(), "/")
					v, ok := m.Data().(uint64)
					if ok {
						stats[n] = v
					}
				}

				assert.Equal(suite.T(), len(metrics), len(stats))

				So(stats["intel"], ShouldEqual, suite.cache*1024)
			})

		})
	})
}

func TestMySQLPluginSuite(t *testing.T) {
	suite.Run(t, &MySQLPluginSuite{})
}