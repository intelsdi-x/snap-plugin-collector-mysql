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
		//panic(">>>>>>>>>>" + err.Error() + "<<<<<<<<<<<<<")
		t.Error(err)
		//t.Fail()
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
			prep := mock.ExpectPrepare("SELECT name, count, type FROM information_schema.innodb_metrics WHERE status = 'enabled'")
			prep.Optional()
		} else {
			prep := mock.ExpectPrepare("SELECT name, count, type FROM information_schema.innodb_metrics WHERE status = 'enabled'")
			prep.Optional()
			prep.WillReturnError(smthErr)
		}
	}

	if !skip[3] {
		prep := mock.ExpectPrepare("SHOW MASTER STATUS")
		prep.Optional()
	}

	if !skip[4] {
		prep := mock.ExpectPrepare("SHOW SLAVE STATUS")
		prep.Optional()
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

			Convey("if it's not", func() {

				tesingNew(mock, "4.6.5-ubu", true, true, 1)
				mock.ExpectPrepare("SHOW STATUS")
				New(testingConnectionString)

				Convey("show status is prepared", func() {

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

		Convey("show master status is prepared", func() {

			tesingNew(mock, "5.5.5-ubu", true, true, 3)
			mock.ExpectPrepare("SHOW MASTER STATUS")
			New(testingConnectionString)

			assert(mock, t)

		})

		Convey("show slave status is prepared", func() {

			tesingNew(mock, "5.5.5-ubu", true, true, 4)
			mock.ExpectPrepare("SHOW SLAVE STATUS")
			New(testingConnectionString)

			assert(mock, t)

		})

	})
}

func TestGetMasterStatus(t *testing.T) {

	Convey("GetMasterStatus", t, func() {

		mock := testingMockConn()

		tesingNew(mock, "5.5.5-ubu", true, true, 3)

		Convey("queries database about master status", func() {

			mock.ExpectPrepare("SHOW MASTER STATUS")
			New(testingConnectionString)

			assert(mock, t)

		})

		Convey("if data is correct", func() {

			Convey("returns desired metric", func() {

				prep := mock.ExpectPrepare("SHOW MASTER STATUS")
				prep.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"File", "Position", "Binlog_Do_DB", "Binlog_Ignore_DB", "Executed_Gtid_Set"}).AddRow("xyz", 123, nil, nil, nil))
				sut, _ := New(testingConnectionString)

				dut, dut_err := sut.GetMasterStatus()

				assert(mock, t)

				Convey("with correct value", func() {

					So(dut["mysql_log_position/master-bin"].Str(), ShouldEqual, counter(123).Str())

				})

				Convey("with no error", func() {

					So(dut_err, ShouldBeNil)

				})

			})
		})

		Convey("if number of columns is incorrect", func() {

			prep := mock.ExpectPrepare("SHOW MASTER STATUS")
			prep.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"File", "Position"}).AddRow("xyz", 123))
			sut, _ := New(testingConnectionString)

			_, dut_err := sut.GetMasterStatus()

			Convey("returns error", func() {

				So(dut_err, ShouldNotBeNil)

			})

		})

		Convey("if query failed", func() {

			prep := mock.ExpectPrepare("SHOW MASTER STATUS")
			prep.ExpectQuery().WillReturnError(smthErr)
			sut, _ := New(testingConnectionString)

			_, dut_err := sut.GetMasterStatus()

			Convey("returns error", func() {

				So(dut_err, ShouldNotBeNil)

			})

		})

	})

}

