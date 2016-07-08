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
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/intelsdi-x/snap-plugin-collector-mysql/stats"
)

type statsMock struct {
	mock.Mock
}

func (self *statsMock) GetStatus(parseInnodb bool) (stats.Stats, error) {
	args := self.Mock.Called(parseInnodb)

	r0 := *args.Get(0).(*interface{})

	if r0 == nil {
		return nil, errors.New("x")
	}

	err, isErr := (r0).(error)

	if isErr {
		return nil, err
	}

	return r0.(stats.Stats), args.Error(1)
}

func (self *statsMock) GetInnodb() (stats.Stats, error) {
	args := self.Mock.Called()

	r0 := *args.Get(0).(*interface{})

	if r0 == nil {
		return nil, errors.New("x")
	}

	err, isErr := (r0).(error)

	if isErr {
		return nil, err
	}

	return r0.(stats.Stats), args.Error(1)
}
func (self *statsMock) GetMasterStatus() (stats.Stats, error) {
	args := self.Mock.Called()

	r0 := *args.Get(0).(*interface{})

	if r0 == nil {
		return nil, errors.New("x")
	}

	err, isErr := (r0).(error)

	if isErr {
		return nil, err
	}

	return r0.(stats.Stats), args.Error(1)
}
func (self *statsMock) GetSlaveStatus() (stats.Stats, error) {
	args := self.Mock.Called()

	r0 := *args.Get(0).(*interface{})

	if r0 == nil {
		return nil, errors.New("x")
	}

	err, isErr := (r0).(error)

	if isErr {
		return nil, err
	}

	return r0.(stats.Stats), args.Error(1)
}
func (self *statsMock) Close() error {
	args := self.Mock.Called()
	return args.Error(0)
}

func TestDiscover(t *testing.T) {

	Convey("Discover", t, func() {

		mockedOrg := newMockedStats()
		mocked := mockedOrg

		source := statsMock{}

		source.On("GetStatus", mock.Anything).Return(mocked.statusPtr, nil)
		source.On("GetInnodb").Return(mocked.innodbPtr, nil)
		source.On("GetMasterStatus").Return(mocked.masterPtr, nil)
		source.On("GetSlaveStatus").Return(mocked.slavePtr, nil)

		sut := NewCollector(&source, true)

		dut, dut_err := sut.Discover()

		Convey("requests status data", func() {

			source.AssertCalled(t, "GetStatus", mock.Anything)

			Convey("with correct innodb flag", func() {

				source.AssertCalled(t, "GetStatus", true)

				sut2 := NewCollector(&source, false)
				sut2.Discover()
				source.AssertCalled(t, "GetStatus", false)

			})

			Convey("returns error when request fails", func() {

				sut2 := NewCollector(&source, false)

				*mocked.statusPtr = nil

				_, dut_err2 := sut2.Discover()

				So(dut_err2, ShouldNotBeNil)

			})

			Convey("if request succeed", func() {

				Convey("exposes data", func() {

					content := map[metric]bool{}

					for _, v := range dut {
						content[v] = true
					}

					So(content[metric{Name: "global/stat1", Call: callGlobal}], ShouldBeTrue)
					So(content[metric{Name: "global/stat2", Call: callGlobal}], ShouldBeTrue)

				})

			})

		})

		Convey("if innodb is enabled", func() {

			Convey("requests innodb data", func() {

				source.AssertCalled(t, "GetInnodb")

				Convey("returns error when request fails", func() {

					sut2 := NewCollector(&source, true)

					*mocked.innodbPtr = nil

					_, dut_err2 := sut2.Discover()

					So(dut_err2, ShouldNotBeNil)

				})

				Convey("if request succeed", func() {

					Convey("exposes data", func() {

						content := map[metric]bool{}

						for _, v := range dut {
							content[v] = true
						}

						So(content[metric{Name: "inno/stat1", Call: callInnoDB}], ShouldBeTrue)
						So(content[metric{Name: "inno/stat2", Call: callInnoDB}], ShouldBeTrue)
					})

				})

			})

		})

		Convey("if innodb is disabled", func() {

			source := statsMock{}

			source.On("GetStatus", mock.Anything).Return(mocked.statusPtr, nil)
			source.On("GetInnodb").Return(mocked.innodbPtr, nil)
			source.On("GetMasterStatus").Return(mocked.masterPtr, nil)
			source.On("GetSlaveStatus").Return(mocked.slavePtr, nil)

			sut := NewCollector(&source, false)

			sut.Discover()

			Convey("does not request innodb data", func() {

				source.AssertNotCalled(t, "GetInnodb")
			})

		})

		Convey("tries to request master data", func() {

			source.AssertCalled(t, "GetMasterStatus")

			Convey("does not fail when master data is unavailable", func() {

				sut2 := NewCollector(&source, false)

				*mocked.slavePtr = nil

				_, dut_err2 := sut2.Discover()

				So(dut_err2, ShouldBeNil)
			})

			Convey("if master data is available", func() {

				Convey("exposes data", func() {

					content := map[metric]bool{}

					for _, v := range dut {
						content[v] = true
					}

					So(content[metric{Name: "master/stat1", Call: callMaster}], ShouldBeTrue)
					So(content[metric{Name: "master/stat2", Call: callMaster}], ShouldBeTrue)
				})

			})

		})

		Convey("tries to request slave data", func() {

			source.AssertCalled(t, "GetSlaveStatus")

			Convey("does not fail when slave data is unavailable", func() {

				sut2 := NewCollector(&source, false)

				*mocked.slavePtr = nil

				_, dut_err2 := sut2.Discover()

				So(dut_err2, ShouldBeNil)
			})

			Convey("if slave data is available", func() {

				Convey("exposes data", func() {

					content := map[metric]bool{}

					for _, v := range dut {
						content[v] = true
					}

					So(content[metric{Name: "slave/stat1", Call: callSlave}], ShouldBeTrue)
					So(content[metric{Name: "slave/stat2", Call: callSlave}], ShouldBeTrue)

				})

			})

		})

		Convey("returns no error if all requests succeed", func() {

			So(dut_err, ShouldBeNil)
		})

	})

}

