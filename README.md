DISCONTINUATION OF PROJECT. 

This project will no longer be maintained by Intel.

This project has been identified as having known security escapes.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.
# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**

# Snap collector plugin - MySQL
This plugin collects metrics from MySQL database.  

It's used in the [Snap framework](http://github.com:intelsdi-x/snap).

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
* [golang 1.6+](https://golang.org/dl/)  - needed only for building

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download mysql plugin binary:
You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-mysql

Clone repo into `$GOPATH/src/github.com/intelsdi-x/`

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-mysql.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/README.md#global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/README.md#examples).

## Documentation

###Global config
Global configuration files are described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPTELD_CONFIGURATION.md). You have to add `"mysql"` section with following entries:

 - `"mysql_connection_string"` -  it's DSN with format described [here](https://github.com/go-sql-driver/mysql#dsn-data-source-name).  ex. `"root:r00tme@tcp(localhost:3306)/"` where `root` is username and `r00tme` is password, `localhost` is host address and `3306` is port where mysql is listening.
 - `"mysql_use_innodb"` - possible values are `true` and `false`. Specifies if InnoDB statistics are collected. If you set this value to true and they are unavailable plugin will fail to start.
 
See exemplary Global configuration files in [examples/configs/] (https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/examples/configs/).

### Collected Metrics

List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/METRICS.md).

### Example
Example running mysql and writing data to a file using [snap-plugin-publisher-file](https://github.com/intelsdi-x/snap-plugin-publisher-file).

Create Global Config, see examples in [examples/configs/] (https://github.com/intelsdi-x/snap-plugin-collector-mysql/blob/master/examples/configs/).

Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started),
in one terminal window, run `snapteld` (in this case with logging set to 1, trust disabled and global configuration saved in config.json):
```
$ snapteld -l 1 -t 0 --config config.json
```

In another terminal window:

Download and load Snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-mysql/latest/linux/x86_64/snap-plugin-collector-mysql
$ snaptel plugin load snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-mysql
```

See available metrics for your system
```
$ snaptel metric list
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
        "/intel/mysql/*": {}
      },
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

Create a task:
```
$ snaptel task create -t examples/tasks/task.json
```

Stop previously created task:
```
$ snaptel task stop <task_id>
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-mysql/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-mysql/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@Lukasz Mroz](https://github.com/lmroz/)
