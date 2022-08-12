---
title: "Suite Bootstrapping configuration"
type: docs
description: >
  Customize the automatic provisioning of devices.
weight: 6
---


### Properties

To control all aspects of the suite bootstrapping behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| provisioningFile | string | provisioning.json | Path to the provisioning file, if Bosch IoT Device Management is in use |
| **Remote connectivity** | | | |
| deviceId | string | | Device unique identifier |
| address | string | mqtts://mqtt.bosch-iot-hub.com:8883 | Address of the MQTT endpoint that the suite connector will connect to the hub, the format is: `scheme://host:port` |
| tenantId | string | | Tenant unique identifier that the device belongs to |
| policyId | string | | Policy unique identifier of the digital twin |
| authId | string | | Authentication unique identifier that is a part of the credentials |
| clientId | string | | Hub client unique identifier |
| password | string | | Password that is a part of the credentials |
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
| preBootstrapFile | string | | Pre-bootstrapping file |
| preBootstrapScript | script | | Pre-bootstrapping script, space separated arguments if any |
| postBootstrapFile | string | | Post-bootstrapping file |
| postBootstrapScript | script | | Post-bootstrapping script, space separated arguments if any |
| bootstrapProvisioningFile | string | | Location where bootstrapping provisioning JSON result is stored |
| maxChunkSize | int | 46080 | Maximum request chunk size in bytes |
| **Logging** | | | |
| logFile | string | log/suite-bootstrapping.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

**Example**

The minimal required configuration that enables automatic provisioning of devices.

```json
{
    "deviceId": "org.eclipse.kanto:exampleDevice",
    "tenantId": "t1b2d1efcd7d84ef4aee3951c2bb7a921_hub",
    "policyId": "org.eclipse.kanto:exampleDevice",
    "authId": "org.eclipse.kanto_example",
    "password": "example_password",
    "preBootstrapScript": "do_something.bat",
    "bootstrapProvisioningJson": "bootstrapping-provisioning.json",
    "logFile": "/var/log/suite-bootstrapping/suite-bootstrapping.log"
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

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
