---
title: "Offline explore edge device"
type: docs
description: >
    Offline receive the structure of your edge device.
weight: 5
---

By following the steps below, you will get the JSON structure of your device with all its properties using Eclipse Kanto. 
A simple client written in Python is provided to display the structure.

### Before you begin

To ensure that all steps in this guide can be executed, you need:

* {{% refn "https://github.com/eclipse/paho.mqtt.python/" %}}`paho-mqtt`{{% /refn %}} installed

  This is the Eclipse Paho MQTT Python client library used in the example to retrieve the structure of your edge device.
  You can install it by executing:

  ```shell
  pip3 install paho-mqtt
  ```

* If you don't have an installed and running Eclipse Kanto on your edge device,
  follow [Install Eclipse Kanto]({{< relref "install" >}})

* Only the Local Digital Twins service must running. If there are running suite-connector, please stop it for correct executing

```shell
sudo systemctl stop suite-connector.service && \
sudo systemctl restart local-digital-twins.service
```

* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/kanto_commands_ldt.py" %}} eclipse kanto client {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono   {{% /relrefn %}} guide are located and execute the following script:

  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/kanto_commands_ldt.py
  ```

### Configure Local digital twins

Local digital twins uses the `/etc/local-digital-twins/config.json` to acquire all the remote communication, identification and
authentication data to establish the remote connection. Because there is no such remote connection the Local digital twins will work offline. Update it with the following:


```json
{
  "caCert": "/etc/local-digital-twins/iothub.crt",
  "logFile": "/var/log/local-digital-twins/local-digital-twins.log",
  "thingsDb": "/var/lib/local-digital-twins/thing.db",
  "tenantId": "demo",
  "deviceId": "demo:device",
  "authId": "demo_device",
  "password": "secret"
}
```

Restart the Local Digital Twins service for the changes to take effect:

```shell
sudo systemctl restart local-digital-twins.service
```

### Receive the structure of the edge device

Now we are ready to request the structure of the application via executing the Eclipse kanto client that requires the Local digital twins tenant (-t) and the device identifier (-d):

```shell
python3 kanto_commands_ldt.py -t demo -d demo:device
```

### Verify

A command with a similar value will be received in the terminal and status 200.

```json
{
  "topic": "_/_/things/twin/commands/retrieve",
  "headers": {
    "content-type": "application/vnd.eclipse.ditto+json",
    "correlation-id": "correlation-id",
    "reply-to": "command/demo",
    "response-required": false
  },
  "path": "/",
  "value": [
    {
      "thingId": "demo:device",
      "features": {
        "AutoUploadable": {...},
        "BackupAndRestore": {...},
        "Metrics": {...},
        "SoftwareUpdatable": {...}
      }
    },
    {
      "thingId": "demo:device:edge:containers",
      "features": {
        "ContainerFactory": {...},
        "Metrics": {...},
        "SoftwareUpdatable": {...}
      }
    }
  ],
  "status": 200
}
```