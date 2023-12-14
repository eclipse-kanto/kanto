---
title: "Bootstrap device"
type: docs
description: >
    Automatically provision your device via bootstrapping.
weight: 5
---

By following the steps below you will automatically provision a new device via a publicly available Eclipse Hono sandbox using Eclipse Kanto.
A simple Eclipse Hono northbound business application written in Python is
provided to explore the capabilities for device bootstrapping and automatically provision a new device.

First a bootstrapping request is sent from the edge.
The custom Python application handles the request by automatically
provisioning a new device. Upon successful automatically provisioning it sends back
all mandatory remote communication, identification and authentication data.
On the edge side, the response is handled by updating the connection configuration with the received data
and by executing a basic
{{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/post_bootstrap.sh" %}}`post_bootstrap.sh`{{% /refn %}}
script to restart the Suite Connector service for the changes to take effect.

### Before you begin

To ensure that your edge device is capable to execute the steps in this guide, you need:

* If you don't have an installed and running Eclipse Kanto, follow {{% relrefn "install" %}} Install Eclipse Kanto {{% /relrefn %}}
* If you don't have a connected Eclipse Kanto to Eclipse Hono sandbox,
  follow {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}
* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/hono_commands_sb.py" %}}
  suite bootstrapping application {{% /refn %}} and {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/post_bootstrap.sh" %}} 
  post bootstrap script{{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}
  guide are located and execute the following script:

  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands_sb.py && \
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/post_bootstrap.sh
  ```
* Grab the post script file and place it in the /var/tmp/suite-bootstrapping directory via executing:

  ```shell
  sudo mkdir -p /var/tmp/suite-bootstrapping/ && sudo cp ./post_bootstrap.sh /var/tmp/suite-bootstrapping/
  ```
* Back up `/etc/suite-connector/config.json` as this file will be modified from this guide
* Stop suite-connector.service. Suite bootstrapping automatically provision device and try to start the suite connector service with new device
  ```shell
  sudo systemctl stop suite-connector.service
  ```

### Configure Suite Bootstrapping

Open file `/etc/suite-connector/config.json`, copy `address`, `tenantId`, `deviceId`, `authId` and `password`.
```json
{
    ...
    "address": "mqtts://hono.eclipseprojects.io:8883",
    "tenantId": "demo",
    "deviceId": "demo:device",
    "authId": "demo_device",
    "password": "secret"
    ...
}
```
Bootstrapping uses the `/etc/suite-bootstrapping/config.json` to acquire all the remote communication, identification and
authentication data to establish the remote connection for bootstrapping.
It is also where you need to specify the path to the post bootstrapping script and where to store received response data.
Update the configuration as shown below and replace `tenantId`, `deviceId`, `authId` and `password`
with the settings that you copied in the previous step.

```json
{
  "logFile": "/var/log/suite-bootstrapping/suite-bootstrapping.log",
  "postBootstrapFile": "/etc/suite-connector/config.json",
  "postBootstrapScript": ["/var/tmp/suite-bootstrapping/post_bootstrap.sh"],
  "caCert": "/etc/suite-bootstrapping/iothub.crt",
  "address": "mqtts://hono.eclipseprojects.io:8883",
  "tenantId": "demo",
  "deviceId": "demo:device",
  "authId": "demo_device",
  "password": "secret"
}
```

Restart the suite bootstrapping service for the changes to take effect:

```shell
sudo systemctl restart suite-bootstrapping.service
```
When configured correctly the Suite Bootstrapping service automatically sends the bootstrapping request.

### Automatically provision via bootstrapping

To explore the suite bootstrapping, we will use a Python script to automatically provision and monitor the new device.
The location where the Python application will run does not have to be your edge device as it communicates remotely
with Eclipse Hono only.

Now we are ready to handle the bootstrapping request via executing the application
that requires the Eclipse Hono tenant (`-t`), the device identifier (`-d`) and the password (`-p`) you wish to use for the new device:

```shell
python3 hono_commands_sb.py -t demo -d demo:device -p secret
```

### Verify

The last event received for the application is with the new tenant id that is automatically provisioning for the Suite Connector.
You can check out that the Suite Connector is now connected to the new device via its status.

```shell
sudo systemctl status suite-connector.service
```

### Clean up

Revert previous back up `/etc/suite-connector/config.json` file.
Remove temporary directory for post bootstrap file /var/tmp/suite-bootstrapping via executing:
```shell
sudo rm -r -f /var/tmp/suite-bootstrapping/
```
Stop suite bootstrapping service and restart suite connector service by executing:
```shell
sudo systemctl stop suite-bootstrapping.service && \
sudo systemctl restart suite-connector.service
```