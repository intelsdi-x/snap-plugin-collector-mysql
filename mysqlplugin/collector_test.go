package mysqlplugin

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/intelsdi-x/snap-plugin-collector-mysql/stats"
)

var smthErr = fmt.Errorf("smth")
var statNil = &stats.Stats{}

type statsMock struct {
	mock.Mock
}

func (self *statsMock) GetStatus(parseInnodb bool) (stats.Stats, error) {
	args := self.Mock.Called(parseInnodb)

	r0 := *args.Get(0).(*interface{})

	if r0 == nil {
		return nil, smthErr
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
		return nil, smthErr
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
		return nil, smthErr
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
		return nil, smthErr
	}

	err, isErr := (r0).(error)

	if isErr {
		return nil, err
	}

	return r0.(stats.Stats), args.Error(1)
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

					So(content[metric{Name: "global/stat1", Call: CALL_GLOBAL}], ShouldBeTrue)
					So(content[metric{Name: "global/stat2", Call: CALL_GLOBAL}], ShouldBeTrue)

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

						So(content[metric{Name: "inno/stat1", Call: CALL_INNODB}], ShouldBeTrue)
						So(content[metric{Name: "inno/stat2", Call: CALL_INNODB}], ShouldBeTrue)

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

					So(content[metric{Name: "master/stat1", Call: CALL_MASTER}], ShouldBeTrue)
					So(content[metric{Name: "master/stat2", Call: CALL_MASTER}], ShouldBeTrue)

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

					So(content[metric{Name: "slave/stat1", Call: CALL_SLAVE}], ShouldBeTrue)
					So(content[metric{Name: "slave/stat2", Call: CALL_SLAVE}], ShouldBeTrue)

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

				sut.Collect(map[int]bool{CALL_GLOBAL: true})
				source.AssertCalled(t, "GetStatus", true)

			})

			Convey("InnoDB", func() {

				sut.Collect(map[int]bool{CALL_INNODB: true})
				source.AssertCalled(t, "GetInnodb")

			})

			Convey("MasterStatus", func() {

				sut.Collect(map[int]bool{CALL_MASTER: true})
				source.AssertCalled(t, "GetMasterStatus")

			})

			Convey("SlaveStatus", func() {

				sut.Collect(map[int]bool{CALL_SLAVE: true})
				source.AssertCalled(t, "GetSlaveStatus")

			})

		})

		Convey("Doesn't do unnecessary calls", func() {

			Convey("Global", func() {

				sut.Collect(map[int]bool{CALL_INNODB: true, CALL_MASTER: true, CALL_SLAVE: true})
				source.AssertNotCalled(t, "GetStatus")

			})

			Convey("InnoDB", func() {

				sut.Collect(map[int]bool{CALL_GLOBAL: true, CALL_MASTER: true, CALL_SLAVE: true})
				source.AssertNotCalled(t, "GetInnodb")

			})

			Convey("MasterStatus", func() {

				sut.Collect(map[int]bool{CALL_GLOBAL: true, CALL_INNODB: true, CALL_SLAVE: true})
				source.AssertNotCalled(t, "GetMasterStatus")

			})

			Convey("SlaveStatus", func() {

				sut.Collect(map[int]bool{CALL_GLOBAL: true, CALL_MASTER: true, CALL_INNODB: true})
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

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 10, Type: stats.TYPE_GAUGE, IsNull: false}
				dut1, _ := sut.Collect(map[int]bool{CALL_GLOBAL: true})

				So(dut1["global/stat1"], ShouldAlmostEqual, 10, 0.1)

				timeNow = func() time.Time { return time.Unix(102, 0) }

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 20, Type: stats.TYPE_GAUGE, IsNull: false}
				dut2, _ := sut.Collect(map[int]bool{CALL_GLOBAL: true})

				So(dut2["global/stat1"], ShouldAlmostEqual, 20, 0.1)

			})

			Convey("Derives are exposed as ratio of change to time", func() {

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 10, Type: stats.TYPE_DERIVE, IsNull: false}
				dut1, _ := sut.Collect(map[int]bool{CALL_GLOBAL: true})

				So(dut1["global/stat1"], ShouldAlmostEqual, 10, 0.1)

				timeNow = func() time.Time { return time.Unix(102, 0) }

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 20, Type: stats.TYPE_DERIVE, IsNull: false}
				dut2, _ := sut.Collect(map[int]bool{CALL_GLOBAL: true})

				So(dut2["global/stat1"], ShouldAlmostEqual, 5, 0.1)

			})

			Convey("Counters are exposed as ratio of change to time", func() {

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 10, Type: stats.TYPE_COUNTER, IsNull: false}
				dut1, _ := sut.Collect(map[int]bool{CALL_GLOBAL: true})

				So(dut1["global/stat1"], ShouldAlmostEqual, 10, 0.1)

				timeNow = func() time.Time { return time.Unix(102, 0) }

				(*mocked.statusPtr).(stats.Stats)["global/stat1"] = stats.Stat{Value: 20, Type: stats.TYPE_COUNTER, IsNull: false}
				dut2, _ := sut.Collect(map[int]bool{CALL_GLOBAL: true})

				So(dut2["global/stat1"], ShouldAlmostEqual, 5, 0.1)

			})

		})

	})

}

func mockStat(ns string) *interface{} {
	s := stats.Stats{}

	s[ns+"/stat1"] = stats.Stat{Value: 1, Type: stats.TYPE_GAUGE, IsNull: false}
	s[ns+"/stat2"] = stats.Stat{Value: 2, Type: stats.TYPE_DERIVE, IsNull: false}

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
