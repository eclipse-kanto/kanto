---
title: "Configuration"
type: docs
description: >
  Customize the update manager.
weight: 1
---

### Properties

To control all aspects of the update manager.

| Property | Type | Default | Description |
| - | - | - | - |
| **General** | | | |
| domain | string | device | The domain of this update agent, used as MQTT topic prefix |
| domains | string | containers| A comma-separated list of domains handled by the update manager |
| phaseTimeout | string | 10m | Timeout as duration string for completing an Update Orchestration phase |
| rebootAfter | string | 30s | Time period in cron format to wait before a reboot process is initiated after successful update operation |
| rebootEnabled | bool | true | Enable the reboot process after successful update operation |
| reportFeedbackInterval | string | 1m | Time interval as duration string for reporting intermediate desired state feedback messages during an active update operation |
| currentStateDelay | string | 30s | Time interval as duration string for reporting current state messages |
| **Domain agents** | | | |
| name | string | | Domain name |
| readTimeout | string | 1m | Timeout as duration string for reading the current state for the domain |
| rebootRequired | bool | false | Require a reboot for the domain |
| **Local connectivity** | | | |
| brokerUrl | string | tcp://localhost:1883 | Address of the MQTT server/broker that the container manager will connect for the local communication, the format is: `scheme://host:port` |
| keepAlive | bool | 20000 | Keep alive duration in milliseconds for the MQTT requests |
| disconnectTimeout | int | 250 | Disconnect timeout in milliseconds for the MQTT server/broker |
| clientUsername | string | | Username that is a part of the credentials |
| clientPassword | string | | Password that is a part of the credentials |
| acknowledgeTimeout | int | 15000 | Acknowledge timeout in milliseconds for the MQTT requests |
| connectTimeout | int | 30000 | Connect timeout in milliseconds for the MQTT server/broker |
| subscribeTimeout | int | 15000 | Subscribe timeout in milliseconds for the MQTT requests |
| unsubscribeTimeout | int | 5000 | Unsubscribe timeout in milliseconds for the MQTT requests |
| **Logging** | | | |
| logFile | string | | Path to the file where the update manager’s log messages are written |
| logLevel | string | INFO | All log messages at this or a higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Example

An example for configuring the update manager with two domains - `containers` and `self-update`, report feedback interval at 30 seconds, and log, written to custom log file `update-manager.log` with
log level `DEBUG`.

```json
{
	"domains": "containers,self-update",
	"log": {
		"logFile": "update-manager.log",
		"logLevel": "DEBUG"
	},
	"reportFeedbackInterval": "30s"
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

```json
{
	"domain": "device",
	"domains": "containers",
	"agents": {
		"containers": {
			"name": "containers",
			"rebootRequired": false,
			"readTimeout": "30s"
		}
	},
	"log": {
		"logFile": "",
		"logLevel": "INFO",
		"logFileCount": 5,
		"logFileMaxAge": 28,
		"logFileSize": 2
	},
	"mqtt": {
		"brokerUrl": "tcp://localhost:1883",
		"keepAlive": 20000,
		"acknowledgeTimeout": 15000,
		"clientUsername": "",
		"clientPassword": "",
		"connectTimeout": 30000,
		"disconnectTimeout": 250,
		"subscribeTimeout": 15000,
		"unsubscribeTimeout": 5000
	},
	"phaseTimeout": "10m",
	"rebootAfter": "30s",
	"rebootEnabled": true,
	"reportFeedbackInterval": "1m",
	"currentStateDelay": "30s"
}
```
