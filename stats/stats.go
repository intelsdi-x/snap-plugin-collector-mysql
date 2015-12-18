package stats

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"reflect"
	"strconv"
	"strings"
)

const (
	TYPE_GAUGE = iota
	TYPE_DERIVE
	TYPE_COUNTER
)

const (
	SLAVE_READ_MASTER_IDX         = 6
	SLAVE_IO_RUNNING_IDX          = 10
	SLAVE_SQL_RUNNING_IDX         = 11
	SLAVE_EXEC_MASTER_LOG_POS_IDX = 21
	SLAVE_SECONDS_BEHIND_IDX      = 32
)

type Stat struct {
	Value  int64
	Type   int
	IsNull bool
}

func toInt(ifc interface{}) int64 {
	if ifc == nil {
		return 0
	}
	if bs, isBs := ifc.([]uint8); isBs {
		v, err := strconv.Atoi(string(bs))
		if err != nil {
			panic(err)
		}
		return int64(v)
	}
	return reflect.ValueOf(ifc).Convert(reflect.TypeOf(int64(0))).Int()
}

func counter(val interface{}) Stat {
	return Stat{Value: toInt(val), Type: TYPE_COUNTER, IsNull: val == nil}
}

func derive(val interface{}) Stat {
	return Stat{Value: toInt(val), Type: TYPE_DERIVE, IsNull: val == nil}
}

func gauge(val interface{}) Stat {
	return Stat{Value: toInt(val), Type: TYPE_GAUGE, IsNull: val == nil}
}

type Stats map[string]Stat

func parseVersion(s string) uint {
	dotsStr := strings.Split(strings.TrimSpace(s), "-")[0]
	var a, b, c uint
	fmt.Sscanf(dotsStr, "%d.%d.%d", &a, &b, &c)
	return a*10000 + b*100 + c
}

type MySQLStats struct {
	db             *sql.DB
	version        uint
	supportsInnodb bool

	stats, innodb, master, slave *sql.Stmt
}

func New(connectionString string) (*MySQLStats, error) {
	var err error
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("sql open failed: %v", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("database connection cannot be established: %v", err)
	}

	resVer := db.QueryRow("SELECT VERSION()")

	var verStr string
	err = resVer.Scan(&verStr)

	if err != nil {
		return nil, fmt.Errorf("version request failed: %v", err)
	}

	ver := parseVersion(verStr)

	res := &MySQLStats{db: db, version: ver}

	if ver >= 50002 {
		res.stats, err = db.Prepare("SHOW GLOBAL STATUS")
	} else {
		res.stats, err = db.Prepare("SHOW STATUS")
	}

	if err != nil {
		return nil, fmt.Errorf("cannot prepare status statement: %v", err)
	}

	if ver >= 50600 {
		res.supportsInnodb = true
		res.innodb, err = db.Prepare("SELECT name, count, type FROM information_schema.innodb_metrics WHERE status = 'enabled'")
	}

	if err != nil {
		return nil, fmt.Errorf("cannot prepare innodb statement: %v", err)
	}

	res.master, err = db.Prepare("SHOW MASTER STATUS")
	if err != nil {
		return nil, fmt.Errorf("cannot prepare master status statement: %v", err)
	}

	res.slave, err = db.Prepare("SHOW SLAVE STATUS")
	if err != nil {
		return nil, fmt.Errorf("cannot prepare slave status statement: %v", err)
	}

	return res, nil

}

