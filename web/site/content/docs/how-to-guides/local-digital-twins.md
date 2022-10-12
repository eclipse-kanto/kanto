---
title: "Offline explore edge device"
type: docs
description: >
    Offline receive the structure of your edge device.
weight: 5
---

By following the steps below, you will get the structure of the edge digital twins with all its features and properties using Eclipse Kanto. 
A simple application written in Python is provided to display the things' and its features' structure.

### Before you begin

To ensure that all steps in this guide can be executed, you need:

* {{% refn "https://github.com/eclipse/paho.mqtt.python/" %}}`paho-mqtt`{{% /refn %}} installed

  This is the MQTT client used in the example to retrieve the structure of your edge devices.
  You can install it by executing:

  ```shell
  pip3 install paho-mqtt
  ```

* If you don't have an installed and running Eclipse Kanto on your edge device,
  follow [Install Eclipse Kanto]({{< relref "install" >}})

* Start local digital twins

  The local digital twins is a replacement of suite connector that's why we have to stop suite connector and start the local digital twins.

  ```shell
  sudo systemctl stop suite-connector.service
  ```

* Local digital twins uses the `/etc/local-digital-twins/config.json` to acquire all the remote communication, identification and
authentication data to establish the remote connection. Update it with the following:

  ```json
  {
    "tenantId": "demo",
    "deviceId": "demo:device",
    "authId": "demo_device",
    "password": "secret",
    "caCert": "/etc/local-digital-twins/iothub.crt",
    "thingsDb": "/var/lib/local-digital-twins/thing.db",
    "logFile": "/var/log/local-digital-twins/local-digital-twins.log"
  }
  ```

  Restart the local digital twins service for the changes to take effect:

  ```shell
  sudo systemctl restart local-digital-twins.service
  ```

* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/kanto_commands_ldt.py" %}} offline explore application {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono   {{% /relrefn %}} guide are located and execute the following script:

  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/kanto_commands_ldt.py
  ```

### Receive the structure of the edge device

Now we are ready to request the structure of the edge digital twins via executing the Eclipse kanto client that requires the local digital twins tenant (-t) and the device identifier (-d):

```shell
python3 kanto_commands_ldt.py -t demo -d demo:device
```

### Verify

On the shell there will be output of the structure of the edge digital twins with all it's features and properties. Things with following identifiers will be presented: 
* demo:device
* demo:device:edge:containers
