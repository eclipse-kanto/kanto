---
title: "Local digital twins configuration"
type: docs
description: >
  Customize the local digital twins persistency, access and synchronization.
weight: 2
---

### Properties

To control all aspects of the local digital twins behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| thingsDb | string | things.db | Path to the file where digital twins will be stored |
| provisioningFile | string ​| provisioning.json | Path to the provisioning file, if {{% relrefn "dmp" %}}Bosch IoT Device Management{{% /relrefn %}} is in use |
| **Remote connectivity** | | | |
| address | string | mqtts://mqtt.bosch-iot-hub.com:8883 | Address of the MQTT endpoint that the local digital twins will connect for the remote communication, the format is: `scheme://host:port` |
| deviceId | string | | Device unique identifier |
| authId | string | | Authentication unique identifier that is a part of the credentials |
| tenantId | string | | Tenant unique identifier that the device belongs to |
| password | string | | Password that is a part of the credentials |
| clientId | string | | MQTT client unique identifier |
| policyId | string | | Policy unique identifier of the digital twin |
| **Remote connectivity - TLS** | | | |
| caCert | string | iothub.crt | PEM encoded CA certificates file |
| cert | string | | PEM encoded certificate file to authenticate to the MQTT endpoint |
| key | string | | PEM encoded unencrypted private key file to authenticate to the MQTT endpoint |
| deviceIdPattern | string | | Pattern to generate the device identifier, `{{subject-dn}}` and `{{subject-cn}}` placeholders can be part of it |
| **Remote connectivity - TLS over TPM** | | | |
| tpmDevice | string | | Path to the device file or the unix socket to access the TPM 2.0 |
| tpmHandle | int | | TPM 2.0 storage root key handle, the type is unsigned 64-bit integer |
| tpmKeyPub | string | | File path to the public part of the TPM 2.0 key |
| tpmKey | string | | File path to the private part of the TPM 2.0 key |
| **Local connectivity** | | | |
| localAddress | string | tcp://localhost:1883 | Address of the MQTT server/broker that the local digital twins will connect for the local communication, the format is: `scheme://host:port` |
| localUsername | string | | Username that is a part of the credentials |
| localPassword | string | | Password that is a part of the credentials |
| **Local connectivity - TLS** | | | |
| localCACert | string | | PEM encoded CA certificates file |
| localCert | string | | PEM encoded certificate file to authenticate to the MQTT server/broker |
| localKey | string | | PEM encoded unencrypted private key file to authenticate to the MQTT server/broker |
| **Logging** | | | |
| logFile | string | log/local-digital-twins.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or a higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

### Example

The minimal required configuration to enable the local digital twins and their synchronization with the publicly available {{% refn "https://www.eclipse.org/hono/sandbox/" %}}Eclipse Hono sandbox{{% /refn %}}.

```json
{
    "thingsDb": "/var/lib/local-digital-twins/thing.db",
    "address": "hono.eclipseprojects.io:1883",
    "tenantId": "org.eclipse.kanto",
    "deviceId": "org.eclipse.kanto:exampleDevice",
    "authId": "org.eclipse.kanto_example",
    "password": "secret",
    "logFile": "/var/log/local-digital-twins/local-digital-twins.log"
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
    "thingsDb": "things.db",
    "provisioningFile": "provisioning.json",
    "address": "mqtts://mqtt.bosch-iot-hub.com:8883",
    "deviceId": "",
    "authId": "",
    "tenantId": "",
    "password": "",
    "clientId": "",
    "policyId": "",
    "caCert": "iothub.crt",
    "cert": "",
    "key": "",
    "deviceIdPattern": "",
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
    "logFile": "log/local-digital-twins.log",
    "logLevel": "INFO",
    "logFileCount": 5,
    "logFileMaxAge": 28,
    "logFileSize": 2
}
```