func (self *MySQLStats) GetStatus(parseInnodb bool) (Stats, error) {
	rows, err := self.stats.Query()
	if err != nil {
		return nil, fmt.Errorf("status request failed: %v", err)
	}

	stats := Stats{}

	for rows.Next() {
		var name string
		var value interface{}

		err = rows.Scan(&name, &value)
		if err != nil {
			return nil, fmt.Errorf("status request failed: %v", err)
		}

		switch {
		case strings.HasPrefix(name, "Com_") && !strings.HasPrefix(name, "Com_stmt_"):
			stats["mysql_commands/"+strings.TrimPrefix(name, "Com_")] = counter(value)

		case strings.HasPrefix(name, "Handler_"):
			stats["mysql_handler/"+strings.TrimPrefix(name, "Handler_")] = counter(value)

		case strings.HasPrefix(name, "Table_locks_"):
			stats["mysql_locks/"+strings.TrimPrefix(name, "Table_locks_")] = counter(value)

		case strings.HasPrefix(name, "Select_"):
			stats["mysql_select/"+strings.TrimPrefix(name, "Select_")] = counter(value)

		case strings.HasPrefix(name, "Sort_"):
			stats["mysql_sort/"+strings.TrimPrefix(name, "Sort_")] = counter(value)

		default:
			switch name {
			case "Qcache_hits":
				stats["cache_result/qcache-hits"] = derive(value)
			case "Qcache_inserts":
				stats["cache_result/qcache-inserts"] = derive(value)
			case "Qcache_not_cached":
				stats["cache_result/qcache-not_cached"] = derive(value)
			case "Qcache_lowmem_pruness":
				stats["cache_result/qcache-prunes"] = derive(value)
			case "Qcache_queries_in_cache":
				stats["cache_size/qcache"] = gauge(value)

			case "Bytes_received":
				stats["mysql_octets/rx"] = gauge(value)
			case "Bytes_sent":
				stats["mysql_octets/tx"] = gauge(value)

			case "Threads_running":
				stats["threads/running"] = gauge(value)
			case "Threads_connected":
				stats["threads/connected"] = gauge(value)
			case "Threads_cached":
				stats["threads/cached"] = gauge(value)
			case "Threads_created":
				stats["total_threads/created"] = derive(value)
			}

			if parseInnodb {
				switch name {
				case "Innodb_buffer_pool_pages_data":
					stats["mysql_bpool_pages/data"] = gauge(value)
				case "Innodb_buffer_pool_pages_dirty":
					stats["mysql_bpool_pages/dirty"] = gauge(value)
				case "Innodb_buffer_pool_pages_flushed":
					stats["mysql_bpool_counters/pages_flushed"] = counter(value)
				case "Innodb_buffer_pool_pages_free":
					stats["mysql_bpool_pages/free"] = gauge(value)
				case "Innodb_buffer_pool_pages_misc":
					stats["mysql_bpool_pages/misc"] = gauge(value)
				case "Innodb_buffer_pool_pages_total":
					stats["mysql_bpool_pages/total"] = gauge(value)
				case "Innodb_buffer_pool_read_ahead_rnd":
					stats["mysql_bpool_counters/read_ahead_rnd"] = counter(value)
				case "Innodb_buffer_pool_read_ahead":
					stats["mysql_bpool_counters/read_ahead"] = counter(value)
				case "Innodb_buffer_pool_read_ahead_evicted":
					stats["mysql_bpool_counters/read_ahead_evicted"] = counter(value)
				case "Innodb_buffer_pool_read_requests":
					stats["mysql_bpool_counters/read_requests"] = counter(value)
				case "Innodb_buffer_pool_reads":
					stats["mysql_bpool_counters/reads"] = counter(value)
				case "Innodb_buffer_pool_write_requests":
					stats["mysql_bpool_counters/write_requests"] = counter(value)
				case "Innodb_buffer_pool_bytes_data":
					stats["mysql_bpool_bytes/data"] = gauge(value)
				case "Innodb_buffer_pool_bytes_dirty":
					stats["mysql_bpool_bytes/dirty"] = gauge(value)
				case "Innodb_data_fsyncs":
					stats["mysql_innodb_data/fsyncs"] = counter(value)
				case "Innodb_data_read":
					stats["mysql_innodb_data/read"] = counter(value)
				case "Innodb_data_reads":
					stats["mysql_innodb_data/reads"] = counter(value)
				case "Innodb_data_writes":
					stats["mysql_innodb_data/writes"] = counter(value)
				case "Innodb_data_written":
					stats["mysql_innodb_data/written"] = counter(value)
				case "Innodb_dblwr_writes":
					stats["mysql_innodb_dblwr/writes"] = counter(value)
				case "Innodb_dblwr_pages_written":
					stats["mysql_innodb_dblwr/written"] = counter(value)
				case "Innodb_log_waits":
					stats["mysql_innodb_log/waits"] = counter(value)
				case "Innodb_log_write_requests":
					stats["mysql_innodb_log/write_requests"] = counter(value)
				case "Innodb_log_writes":
					stats["mysql_innodb_log/writes"] = counter(value)
				case "Innodb_os_log_fsyncs":
					stats["mysql_innodb_log/fsyncs"] = counter(value)
				case "Innodb_os_log_written":
					stats["mysql_innodb_log/written"] = counter(value)
				case "Innodb_pages_created":
					stats["mysql_innodb_pages/created"] = counter(value)
				case "Innodb_pages_read":
					stats["mysql_innodb_pages/read"] = counter(value)
				case "Innodb_pages_written":
					stats["mysql_innodb_pages/written"] = counter(value)
				case "Innodb_row_lock_time":
					stats["mysql_innodb_row_lock/time"] = counter(value)
				case "Innodb_row_lock_waits":
					stats["mysql_innodb_row_lock/waits"] = counter(value)
				case "Innodb_rows_deleted":
					stats["mysql_innodb_rows/deleted"] = counter(value)
				case "Innodb_rows_inserted":
					stats["mysql_innodb_rows/inserted"] = counter(value)
				case "Innodb_rows_read":
					stats["mysql_innodb_rows/read"] = counter(value)
				case "Innodb_rows_updated":
					stats["mysql_innodb_rows/updated"] = counter(value)
				}
			}
		}
	}

	return stats, nil

}

