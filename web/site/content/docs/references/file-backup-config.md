---
title: "File backup configuration"
type: docs
description: >
  Customize the files backup and restore to and from a backend storage.
weight: 6
---

### Properties

To control all aspects of the file backup behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| featureId | string | BackupAndRestore | Feature unique identifier in the scope of the edge digital twin |
| type | string | file | Type of the files that are backed up by this feature |
| context | string | edge | Context of the files backed up by this feature, unique in the scope of the `type` |
| dir | string | | Directory to be backed up |
| mode | string | strict | Restriction on directories that can be dynamically selected for a backup, the supported modes are: strict, lax and scoped |
| backupCmd | string | | Command to be executed before the backup is done |
| restoreCmd | string | | Command to be executed after the restore |
| singleUpload | bool | false | Forbid triggering of new backups when there is a backup in progress |
| checksum | bool | false | Send MD5 checksum for backed up files to ensure data integrity |
| stopTimeout | string | 30s | Time to wait for running backups to finish as a sequence of decimal numbers, each with optional fraction and a unit suffix, such as: 300ms, 1.5h, 10m30s, etc., time units are: ns, us (or µs), ms, s, m, h |
| keepUploaded | bool | false | Keep successfully uploaded backups locally |
| storage | string | ./storage | Directory where backups and downloads will be stored |
| **Upload/Download - TLS** | | | |
| serverCert| string | | PEM encoded certificate file for secure uploads and downloads |
| **Auto backup** | | | |
| active | bool | false | Activate periodic backups |
| activeFrom | string | | Time from which periodic backups should be active, in RFC 3339 format, if omitted (and `active` flag is set) current time will be used as start of the periodic backups |
| activeTill | string | | Time till which periodic backups should be active, in RFC 3339 format, if omitted (and `active` flag is set) periodic backups will be active indefinitely |
| period | string | 10h| Backup period as a sequence of decimal numbers, each with optional fraction and a unit suffix, such as: 300ms, 1.5h, 10m30s, etc., time units are: ns, us (or µs), ms, s, m, h |
| **Local connectivity** | | | |
| broker | string | tcp://localhost:1883 | Address of the MQTT server/broker that the file backup will connect for the local communication, the format is: `scheme://host:port` |
| username | string | | Username that is a part of the credentials |
| password | string | | Password that is a part of the credentials |
| **Local connectivity - TLS** | | | |
| caCert | string | | PEM encoded CA certificates file |
| cert | string | | PEM encoded certificate file to authenticate to the MQTT server/broker |
| key | string | | PEM encoded unencrypted private key file to authenticate to the MQTT server/broker |
| **Logging** | | | |
| logFile | string | log/file-backup.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or a higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Example

The minimal required configuration that enables backing up a directory and sets the file type to config.

```json
{
    "type": "config",
    "dir": "/var/tmp/file-backup",
    "mode": "scoped",
    "storage": "/var/lib/file-backup",
    "logFile": "/var/log/file-backup/file-backup.log"
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

{{% warn %}}
Be aware that some combinations may be incompatible.
{{% /warn %}}

```json
{
  "featureId": "BackupAndRestore",
  "type": "file",
  "context": "edge",
  "dir": "",
  "mode": "strict",
  "backupCmd": "",
  "restoreCmd": "",
  "singleUpload": false,
  "checksum": false,
  "stopTimeout": "30s",
  "keepUploaded": false,
  "storage": "./storage",
  "serverCert": "",
  "active": false,
  "activeFrom": "",
  "activeTill": "",
  "period": "10h",
  "broker": "tcp://localhost:1883",
  "username": "",
  "password": "",
  "caCert": "",
  "cert": "",
  "key": "",
  "logFile": "log/file-backup.log",
  "logLevel": "INFO",
  "logFileCount": 5,
  "logFileMaxAge": 28,
  "logFileSize": 2
}
```
