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

package stats

import (
	"fmt"
	"testing"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func (self Stat) Str() string {
	return fmt.Sprintf("%v", self)
}

const testingConnectionString = "x:x@/asfd"

var smthErr = fmt.Errorf("smth")

func assert(mock sqlmock.Sqlmock, t *testing.T) {
	err := mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func testingMockConn() sqlmock.Sqlmock {
	orgSqlOpen := sqlOpen

	db, mock, err := sqlmock.New()
	mock.MatchExpectationsInOrder(false)

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, err
	}

	Reset(func() {
		sqlOpen = orgSqlOpen
	})

	return mock
}

func tesingNew(mock sqlmock.Sqlmock, version interface{}, innodb, global bool, skipStep ...int) {

	skip := map[int]bool{}

	for _, v := range skipStep {
		skip[v] = true
	}

	if !skip[0] {
		if version == nil {
			version = "5.9.9"
		}
		ver, ver_ok := version.(string)
		if ver_ok {
			mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow(ver))
		} else {
			mock.ExpectQuery("SELECT VERSION()").WillReturnError(smthErr)
		}
	}

	if !skip[1] {
		if global {
			mock.ExpectPrepare("SHOW GLOBAL STATUS")
		} else {
			mock.ExpectPrepare("SHOW STATUS")
		}
	}

	if !skip[2] {
		if innodb {
			mock.ExpectPrepare("SELECT name, count, type FROM information_schema.innodb_metrics WHERE status = 'enabled'")
			//prep.Optional()
		} else {
			prep := mock.ExpectPrepare("SELECT name, count, type FROM information_schema.innodb_metrics WHERE status = 'enabled'")
			//prep.Optional()
			prep.WillReturnError(smthErr)
		}
	}

	if !skip[3] {
		mock.ExpectPrepare("SHOW MASTER STATUS")
		//prep.Optional()
	}

	if !skip[4] {
		mock.ExpectPrepare("SHOW SLAVE STATUS")
		//prep.Optional()
	}
}

func TestMNew(t *testing.T) {
	Convey("New", t, func() {

		mock := testingMockConn()

		Convey("asks about mysql version", func() {

			Convey("when version is returned", func() {
				tesingNew(mock, "5.6.5-ubu", true, true, 0)
				mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("5.6.5-ubu"))
				_, dut := New(testingConnectionString)
				Convey("no error is returned", func() {
					So(dut, ShouldBeNil)
					assert(mock, t)

				})

			})

			Convey("when version is not returned", func() {

				tesingNew(mock, fmt.Errorf("xxx"), true, true, 0)
				_, dut := New(testingConnectionString)

				Convey("error is returned", func() {

					So(dut, ShouldNotBeNil)

				})

			})

		})

		Convey("prepares statements", func() {

			Convey("when version >= 5.0.2", func() {
				tesingNew(mock, "5.6.5-ubu", true, true, 1)
				mock.ExpectPrepare("SHOW GLOBAL STATUS")
				New(testingConnectionString)

				Convey("show global status is prepared", func() {

					assert(mock, t)

				})

			})
		})

		Convey("if mysql version >= 5.6.0", func() {

			tesingNew(mock, "5.6.5-ubu", true, true, 2)
			mock.ExpectPrepare(".*innodb.*")
			New(testingConnectionString)

			Convey("innodb query is prepared", func() {

				assert(mock, t)

			})

		})

		Convey("if mysql version < 5.6.0", func() {

			tesingNew(mock, "5.5.5-ubu", true, true, 2)
			mock.ExpectPrepare(".*innodb.*").WillReturnError(smthErr)
			_, dut := New(testingConnectionString)

			Convey("innodb query is not prepared", func() {

				Convey("so New() should not return error when preparation fails", func() {

					So(dut, ShouldBeNil)

				})

			})

		})

	})
}
