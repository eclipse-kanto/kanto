---
title: "File upload configuration"
type: docs
description: >
  Customize the files transfer to a backend storage.
weight: 4
---

### Properties

To control all aspects of the file upload behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| featureId | string | AutoUploadable | Feature unique identifier in the scope of the edge digital twin |
| type | string | file | Type of the files that are uploaded by this feature |
| context | string | edge | Context of the files that are uploaded by this feature, unique in the scope of the `type` |
| files | string | | Glob pattern to select the files for upload |
| singleUpload | bool | false | Forbid triggering of new uploads when there is an upload in progress |
| checksum | bool | false | Send MD5 checksum for uploaded files to ensure data integrity |
| stopTimeout | string | 30s | Time to wait for running uploads to finish as a sequence of decimal numbers, each with optional fraction and a unit suffix, such as: 300ms, 1.5h, 10m30s, etc., time units are: ns, us (or µs), ms, s, m, h |
| delete | bool | false | Delete successfully uploaded files |
| **Auto upload** | | | |
| active | bool | false | Activate periodic uploads |
| activeFrom | string | | Time from which periodic uploads should be active, in RFC 3339 format, if omitted (and `active` flag is set) current time will be used as start of the periodic uploads |
| activeTill | string | | Time till which periodic uploads should be active, in RFC 3339 format, if omitted (and `active` flag is set) periodic uploads will be active indefinitely |
| period | string | 10h | Upload period as a sequence of decimal numbers, each with optional fraction and a unit suffix, such as: 300ms, 1.5h, 10m30s, etc., time units are: ns, us (or µs), ms, s, m, h |
| **Local connectivity** | | | |
| broker | string | tcp://localhost:1883 | Address of the MQTT server/broker that the file upload will connect for the local communication, the format is: `scheme://host:port` |
| username | string | | Username that is a part of the credentials |
| password | string | | Password that is a part of the credentials |
| **Logging** | | | |
| logFile | string | log/file-upload.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Example

The minimal required configuration that sets the file type to log.

```json
{
    "type": "log",
    "files": "/var/tmp/file-upload/*.*",
    "logFile": "/var/log/file-upload/file-upload.log"
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

```json
{
    "featureId": "AutoUploadable",
    "type": "file",
    "context": "edge",
    "files": "",
    "dumpFiles": false,
    "singleUpload": false,
    "checksum": false,
    "stopTimeout": "30s",
    "delete": false,
    "active": false,
    "activeFrom": "",
    "activeTill": "",
    "period": "10h",
    "broker": "tcp://localhost:1883",
    "username": "",
    "password": "",
    "logFile": "log/file-upload.log",
    "logLevel": "INFO",
    "logFileCount": 5,
    "logFileMaxAge": 28,
    "logFileSize": 2
}
```
