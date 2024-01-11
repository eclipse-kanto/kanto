---
title: "Offline explore edge device"
type: docs
description: >
    Offline receive the structure of your edge device.
weight: 6
---

By following the steps below you will get the structure of the edge digital twins with all its features and properties using Eclipse Kanto.
A simple Eclipse Hono northbound business application written in Python is provided to display the things' and their features' structure.

### Before you begin

To ensure that your edge device is capable to execute the steps in this guide, you need:

* If you don't have an installed and running Eclipse Kanto on your edge device,
  follow {{% relrefn "install" %}} Install Eclipse Kanto {{% /relrefn %}}

* If you don't have a connected Eclipse Kanto to Eclipse Hono sandbox,
  follow {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}

* Stop `suite-connector.service`. The local digital twins is a replacement of suite connector that's why only one of the services need to running.

  ```shell
  sudo systemctl stop suite-connector.service
  ```

* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/hono_commands_ldt.py" %}} offline explore application {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono  {{% /relrefn %}} guide are located and execute the following script:

  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands_ldt.py
  ```

### Configure Local digital twins

Open file `/etc/suite-connector/config.json`, copy `tenantId`, `deviceId`, `authId` and `password`.
```json
{
    ...
    "tenantId": "demo",
    "deviceId": "demo:device",
    "authId": "demo_device",
    "password": "secret"
    ...
}
```
Local digital twins uses the `/etc/local-digital-twins/config.json` to acquire all the remote communication, identification and
authentication data to establish the remote connection. Update the configuration as shown below and
replace `tenantId`, `deviceId`, `authId` and `password` with the settings that you copied in the previous step.

```json
  {
    "logFile": "/var/log/local-digital-twins/local-digital-twins.log",
    "caCert": "/etc/local-digital-twins/iothub.crt",
    "thingsDb": "/var/lib/local-digital-twins/thing.db",
    "tenantId": "demo",
    "deviceId": "demo:device",
    "authId": "demo_device",
    "password": "secret"
  }
```

Start the local digital twins service with previous set configuration:

```shell
sudo systemctl start local-digital-twins.service
```

### Receive the structure of the edge device

Now we are ready to request the structure of the edge digital twins via executing the offline explore application that requires the local digital twins tenant (`-t`) and the device identifier (`-d`):

```shell
python3 hono_commands_ldt.py -t demo -d demo:device
```

### Verify

On the shell there will be output of the structure of the edge digital twins with all it's features and properties. Things with following identifiers will be presented:
* demo:device
* demo:device:edge:containers

### Clean up

Stop local digital twins service and start suite connector service by executing:
```shell
sudo systemctl stop local-digital-twins.service && \
sudo systemctl restart suite-connector.service
```