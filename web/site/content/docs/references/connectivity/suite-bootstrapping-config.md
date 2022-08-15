---
title: "Suite bootstrapping configuration"
type: docs
description: >
  Customize the automatic provisioning.
weight: 2
---

### Properties

To control all aspects of the suite bootstrapping behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| provisioningFile | string | provisioning.json | Path to the provisioning file, if {{% relrefn "dmp" %}}Bosch IoT Device Management{{% /relrefn %}} is in use |
| **Remote connectivity** | | | |
| address | string | mqtts://mqtt.bosch-iot-hub.com:8883 | Address of the MQTT endpoint that the suite connector will connect for the remote communication, the format is: `scheme://host:port` |
| deviceId | string | | Device unique identifier |
| authId | string | | Authentication unique identifier that is a part of the credentials |
| tenantId | string | | Tenant unique identifier that the device belongs to |
| password | string | | Password that is a part of the credentials |
| clientId | string | | MQTT client unique identifier |
| policyId | string | | Policy unique identifier of the digital twin |
| **Remote connectivity - TLS** | | | |
| cacert | string | iothub.crt | A PEM encoded CA certificates file |
| cert | string | | A PEM encoded certificate file for cloud access |
| key | string | | A PEM encoded unencrypted private key file for cloud access |
| deviceIdPattern | string | | Pattern to generate the device identifier, `{{subject-dn}}` and `{{subject-cn}}` placeholders can be part of it |
| **Remote connectivity - TLS over TPM** | | | |
| tpmDevice | string | | Path to the device file or the unix socket to access the TPM 2.0 |
| tpmHandle | int | | TPM 2.0 storage root key handle, the type is unsigned 64-bit integer |
| tpmKeyPub | string | | File path to the public part of the TPM 2.0 key |
| tpmKey | string | | File path to the private part of the TPM 2.0 key |
| **Local connectivity** | | | |
| localAddress | string | tcp://localhost:1883 | Address of the MQTT server/broker that the suite connector will connect for the local communication, the format is: `scheme://host:port` |
| localUsername | string | | Username that is a part of the credentials |
| localPassword | string | | Password that is a part of the credentials |
| **Bootstrapping** | | | |
| preBootstrapScript | string | | A file(s), containing pre-bootstrapping script(s), space separated arguments if any, executed before a bootstrapping request |
| preBootstrapFile | string | | A file path, that is used for output of pre-bootstrapping script(s) |
| postBootstrapScript | string | | A file(s), containing post-bootstrapping script(s), space separated arguments if any, executed after a bootstrapping request |
| postBootstrapFile | string | | A file path, that is used for output of post-bootstrapping script(s) |
| bootstrapProvisioningFile | string | | A file path, where bootstrapping response data is stored as JSON |
| maxChunkSize | int | 46080 | Maximum chunk size of the request data in bytes |
| **Logging** | | | |
| logFile | string | log/suite-bootstrapping.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

**Example**

The minimal required configuration that enables automatic provisioning.

```json
{
    "address":"hono.eclipseprojects.io:1883",
    "cacert": "/etc/suite-bootstrapping/iothub.crt",
    "tenantId": "org.eclipse.kanto",
    "deviceId": "org.eclipse.kanto:exampleDevice",
    "authId": "org.eclipse.kanto_example",
    "password": "secret",
    "logFile": "/var/log/suite-bootstrapping/suite-bootstrapping.log"
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
    "provisioningFile": "provisioning.json",
    "deviceId": "",
    "deviceIdPattern": "",
    "address": "mqtts://mqtt.bosch-iot-hub.com:8883",
    "tenantId": "",
    "policyId": "",
    "authId": "",
    "clientId": "",
    "cacert": "iothub.crt",
    "password": "",
    "cert": "",
    "key": "",
    "tpmDevice": "",
    "tpmHandle": 0,
    "tpmKeyPub": "",
    "tpmKey": "",
    "localAddress": "tcp://localhost:1883",
    "localUsername": "",
    "localPassword": "",
    "preBootstrapFile": "",
    "preBootstrapScript": "",
    "postBootstrapFile": "",
    "postBootstrapScript": "",
    "bootstrapProvisioningJson": "",
    "maxChunkSize": 46080,
    "logFile": "log/suite-bootstrapping.log",
    "logLevel": "INFO",
    "logFileCount": 5,
    "logFileMaxAge": 28,
    "logFileSize": 2
}
```
