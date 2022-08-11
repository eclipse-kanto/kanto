---
title: "System metrics configuration"
type: docs
description: >
  Customize the reporting of system metrics.
weight: 5
---

### Properties

To control all aspects of the system metrics behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| frequency | string | | Initial frequency of reporting system metrics as a sequence of decimal numbers, each with optional fraction and a unit suffix, such as: 300ms, 1.5h, 10m30s, etc., time units are: ns, us (or µs), ms, s, m, h |
| **Local connectivity** | | | |
| broker | string | tcp://localhost:1883 | Address of the MQTT server/broker that the system metrics will connect for the local communication, the format is: `scheme://host:port` |
| username | string | | Username that is a part of the credentials |
| password | string | | Password that is a part of the credentials |
| **Logging** | | | |
| logFile | string | log/system-metrics.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Example

The minimal required configuration that enables reporting of system metrics.

```json
{
    "frequency": "60s",
    "logFile": "/var/log/system-metrics/system-metrics.log"
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

```json
{
    "frequency" : "0s"
    "broker": "tcp://localhost:1883",
    "username": "",
    "password": "",
    "logFile": "log/system-metrics.log",
    "logLevel": "INFO",
    "logFileCount": 5,
    "logFileMaxAge": 28,
    "logFileSize": 2
}
```
