---
title: "Software update configuration"
type: docs
description: >
  Customize the deployment and management of software artifacts.
weight: 4
---

### Properties

To control all aspects of the software update behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| featureId | string | SoftwareUpdatable | Feature unique identifier in the scope of the edge digital twin |
| moduleType | string | software | Type of the software that is managed by this feature |
| artifactType | string | archive | Type of the artifact that is to be processed: archive or plain |
| install | string[] | | Absolute path to the install script/command and an optional sequence of additional flags/parameters |
| storageLocation | string | ./ | Path to the storage directory where the working files are stored |
| installDirs | string[] | | Directories where the artifacts will be stored |
| mode | string | strict | Restriction where the artifacts can be located, the supported modes are: strict, lax and scope |
| **Download - TLS** | | | |
| serverCert | string | | PEM encoded certificate file for secure downloads |
| **Download** | | | |
| downloadRetryCount | int| 0 | Number of retries, in case of a failed download |
| downloadRetryInterval | string | 5s | Interval between retries, in case of a failed download as a sequence of decimal numbers, each with optional fraction and a unit suffix, such as: 300ms, 1.5h, 10m30s, etc., time units are: ns, us (or µs), ms, s, m, h |
| **Local connectivity** | | | |
| broker | string | tcp://localhost:1883 | Address of the MQTT server/broker that the software update will connect for the local communication, the format is: `scheme://host:port` |
| username | string | | Username that is a part of the credentials |
| password | string | | Password that is a part of the credentials |
| **Logging** | | | |
| logFile | string | log/software-update.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Example

The minimal required configuration that sets the software type to firmware.

```json
{
    "moduleType": "firmware",
    "storageLocation": "/var/lib/software-update",
    "logFile": "/var/log/software-update/software-update.log"
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

```json
{
    "featureId": "SoftwareUpdatable",
    "moduleType": "software",
    "artifactType": "archive",
    "install": [],
    "storageLocation": "./",
    "installDirs": [],
    "mode": "strict",
    "serverCert": "",
    "downloadRetryCount": 0,
    "downloadRetryInterval": "5s",
    "broker": "tcp://localhost:1883",
    "username": "",
    "password": "",
    "logFile": "log/software-update.log",
    "logLevel": "INFO",
    "logFileCount": 5,
    "logFileMaxAge": 28,
    "logFileSize": 2
}
```
