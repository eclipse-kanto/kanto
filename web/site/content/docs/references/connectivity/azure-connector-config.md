---
title: "Azure Connector configuration"
type: docs
description: >
  Customize the remote connectivity.
weight: 2
---

### Properties

To control all aspects of the azure connector behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| tenantId | string | defaultTenant | Tenant unique identifier that the device belongs to |
| connectionString | string ​| | The connection string for connectivity to Azure IoT Hub, the format is: `"HostName=newHostName.azure-devices.net;DeviceId=deviceId;SharedAccessKey=accessKey"` |
| sasTokenValidity | string | 1h | The validity period for the generated SAS token for device authentication. Positive integer number followed by a unit suffix, such as '300m', '1h', etc., time units are: m, h, d |
| idScope | string | | ID scope for Azure Device Provisioning service |
| **Remote connectivity - TLS** | | | |
| alpn | string[] | | TLS application layer protocol negotiation options space separated for cloud access |
| caCert | string | iothub.crt | PEM encoded CA certificates file |
| cert | string | | PEM encoded certificate file to authenticate to the MQTT endpoint |
| key | string | | PEM encoded unencrypted private key file to authenticate to the MQTT endpoint |
| **Remote connectivity - TLS over TPM** | | | |
| tpmDevice | string | | Path to the device file or the unix socket to access the TPM 2.0 |
| tpmHandle | int | | TPM 2.0 storage root key handle, the type is unsigned 64-bit integer |
| tpmKeyPub | string | | File path to the public part of the TPM 2.0 key |
| tpmKey | string | | File path to the private part of the TPM 2.0 key |
| **Local connectivity** | | | |
| localAddress | string | tcp://localhost:1883 | Address of the MQTT server/broker that the azure connector will connect for the local communication, the format is: `scheme://host:port` |
| localUsername | string | | Username that is a part of the credentials |
| localPassword | string | | Password that is a part of the credentials |
| **Local connectivity - TLS** | | | |
| localCACert | string | | PEM encoded CA certificates file |
| localCert | string | | PEM encoded certificate file to authenticate to the MQTT server/broker |
| localKey | string | | PEM encoded unencrypted private key file to authenticate to the MQTT server/broker |
| **Logging** | | | |
| logFile | string | logs/azure-connector.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or a higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Example

The minimal required configuration to connect.

```json
{
    "connectionString": "HostName=hostName.azure-devices.net;DeviceId=deviceId;SharedAccessKey=cGFzc3AvcKQ=",
    "caCert": "/etc/azure-connector/iothub.crt",
    "logFile": "/var/log/azure-connector/azure-connector.log"
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

{{% warn %}}
Be aware that some combinations may be incompatible
{{% /warn %}}

```json
{
    "tenantId": "defaultTenant",
    "connectionString": "",
    "sasTokenValidity": "1h",
    "idScope": "",
    "alpn" : [],
    "caCert": "iothub.crt",
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
    "logFile": "logs/azure-connector.log",
    "logLevel": "INFO",
    "logFileCount": 5,
    "logFileMaxAge": 28,
    "logFileSize": 2
}
```
