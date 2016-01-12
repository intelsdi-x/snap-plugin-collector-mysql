package mysqlplugin

import (
	"fmt"
	"testing"

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
		return nil, args.Error(1)
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
		return nil, args.Error(1)
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
		return nil, args.Error(1)
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
		return nil, args.Error(1)
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

				mocked.status = nil

				_, dut_err2 := sut2.Discover()

				So(dut_err2, ShouldBeNil)

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

				Convey("returns no error", func() {

					So(dut_err, ShouldBeNil)

				})

			})

		})

		Convey("if innodb is enabled", func() {

			Convey("requests innodb data", func() {

				source.AssertCalled(t, "GetInnodb")

				Convey("returns error when request fails", func() {

					sut2 := NewCollector(&source, false)

					mocked.innodb = nil

					_, dut_err2 := sut2.Discover()

					So(dut_err2, ShouldBeNil)

				})

				Convey("if request succeed", func() {

					Convey("exposes data", func() {

						Convey("-", nil)

					})

					Convey("returns no error", func() {

						Convey("-", nil)

					})

				})

			})

		})

		Convey("if innodb is disabled", func() {

			//sut := NewCollector(&source, false)

			Convey("does not request innodb data", func() {

				Convey("-", nil)

			})

		})

		Convey("tries to request master data", func() {

			Convey("does not fail when master data is unavailable", func() {

				Convey("-", nil)

			})

			Convey("if master data is available", func() {

				Convey("exposes data", func() {

					Convey("-", nil)

				})

				Convey("returns no error", func() {

					Convey("-", nil)

				})

			})

		})

		Convey("tries to request slave data", func() {

			Convey("does not fail when slave data is unavailable", func() {

				Convey("-", nil)

			})

			Convey("if slave data is available", func() {

				Convey("exposes data", func() {

					Convey("-", nil)

				})

				Convey("returns no error", func() {

					Convey("-", nil)

				})

			})

		})

	})

}

func TestCollect(t *testing.T) {

}

/*
func (self *statsMock) mockStat(call, ns string, ptr *stats.Stats, args ...interface{}) {

	if ptr == statsError {
		self.On(call, args...).Return(nil, smthErr)
	} else {
		if ptr == nil {
			ptr = &stats.Stats{}
		}

		(*ptr)[ns+"/stat1"] = stats.Stat{Value: 1, Type: stats.TYPE_GAUGE, IsNull: false}
		(*ptr)[ns+"/stat2"] = stats.Stat{Value: 2, Type: stats.TYPE_DERIVE, IsNull: false}

		self.On(call, mock.Anything).Return(*ptr, nil)

	}
}*/

func mockStat(ns string) (stats.Stats, *interface{}) {
	s := stats.Stats{}

	s[ns+"/stat1"] = stats.Stat{Value: 1, Type: stats.TYPE_GAUGE, IsNull: false}
	s[ns+"/stat2"] = stats.Stat{Value: 2, Type: stats.TYPE_DERIVE, IsNull: false}

	var ifc interface{} = s

	return s, &ifc
}

type statMockData struct {
	status, innodb, master, slave             stats.Stats
	statusPtr, innodbPtr, masterPtr, slavePtr *interface{}
}

func newMockedStats() statMockData {
	self := statMockData{}
	self.status, self.statusPtr = mockStat("global")
	self.innodb, self.innodbPtr = mockStat("inno")
	self.master, self.masterPtr = mockStat("master")
	self.slave, self.slavePtr = mockStat("slave")

	return self
}

/*
func newMockedStats(status, innodb, master, slave *stats.Stats) *statsMock {
	self := new(statsMock)

	self.mockStat("GetStatus", "global", status, status, mock.Anything)
	self.mockStat("GetInnodb", "inno", innodb)
	self.mockStat("GetMasterStatus", "master", master)
	self.mockStat("GetSlaveStatus", "slave", slave)

	return self
}
*/
