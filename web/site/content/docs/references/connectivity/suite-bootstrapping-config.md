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
| preBootstrapScript | string | | Path to the script/command with optional space-separated arguments that is executed before a bootstrapping request, optionally producing `preBootstrapFile` |
| preBootstrapFile | string | | Path to the file used as a bootstrapping request data |
| postBootstrapScript | string | | Path to the script/command with optional space-separated arguments that is executed after a bootstrapping response, optionally consuming `postBootstrapFile` |
| postBootstrapFile | string | | Path to the file used for a bootstrapping response data |
| bootstrapProvisioningFile | string | | Path to the file that stores provisioning info from bootstrapping response |
| maxChunkSize | int | 46080 | Maximum chunk size of the request data in bytes |
| provisioningFile | string | provisioning.json | Path to the provisioning file, if {{% relrefn "dmp" %}}Bosch IoT Device Management{{% /relrefn %}} is in use |
| **Remote connectivity** | | | |
| address | string | mqtts://mqtt.bosch-iot-hub.com:8883 | Address of the MQTT endpoint that the suite bootstrapping will connect for the remote communication, the format is: `scheme://host:port` |
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
| **Logging** | | | |
| logFile | string | log/suite-bootstrapping.log | Path to the file where log messages are written |
| logLevel | string | INFO | All log messages at this or higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| logFileCount | int | 5 | Log file maximum rotations count |
| logFileMaxAge | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| logFileSize | int | 2 | Log file size in MB before it gets rotated |

**Example**

The minimal required configuration to connect the publicly available
{{% refn "https://www.eclipse.org/hono/sandbox/" %}}Eclipse Hono sandbox{{% /refn %}} and request automatic provisioning.

```json
{
    "address":"hono.eclipseprojects.io:1883",
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
    "preBootstrapScript": "",
    "preBootstrapFile": "",
    "postBootstrapScript": "",
    "postBootstrapFile": "",
    "bootstrapProvisioningFile": "",
    "maxChunkSize": 46080,
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
    "logFile": "log/suite-bootstrapping.log",
    "logLevel": "INFO",
    "logFileCount": 5,
    "logFileMaxAge": 28,
    "logFileSize": 2
}
```