func (self *MySQLStats) GetInnodb() (Stats, error) {
	if !self.supportsInnodb {
		return nil, fmt.Errorf("innodb stats not supported on current version of mysql server")
	}

	rows, err := self.innodb.Query()
	if err != nil {
		return nil, fmt.Errorf("innodb request failed: %v", err)
	}

	stats := Stats{}

	for rows.Next() {
		var name string
		var value, type_dummy interface{}

		err = rows.Scan(&name, &value, &type_dummy)
		if err != nil {
			return nil, fmt.Errorf("innodb request failed: %v", err)
		}

		switch name {
		case "metadata_mem_pool_size":
			stats["bytes/metadata_mem_pool_size"] = gauge(value)
		case "lock_deadlocks":
			stats["mysql_locks/lock_deadlocks"] = derive(value)
		case "lock_timeouts":
			stats["mysql_locks/lock_timeouts"] = derive(value)
		case "lock_row_lock_current_waits":
			stats["mysql_locks/lock_row_lock_current_waits"] = derive(value)
		case "buffer_pool_size":
			stats["bytes/buffer_pool_size"] = gauge(value)
		case "buffer_pool_reads":
			stats["operations/buffer_pool_reads"] = derive(value)
		case "buffer_pool_read_requests":
			stats["operations/buffer_pool_read_requests"] = derive(value)
		case "buffer_pool_write_requests":
			stats["operations/buffer_pool_write_requests"] = derive(value)
		case "buffer_pool_wait_free":
			stats["operations/buffer_pool_wait_free"] = derive(value)
		case "buffer_pool_read_ahead":
			stats["operations/buffer_pool_read_ahead"] = derive(value)
		case "buffer_pool_read_ahead_evicted":
			stats["operations/buffer_pool_read_ahead_evicted"] = derive(value)
		case "buffer_pool_pages_total":
			stats["gauge/buffer_pool_pages_total"] = gauge(value)
		case "buffer_pool_pages_misc":
			stats["gauge/buffer_pool_pages_misc"] = gauge(value)
		case "buffer_pool_pages_data":
			stats["gauge/buffer_pool_pages_data"] = gauge(value)
		case "buffer_pool_bytes_data":
			stats["gauge/buffer_pool_bytes_data"] = gauge(value)
		case "buffer_pool_pages_dirty":
			stats["gauge/buffer_pool_pages_dirty"] = gauge(value)
		case "buffer_pool_bytes_dirty":
			stats["gauge/buffer_pool_bytes_dirty"] = gauge(value)
		case "buffer_pool_pages_free":
			stats["gauge/buffer_pool_pages_free"] = gauge(value)
		case "buffer_pages_created":
			stats["operations/buffer_pages_created"] = derive(value)
		case "buffer_pages_written":
			stats["operations/buffer_pages_written"] = derive(value)
		case "buffer_pages_read":
			stats["operations/buffer_pages_read"] = derive(value)
		case "buffer_data_reads":
			stats["operations/buffer_data_reads"] = derive(value)
		case "buffer_data_written":
			stats["operations/buffer_data_written"] = derive(value)
		case "os_data_reads":
			stats["operations/os_data_reads"] = derive(value)
		case "os_data_writes":
			stats["operations/os_data_writes"] = derive(value)
		case "os_data_fsyncs":
			stats["operations/os_data_fsyncs"] = derive(value)
		case "os_log_bytes_written":
			stats["operations/os_log_bytes_written"] = derive(value)
		case "os_log_fsyncs":
			stats["operations/os_log_fsyncs"] = derive(value)
		case "os_log_pending_fsyncs":
			stats["operations/os_log_pending_fsyncs"] = derive(value)
		case "os_log_pending_writes":
			stats["operations/os_log_pending_writes"] = derive(value)
		case "trx_rseg_history_len":
			stats["gauge/trx_rseg_history_len"] = gauge(value)
		case "log_waits":
			stats["operations/log_waits"] = derive(value)
		case "log_write_requests":
			stats["operations/log_write_requests"] = derive(value)
		case "log_writes":
			stats["operations/log_writes"] = derive(value)
		case "adaptive_hash_searches":
			stats["operations/adaptive_hash_searches"] = derive(value)
		case "file_num_open_files":
			stats["gauge/file_num_open_files"] = gauge(value)
		case "ibuf_merges_insert":
			stats["operations/ibuf_merges_insert"] = derive(value)
		case "ibuf_merges_delete_mark":
			stats["operations/ibuf_merges_delete_mark"] = derive(value)
		case "ibuf_merges_delete":
			stats["operations/ibuf_merges_delete"] = derive(value)
		case "ibuf_merges_discard_insert":
			stats["operations/ibuf_merges_discard_insert"] = derive(value)
		case "ibuf_merges_discard_delete_mark":
			stats["operations/ibuf_merges_discard_delete_mark"] = derive(value)
		case "ibuf_merges_discard_delete":
			stats["operations/ibuf_merges_discard_delete"] = derive(value)
		case "ibuf_merges_discard_merges":
			stats["operations/ibuf_merges_discard_merges"] = derive(value)
		case "ibuf_size":
			stats["bytes/ibuf_size"] = gauge(value)
		case "innodb_activity_count":
			stats["gauge/innodb_activity_count"] = gauge(value)
		case "innodb_dblwr_writes":
			stats["operations/innodb_dblwr_writes"] = derive(value)
		case "innodb_dblwr_pages_written":
			stats["operations/innodb_dblwr_pages_written"] = derive(value)
		case "innodb_dblwr_page_size":
			stats["gauge/innodb_dblwr_page_size"] = gauge(value)
		case "innodb_rwlock_s_spin_waits":
			stats["operations/innodb_rwlock_s_spin_waits"] = derive(value)
		case "innodb_rwlock_x_spin_waits":
			stats["operations/innodb_rwlock_x_spin_waits"] = derive(value)
		case "innodb_rwlock_s_spin_rounds":
			stats["operations/innodb_rwlock_s_spin_rounds"] = derive(value)
		case "innodb_rwlock_x_spin_rounds":
			stats["operations/innodb_rwlock_x_spin_rounds"] = derive(value)
		case "innodb_rwlock_s_os_waits":
			stats["operations/innodb_rwlock_s_os_waits"] = derive(value)
		case "innodb_rwlock_x_os_waits":
			stats["operations/innodb_rwlock_x_os_waits"] = derive(value)
		case "dml_reads":
			stats["operations/dml_reads"] = derive(value)
		case "dml_inserts":
			stats["operations/dml_inserts"] = derive(value)
		case "dml_deletes":
			stats["operations/dml_deletes"] = derive(value)
		case "dml_updates":
			stats["operations/dml_updates"] = derive(value)
		}

	}

	return stats, nil

}