func TestGetSlaveStats(t *testing.T) {
	Convey("GetSlaveStatus", t, func() {

		mock := testingMockConn()

		tesingNew(mock, "5.5.5-ubu", true, true, 4)

		Convey("queries database about slave status", func() {

			mock.ExpectPrepare("SHOW SLAVE STATUS")
			New(testingConnectionString)

			assert(mock, t)

		})

		Convey("if data is correct", func() {

			Convey("returns desired metric", func() {

				prep := mock.ExpectPrepare("SHOW SLAVE STATUS")
				prep.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"Slave_IO_State",
					"Master_Host",
					"Master_User",
					"Master_Port",
					"Connect_Retry",
					"Master_Log_File",
					"Read_Master_Log_Pos",
					"Relay_Log_File",
					"Relay_Log_Pos",
					"Relay_Master_Log_File",
					"Slave_IO_Running",
					"Slave_SQL_Running",
					"Replicate_Do_DB",
					"Replicate_Ignore_DB",
					"Replicate_Do_Table",
					"Replicate_Ignore_Table",
					"Replicate_Wild_Do_Table",
					"Replicate_Wild_Ignore_Table",
					"Last_Errno",
					"Last_Error",
					"Skip_Counter",
					"Exec_Master_Log_Pos",
					"Relay_Log_Space",
					"Until_Condition",
					"Until_Log_File",
					"Until_Log_Pos",
					"Master_SSL_Allowed",
					"Master_SSL_CA_File",
					"Master_SSL_CA_Path",
					"Master_SSL_Cert",
					"Master_SSL_Cipher",
					"Master_SSL_Key",
					"Seconds_Behind_Master",
					"Master_SSL_Verify_Server_Cert",
					"Last_IO_Errno",
					"Last_IO_Error",
					"Last_SQL_Errno",
					"Last_SQL_Error",
					"Replicate_Ignore_Server_Ids",
					"Master_Server_Id",
					"Master_UUID",
					"Master_Info_File",
					"SQL_Delay",
					"SQL_Remaining_Delay",
					"Slave_SQL_Running_State",
					"Master_Retry_Count",
					"Master_Bind",
					"Last_IO_Error_Timestamp",
					"Last_SQL_Error_Timestamp",
					"Master_SSL_Crl",
					"Master_SSL_Crlpath",
					"Retrieved_Gtid_Set",
					"Executed_Gtid_Set",
					"Auto_Position",
					"Replicate_Rewrite_DB",
					"Channel_Name",
					"Master_TLS_Version"}).AddRow(nil,
					"192.168.56.105",
					"repl",
					3306,
					60,
					"mysql-bin.000001",
					123,
					"ubuntu-relay-bin.000001",
					4,
					"mysql-bin.000001",
					"No",
					"Yes",
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
					0,
					nil,
					0,
					747,
					154,
					"None",
					nil,
					0,
					"No",
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
					"No",
					1593,
					"Fatalerror:TheslaveI/OthreadstopsbecausemasterandslavehaveequalMySQLserverUUIDs;theseUUIDsmustbedifferentforreplicationtowork.",
					0,
					nil,
					nil,
					1,
					nil,
					"/var/lib/mysql/master.info",
					0,
					nil,
					"Slavehasreadallrelaylog;waitingformoreupdates",
					86400,
					nil,
					"16010715:42:22",
					nil,
					nil,
					nil,
					nil,
					nil,
					0,
					nil,
					nil,
					nil))
				sut, _ := New(testingConnectionString)

				dut, dut_err := sut.GetSlaveStatus()

				assert(mock, t)

				Convey("with correct value", func() {

					So(dut["mysql_log_position/slave-read"].Str(), ShouldEqual, counter(123).Str())
					So(dut["mysql_log_position/slave-exec"].Str(), ShouldEqual, counter(747).Str())
					So(dut["mysql_log_position/time_offset"].Str(), ShouldEqual, gauge(nil).Str())

				})

				Convey("with no error", func() {

					So(dut_err, ShouldBeNil)

				})

			})
		})

		Convey("if number of columns is incorrect", func() {

			prep := mock.ExpectPrepare("SHOW SLAVE STATUS")
			prep.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"File", "Position"}).AddRow("xyz", 123))
			sut, _ := New(testingConnectionString)

			_, dut_err := sut.GetSlaveStatus()

			Convey("returns error", func() {

				So(dut_err, ShouldNotBeNil)

			})

		})

		Convey("if query failed", func() {

			prep := mock.ExpectPrepare("SHOW SLAVE STATUS")
			prep.ExpectQuery().WillReturnError(smthErr)
			sut, _ := New(testingConnectionString)

			_, dut_err := sut.GetSlaveStatus()

			Convey("returns error", func() {

				So(dut_err, ShouldNotBeNil)

			})

		})

	})

}