func TestCollect(t *testing.T) {

	Convey("Collect", t, func() {

		mockedOrg := newMockedStats()
		mocked := mockedOrg

		source := statsMock{}

		source.On("GetStatus", mock.Anything).Return(mocked.statusPtr, nil)
		source.On("GetInnodb").Return(mocked.innodbPtr, nil)
		source.On("GetMasterStatus").Return(mocked.masterPtr, nil)
		source.On("GetSlaveStatus").Return(mocked.slavePtr, nil)

		sut := NewCollector(&source, true)

		Convey("Should do each requested call", func() {

			Convey("Global", func() {

				sut.Collect(map[int]bool{callGlobal: true})
				source.AssertCalled(t, "GetStatus", true)

			})

			Convey("InnoDB", func() {

				sut.Collect(map[int]bool{callInnoDB: true})
				source.AssertCalled(t, "GetInnodb")

			})

			Convey("MasterStatus", func() {

				sut.Collect(map[int]bool{callMaster: true})
				source.AssertCalled(t, "GetMasterStatus")

			})

			Convey("SlaveStatus", func() {

				sut.Collect(map[int]bool{callSlave: true})
				source.AssertCalled(t, "GetSlaveStatus")

			})

		})

		Convey("Doesn't do unnecessary calls", func() {

			Convey("Global", func() {

				sut.Collect(map[int]bool{callInnoDB: true, callMaster: true, callSlave: true})
				source.AssertNotCalled(t, "GetStatus")

			})

			Convey("InnoDB", func() {

				sut.Collect(map[int]bool{callGlobal: true, callMaster: true, callSlave: true})
				source.AssertNotCalled(t, "GetInnodb")

			})

			Convey("MasterStatus", func() {

				sut.Collect(map[int]bool{callGlobal: true, callInnoDB: true, callSlave: true})
				source.AssertNotCalled(t, "GetMasterStatus")

			})

			Convey("SlaveStatus", func() {

				sut.Collect(map[int]bool{callGlobal: true, callMaster: true, callInnoDB: true})
				source.AssertNotCalled(t, "GetSlaveStatus")

			})

		})

		Convey("After subsequent calls", func() {

			orgTimeNow := timeNow

			Reset(func() {
				timeNow = orgTimeNow

			})

			timeNow = func() time.Time { return time.Unix(100, 0) }

			Convey("Gauges are exposed as raw value", func() {

				(*mocked.statusPtr).(stats.Stats)["global/stat0"] = stats.Stat{Value: 10, Type: stats.Gauge, IsNull: false}
				dut1, _ := sut.Collect(map[int]bool{callGlobal: true})

				So(dut1["global/stat0"], ShouldAlmostEqual, 10, 0.1)
				_, ok := dut1["global/stat0"].(int64)
				So(ok, ShouldBeTrue)

				timeNow = func() time.Time { return time.Unix(102, 0) }

				(*mocked.statusPtr).(stats.Stats)["global/stat0"] = stats.Stat{Value: 20, Type: stats.Gauge, IsNull: false}
				dut2, _ := sut.Collect(map[int]bool{callGlobal: true})

				So(dut2["global/stat0"], ShouldAlmostEqual, 20, 0.1)
				_, ok = dut2["global/stat0"].(int64)
				So(ok, ShouldBeTrue)
			})

			Convey("Derives are exposed as ratio of change to time", func() {

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 10, Type: stats.Derive, IsNull: false}
				dut1, _ := sut.Collect(map[int]bool{callGlobal: true})

				// derive's rate should be nil on the first measurement
				So(dut1["global/stat1"], ShouldBeNil)

				timeNow = func() time.Time { return time.Unix(102, 0) }

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 20, Type: stats.Derive, IsNull: false}
				dut2, _ := sut.Collect(map[int]bool{callGlobal: true})

				So(dut2["global/stat1"], ShouldAlmostEqual, 5, 0.1)
				_, ok := dut2["global/stat1"].(float64)
				So(ok, ShouldBeTrue)

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 20, Type: stats.Derive, IsNull: true}
				dut3, _ := sut.Collect(map[int]bool{callGlobal: true})
				// derive's rate should be nil when the current value is also nil
				So(dut3["global/stat1"], ShouldBeNil)
			})

			Convey("Counters are exposed as ratio of change to time", func() {

				(*mocked.statusPtr).(stats.Stats)["global/stat2"] = stats.Stat{Value: 10, Type: stats.Counter, IsNull: false}
				dut1, _ := sut.Collect(map[int]bool{callGlobal: true})
				// counter's rate should be nil on the first measurement
				So(dut1["global/stat2"], ShouldBeNil)

				timeNow = func() time.Time { return time.Unix(102, 0) }

				(*mocked.statusPtr).(stats.Stats)["global/stat2"] = stats.Stat{Value: 20, Type: stats.Counter, IsNull: false}
				dut2, _ := sut.Collect(map[int]bool{callGlobal: true})

				So(dut2["global/stat2"], ShouldAlmostEqual, 5, 0.1)
				_, ok := dut2["global/stat2"].(float64)
				So(ok, ShouldBeTrue)

				(*mocked.statusPtr).(stats.Stats)["global/stat2"] = stats.Stat{Value: 20, Type: stats.Derive, IsNull: true}
				dut3, _ := sut.Collect(map[int]bool{callGlobal: true})
				// counter's rate should be nil when the current value is also nil
				So(dut3["global/stat2"], ShouldBeNil)

			})

		})

	})

}

func mockStat(ns string) *interface{} {
	s := stats.Stats{}

	s[ns+"/stat1"] = stats.Stat{Value: 1, Type: stats.Gauge, IsNull: false}
	s[ns+"/stat2"] = stats.Stat{Value: 2, Type: stats.Derive, IsNull: false}

	var ifc interface{} = s

	return &ifc
}

type statMockData struct {
	statusPtr, innodbPtr, masterPtr, slavePtr *interface{}
}

func newMockedStats() statMockData {
	self := statMockData{}
	self.statusPtr = mockStat("global")
	self.innodbPtr = mockStat("inno")
	self.masterPtr = mockStat("master")
	self.slavePtr = mockStat("slave")

	return self
}
