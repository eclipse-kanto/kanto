---
title: "Update manager configuration"
type: docs
description: >
  Customize the update manager.
weight: 6
---

### Properties

To control all aspects of the update manager.

| Property | Type | Default | Description |
| - | - | - | - |
| **General** | | | |
| domain | string | device | The domain of the update manager, used as MQTT topic prefix |
| domains | string | containers| A comma-separated list of domains handled by the update manager. This configuration option is available only as a flag, but not inside the JSON config file. In JSON config file, the keys inside the Domain agents structure serve as domain names. |
| phaseTimeout | string | 10m | Timeout as duration string for completing an Update Orchestration phase |
| rebootAfter | string | 30s | Time period as duration string to wait before a reboot process is initiated after successful update operation |
| rebootEnabled | bool | true | Enable the reboot process after successful update operation |
| reportFeedbackInterval | string | 1m | Time interval as duration string for reporting intermediate desired state feedback messages during an active update operation |
| currentStateDelay | string | 30s | Time interval as duration string for reporting current state messages |
| thingsEnabled | bool | true | Enable the Update Manager to behave as a thing's feature |
| ownerConsentCommands | []string | | List of commands for which an owner consent should be granted. Possible values are: 'DOWNLOAD', 'UPDATE', 'ACTIVATE' |
| ownerConsentTimeout | string | 30m | Timeout as duration string to wait for owner consent" |
| **Domain agents** | | | Holds a map structure (_agents_) with update agent configurations where each map key is treated as domain name |
| readTimeout | string | 1m | Timeout as duration string for reading the current state for the domain |
| rebootRequired | bool | false | Require a reboot for the domain after successful update |
| **Local connectivity** | | | |
| broker | string | tcp://localhost:1883 | Address of the MQTT server/broker that the container manager will connect for the local communication, the format is: `scheme://host:port` |
| keepAlive | string | 20s | Keep alive duration for the MQTT requests as duration string |
| disconnectTimeout | string | 250ms | Disconnect timeout for the MQTT server/broker as duration string |
| username | string | | Username that is a part of the credentials |
| password | string | | Password that is a part of the credentials |
| acknowledgeTimeout | string | 15s | Acknowledge timeout for the MQTT requests as duration string |
| connectTimeout | string | 30s | Connect timeout for the MQTT server/broker as duration string |
| subscribeTimeout | string | 15s | Subscribe timeout for the MQTT requests as duration string |
| unsubscribeTimeout | string | 5s | Unsubscribe timeout for the MQTT requests as duration string |
| **Logging** | | | |
| logFile | string | | Path to the file where the update manager’s log messages are written |
| logLevel | string | INFO | All log messages at this or a higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Example

An example for configuring the update manager with two domains - `containers` and `custom-domain`, report feedback interval at 30 seconds, and log, written to custom log file `update-manager.log` with
log level `DEBUG`.

```json
{
	"log": {
		"logFile": "update-manager.log",
		"logLevel": "DEBUG"
	},
	"agents": {
		"containers": {
			"readTimeout": "30s"
		},
		"custom-domain": {
			"rebootRequired": true
		}
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
	"agents": {
		"containers": {
			"rebootRequired": false,
			"readTimeout": "1m"
		}
	},
	"log": {
		"logFile": "",
		"logLevel": "INFO",
		"logFileCount": 5,
		"logFileMaxAge": 28,
		"logFileSize": 2
	},
	"connection": {
		"broker": "tcp://localhost:1883",
		"keepAlive": "20s",
		"acknowledgeTimeout": "15s",
		"username": "",
		"password": "",
		"connectTimeout": "30a",
		"disconnectTimeout": "250ms",
		"subscribeTimeout": "15s",
		"unsubscribeTimeout": "5s"
	},
	"phaseTimeout": "10m",
	"rebootAfter": "30s",
	"rebootEnabled": true,
	"reportFeedbackInterval": "1m",
	"currentStateDelay": "30s",
	"thingsEnabled": true,
	"ownerConsentCommands": ["DOWNLOAD"]
}
```
