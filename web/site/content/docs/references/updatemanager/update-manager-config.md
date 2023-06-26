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
| domain | string | device | Specify the Domain of this update agent, used as MQTT topic prefix |
| domains | string | containers,self-update,safety-domain| Specify a comma-separated list of domains handled by the update manager |
| phaseTimeout | string | 10m | Specify the timeout for completing an Update Orchestration phase |
| rebootAfter | string | 30s | Specify the timeout in cron format to wait before a reboot process is initiated after successful update operation |
| rebootEnabled | bool | true | Specify a flag that controls the enabling/disabling of the reboot process after successful update operation |
| reportFeedbackInterval | string | 1m | Specify the time interval for reporting intermediate desired state feedback messages during an active update operation |
| currentStateDelay | string | 30s | Specify the time delay for reporting current state messages |
| **Domain agent** | | | |
| readTimeout | int | 1m | Specify the read timeout for the given domain |
| rebootRequired | bool | false | Specify the reboot required flag for the given domain |
| **Local connectivity** | | | |
| brokerUrl | string | tcp://localhost:1883 | Specify the MQTT broker URL to connect to |
| keepAlive | bool | 20000 | Specify the keep alive duration for the MQTT requests in milliseconds |
| acknowledgeTimeout | int | 15000 | Specify the acknowledgement timeout for the MQTT requests in milliseconds |
| clientUsername | string | | Specify the MQTT client username to authenticate with |
| clientPassword | string | | Specify the MQTT client password to authenticate with |
| connectTimeout | int | 30000 | Specify the connect timeout for the MQTT in milliseconds |
| disconnectTimeout | int | 250 | Specify the disconnection timeout for the MQTT connection in milliseconds |
| subscribeTimeout | int | 15000 | Specify the subscribe timeout for the MQTT requests in milliseconds |
| unsubscribeTimeout | int | 5000 | Specify the unsubscribe timeout for the MQTT requests in milliseconds |
| **Logging** | | | |
| logFile | string | | Set the log file |
| logLevel | string | INFO | Set the log level - possible values are ERROR, WARN, INFO, DEBUG, TRACE |
| logFileCount | int | 5 | Set the maximum number of old log files to retain |
| logFileMaxAge | int | 28 | Set the maximum number of days to retain old log files based on the timestamp encoded in their filename |
| logFileSize | int | 2 | Set the maximum size in megabytes of the log file before it gets rotated |

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

```json
{
	"domain": "device",
	"domains": "containers,self-update,safety-domain",
	"agents": {
		"containers": {
			"name": "containers",
			"rebootRequired": false,
			"readTimeout": "30s"
		},
		"self-update": {
			"name": "self-update",
			"rebootRequired": true,
			"readTimeout": "30s"
		},
		"safety-domain": {
			"name": "safety-domain",
			"rebootRequired": true,
			"readTimeout": "30s"
		}
	},
	"log": {
		"logFile": "",
		"logLevel": "DEBUG",
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
