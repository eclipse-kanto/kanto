---
title: "AWS Connector configuration"
type: docs
description: >
  Customize the remote connectivity.
weight: 1
---

### Properties

To control all aspects of the suite connector behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| topicFilter | string ​| | Regex filter used to block incoming messages by their topic |
| payloadFilters | string ​| | Regex filters used to exclude parts of the incoming messages payload |
| **Remote connectivity** | | | |
| address | string | | Address of the MQTT endpoint that the connector will connect for the remote communication, the format is: `scheme://host:port` |
| alpn | string[] | | TLS application layer protocol negotiation options space separated for cloud access |
| tenantId | string | default-tenant-id | Tenant unique identifier that the device belongs to |
| clientId | string | | MQTT client unique identifier |
| **Remote connectivity - TLS** | | | |
| caCert | string | aws.crt | PEM encoded CA certificates file |
| cert | string | | PEM encoded certificate file to authenticate to the MQTT endpoint |
| key | string | | PEM encoded unencrypted private key file to authenticate to the MQTT endpoint |
| **Remote connectivity - TLS over TPM** | | | |
| tpmDevice | string | | Path to the device file or the unix socket to access the TPM 2.0 |
| tpmHandle | int | | TPM 2.0 storage root key handle, the type is unsigned 64-bit integer |
| tpmKeyPub | string | | File path to the public part of the TPM 2.0 key |
| tpmKey | string | | File path to the private part of the TPM 2.0 key |
| **Local connectivity** | | | |
| localAddress | string | tcp://localhost:1883 | Address of the MQTT server/broker that the suite connector will connect for the local communication, the format is: `scheme://host:port` |
| localUsername | string | | Username that is a part of the credentials |
| localPassword | string | | Password that is a part of the credentials |
| **Local connectivity - TLS** | | | |
| localCACert | string | | PEM encoded CA certificates file |
| localCert | string | | PEM encoded certificate file to authenticate to the MQTT server/broker |
| localKey | string | | PEM encoded unencrypted private key file to authenticate to the MQTT server/broker |
| **Logging** | | | |
| logFile | string | log/suite-connector.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or a higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

{{% warn %}}
Be aware that some combinations may be incompatible
{{% /warn %}}

```json
{
    "topicFilter": "",
    "payloadFilters": [],
    "address": "",
    "alpn" : [],
    "tenantId": "default-tenant-id",
    "clientId": "",
    "caCert": "aws.crt",
    "cert": "",
    "key": "",
    "tpmDevice": "",
    "tpmHandle": 0,
    "tpmKeyPub": "",
    "tpmKey": "",
    "localAddress": "tcp://localhost:1883",
    "localUsername": "",
    "localPassword": "",
    "localCACert": "",
    "localCert": "",
    "localKey": "",
    "logFile": "logs/aws-connector.log",
    "logLevel": "INFO",
    "logFileCount": 5,
    "logFileMaxAge": 28,
    "logFileSize": 2
}
```
