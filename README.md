# snap collector plugin - mysql
This plugin collects metrics from MySQL database.  

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.4+](https://golang.org/dl/)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download mysql plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-mysql  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-mysql.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

####Global config
Global configuration files are described in snap's documentation. You have to add `"mysql"` section with following entries:

 - `"mysql_connection_string"` -  it's DSN with format described [here](https://github.com/go-sql-driver/mysql#dsn-data-source-name).  ex. `"root:r00tme@tcp(localhost:3306)/"` where `root` is username and `r00tme` is password, `localhost` is host address and `3306` is port where mysql is listening.
 -  `"mysql_use_innodb"` - possible values are `true` and `false`. Specifies if InnoDB statistics are collected. If you set this value to true and they are unavailable plugin will fail to start.

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Description (optional)
----------|-----------------------
/intel/mysql/bytes/buffer_pool_size |
/intel/mysql/bytes/ibuf_size |
/intel/mysql/bytes/metadata_mem_pool_size |
/intel/mysql/cache_result/qcache-hits |
/intel/mysql/cache_result/qcache-inserts |
/intel/mysql/cache_result/qcache-not_cached |
/intel/mysql/cache_result/qcache-prunes |
/intel/mysql/cache_size/qcache |
/intel/mysql/gauge/buffer_pool_bytes_data |
/intel/mysql/gauge/buffer_pool_bytes_dirty |
/intel/mysql/gauge/buffer_pool_pages_data |
/intel/mysql/gauge/buffer_pool_pages_dirty |
/intel/mysql/gauge/buffer_pool_pages_free |
/intel/mysql/gauge/buffer_pool_pages_misc |
/intel/mysql/gauge/buffer_pool_pages_total |
/intel/mysql/gauge/file_num_open_files |
/intel/mysql/gauge/innodb_activity_count |
/intel/mysql/gauge/innodb_dblwr_page_size |
/intel/mysql/gauge/trx_rseg_history_len |
/intel/mysql/mysql_bpool_bytes/data |
/intel/mysql/mysql_bpool_bytes/dirty |
/intel/mysql/mysql_bpool_counters/pages_flushed |
/intel/mysql/mysql_bpool_counters/read_ahead |
/intel/mysql/mysql_bpool_counters/read_ahead_evicted |
/intel/mysql/mysql_bpool_counters/read_ahead_rnd |
/intel/mysql/mysql_bpool_counters/read_requests |
/intel/mysql/mysql_bpool_counters/reads |
/intel/mysql/mysql_bpool_counters/write_requests |
/intel/mysql/mysql_bpool_pages/data |
/intel/mysql/mysql_bpool_pages/dirty |
/intel/mysql/mysql_bpool_pages/free |
/intel/mysql/mysql_bpool_pages/misc |
/intel/mysql/mysql_bpool_pages/total |
/intel/mysql/mysql_innodb_data/fsyncs |
/intel/mysql/mysql_innodb_data/read |
/intel/mysql/mysql_innodb_data/reads |
/intel/mysql/mysql_innodb_data/writes |
/intel/mysql/mysql_innodb_data/written |
/intel/mysql/mysql_innodb_dblwr/writes |
/intel/mysql/mysql_innodb_dblwr/written |
/intel/mysql/mysql_innodb_log/fsyncs |
/intel/mysql/mysql_innodb_log/waits |
/intel/mysql/mysql_innodb_log/write_requests |
/intel/mysql/mysql_innodb_log/writes |
/intel/mysql/mysql_innodb_log/written |
/intel/mysql/mysql_innodb_pages/created |
/intel/mysql/mysql_innodb_pages/read |
/intel/mysql/mysql_innodb_pages/written |
/intel/mysql/mysql_innodb_row_lock/time |
/intel/mysql/mysql_innodb_row_lock/waits |
/intel/mysql/mysql_innodb_rows/deleted |
/intel/mysql/mysql_innodb_rows/inserted |
/intel/mysql/mysql_innodb_rows/read |
/intel/mysql/mysql_innodb_rows/updated |
/intel/mysql/mysql_locks/lock_deadlocks |
/intel/mysql/mysql_locks/lock_row_lock_current_waits |
/intel/mysql/mysql_locks/lock_timeouts |
/intel/mysql/mysql_log_position/master-bin |
/intel/mysql/mysql_log_position/slave-exec |
/intel/mysql/mysql_log_position/slave-read |
/intel/mysql/mysql_log_position/time_offset |
/intel/mysql/mysql_octets/rx |
/intel/mysql/mysql_octets/tx |
/intel/mysql/operations/adaptive_hash_searches |
/intel/mysql/operations/buffer_data_reads |
/intel/mysql/operations/buffer_data_written |
/intel/mysql/operations/buffer_pages_created |
/intel/mysql/operations/buffer_pages_read |
/intel/mysql/operations/buffer_pages_written |
/intel/mysql/operations/buffer_pool_read_ahead |
/intel/mysql/operations/buffer_pool_read_ahead_evicted |
/intel/mysql/operations/buffer_pool_read_requests |
/intel/mysql/operations/buffer_pool_reads |
/intel/mysql/operations/buffer_pool_wait_free |
/intel/mysql/operations/buffer_pool_write_requests |
/intel/mysql/operations/dml_deletes |
/intel/mysql/operations/dml_inserts |
/intel/mysql/operations/dml_reads |
/intel/mysql/operations/dml_updates |
/intel/mysql/operations/ibuf_merges_delete |
/intel/mysql/operations/ibuf_merges_delete_mark |
/intel/mysql/operations/ibuf_merges_discard_delete |
/intel/mysql/operations/ibuf_merges_discard_delete_mark |
/intel/mysql/operations/ibuf_merges_discard_insert |
/intel/mysql/operations/ibuf_merges_discard_merges |
/intel/mysql/operations/ibuf_merges_insert |
/intel/mysql/operations/innodb_dblwr_pages_written |
/intel/mysql/operations/innodb_dblwr_writes |
/intel/mysql/operations/innodb_rwlock_s_os_waits |
/intel/mysql/operations/innodb_rwlock_s_spin_rounds |
/intel/mysql/operations/innodb_rwlock_s_spin_waits |
/intel/mysql/operations/innodb_rwlock_x_os_waits |
/intel/mysql/operations/innodb_rwlock_x_spin_rounds |
/intel/mysql/operations/innodb_rwlock_x_spin_waits |
/intel/mysql/operations/log_waits |
/intel/mysql/operations/log_write_requests |
/intel/mysql/operations/log_writes |
/intel/mysql/operations/os_data_fsyncs |
/intel/mysql/operations/os_data_reads |
/intel/mysql/operations/os_data_writes |
/intel/mysql/operations/os_log_bytes_written |
/intel/mysql/operations/os_log_fsyncs |
/intel/mysql/operations/os_log_pending_fsyncs |
/intel/mysql/operations/os_log_pending_writes |
/intel/mysql/threads/cached |
/intel/mysql/threads/connected |
/intel/mysql/threads/running |
/intel/mysql/total_threads/created |
/intel/mysql/mysql_commands/[subnamespace] | available namespaces are evaluated in runtime
/intel/mysql/mysql_handler/[subnamespace] | available namespaces are evaluated in runtime
/intel/mysql/mysql_locks/[subnamespace] | available namespaces are evaluated in runtime
/intel/mysql/mysql_select/[subnamespace] | available namespaces are evaluated in runtime
/intel/mysql/mysql_sort/[subnamespace] | available namespaces are evaluated in runtime


### Examples
Example running mysql, passthru processor, and writing data to a file.

This is done from the snap directory.

In one terminal window, open the snap daemon (in this case with logging set to 1 and trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```

In another terminal window:
Load mysql plugin
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-mysql
```
See available metrics for your system
```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file (e.g. `mem-file.json`):
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/mysql/threads/running": {},
                "/intel/mysql/mysql_commands/alter_tablespace": {},
                "/intel/mysql/operations/innodb_rwlock_s_spin_rounds": {}
            },
            "config": {},
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_mysql"
                    }
                }
            ]
        }
    }
}
```
Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
```

Create task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/mem-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

Stop task:
```
$ $SNAP_PATH/bin/snapctl task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-mysql/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-mysql/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@Lukasz Mroz](https://github.com/lmroz/)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