func (self *MySQLStats) GetMasterStatus() (Stats, error) {
	rows, err := self.master.Query()
	if err != nil {
		return nil, fmt.Errorf("master request failed: %v", err)
	}

	stats := Stats{}

	for rows.Next() {
		var dummy_0, position, dummy_2, dummy_3, dummy_4 interface{}

		err = rows.Scan(&dummy_0, &position, &dummy_2, &dummy_3, &dummy_4)
		if err != nil {
			return nil, fmt.Errorf("master request failed: %v", err)
		}
		stats["mysql_log_position/master-bin"] = counter(position)

		return stats, nil
	}

	return nil, fmt.Errorf("master request returned 0 rows")
}

func (self *MySQLStats) GetSlaveStatus() (Stats, error) {
	rows, err := self.slave.Query()
	if err != nil {
		return nil, fmt.Errorf("slave request failed: %v", err)
	}

	stats := Stats{}

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("slave request failed: %v", err)
	}

	if len(cols) < 33 {
		return nil, fmt.Errorf("slave request failed: number of columns < 33: %d", len(cols))
	}

	for rows.Next() {

		fields := make([]interface{}, len(cols))
		fieldPtrs := make([]interface{}, len(fields))

		for i, _ := range fieldPtrs {
			fieldPtrs[i] = &fields[i]
		}

		err = rows.Scan(fieldPtrs...)
		if err != nil {
			return nil, fmt.Errorf("slave request failed: %v", err)
		}
		stats["mysql_log_position/slave-read"] = counter(fields[SLAVE_READ_MASTER_IDX])
		stats["mysql_log_position/slave-exec"] = counter(fields[SLAVE_EXEC_MASTER_LOG_POS_IDX])
		stats["mysql_log_position/time_offset"] = gauge(fields[SLAVE_SECONDS_BEHIND_IDX])

		return stats, nil
	}

	return nil, fmt.Errorf("slave request returned 0 rows")
}

func (self *MySQLStats) Close() error {
	return self.db.Close()
}
