# snap collector plugin - mysql
This plugin collects metrics from MySQL database.  

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Global Config](#global-config)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/)  - needed only for building

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
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/README.md#global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/README.md#examples).

## Documentation

###Global config
Global configuration files are described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add `"mysql"` section with following entries:

 - `"mysql_connection_string"` -  it's DSN with format described [here](https://github.com/go-sql-driver/mysql#dsn-data-source-name).  ex. `"root:r00tme@tcp(localhost:3306)/"` where `root` is username and `r00tme` is password, `localhost` is host address and `3306` is port where mysql is listening.
 -  `"mysql_use_innodb"` - possible values are `true` and `false`. Specifies if InnoDB statistics are collected. If you set this value to true and they are unavailable plugin will fail to start.
 
See exemplary Global configuration files in [examples/configs/] (https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/examples/configs/).

### Collected Metrics

List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/METRICS.md).

### Example
Example running mysql and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

Create Global Config, see examples in [examples/configs/] (https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/examples/configs/).

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled and global configuration saved in config.json ):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0 --config config.json
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

Create a task manifest file  (exemplary files in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/examples/tasks/)):
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
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
```

Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/task.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

Stop previously created task:
```
$ $SNAP_PATH/bin/snapctl task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-mysql/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-mysql/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [snap Gitter channel](https://gitter.im/intelsdi-x/snap).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@Lukasz Mroz](https://github.com/lmroz/)
