---
title: "File backup configuration"
type: docs
description: >
  Customize the file backup to and from a backend storage.
weight: 7
---

### Properties

To control all aspects of the file backup behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| featureId | string | BackupAndRestore | Feature unique identifier in the scope of the edge digital twin. |
| dir | string | | Directory to be backed up. |
| storage | string | ./storage | Directory where backups and downloads will be stored. |
| backupCmd | string | | Command to be executed before backup is done. |
| restoreCmd | string | | Command to be executed after restore. |
| keepUploaded | bool | false | Keep locally successfully uploaded backups. |
| type | string | file | Type of the files that are uploaded by this feature. |
| context | string | edge |Context of the files uploaded by the feature, unique in the scope of the type. |
| mode | string | strict | Restriction on files that can be dynamically selected for a backup, the supported modes are: strict, lax and scoped.|
| singleUpload | bool | false | Forbid triggering of new uploads when there is upload in progress. Trigger can be forced from the backend with the 'force' option. |
| checksum | bool | false | Send MD5 checksum for uploaded files to ensure data integrity. Computing checksums incurs additional CPU/disk usage. |
| stopTimeout | string | 30s | Time to wait for running uploads/backups to finish as a sequence of decimal numbers, each with optional fraction and a unit suffix, such as: 300ms, 1.5h, 10m30s, etc., time units are: ns, us (or µs), ms, s, m, h |
| **Auto upload** | | | |
| active | bool | false | Activate periodic backups. |
| activeFrom | string | | Time from which periodic backups should be active, in RFC 3339 format, if omitted (and `active` flag is set) current time will be used as start of the periodic backups. |
| activeTill | string | | Time till which periodic backups should be active, in RFC 3339 format, if omitted (and `active` flag is set) periodic backups will be active indefinitely. |
| period | string | 10h| Backup period as a sequence of decimal numbers, each with optional fraction and a unit suffix, such as: 300ms, 1.5h, 10m30s, etc., time units are: ns, us (or µs), ms, s, m, h. |
| **Local connectivity** | | | |
| broker | string | tcp://localhost:1883 | Address of the MQTT server/broker that the file upload will connect for the local communication, the format is: `scheme://host:port`. |
| username | string | | Username that is a part of the credentials. |
| password | string | | Password that is a part of the credentials. |
| **Logging** | | | |
| logFile | string | log/file-upload.log | Path to the file where log messages are written. |
| logLevel | string | INFO | All log messages at this or higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE. |
| logFileCount | int | 5 | Log file maximum rotations count. |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files. |
| logFileSize | int | 2 | Log file size in MB before it gets rotated. |


### Example

The minimal required configuration that backs up a directory from the file system.

```json
{
    "type": "directory",
    "dir": "/var/tmp/file-backup/myBackupDirectory",
    "logFile": "/var/log/file-backup/file-backup.log"
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

{{% warn %}}
Be aware that some combinations may require property removal.
{{% /warn %}}

```json
{
  "featureId": "BackupAndRestore",
  "dir": "",
  "storage": "./storage",
  "backupCmd": "",
  "restoreCmd": "",
  "keepUploaded": false,
  "type": "file",
  "context": "edge",
  "mode": "strict",
  "singleUpload": false,
  "checksum": false,
  "stopTimeout": "30s",
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
