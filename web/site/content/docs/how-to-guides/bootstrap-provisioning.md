---
title: "Provision a device"
type: docs
description: >
    Provision a new device from the edge via bootstrapping.
weight: 4
---

Following the steps below you will provision a new device via a publicly available Eclipse Hono sandbox using Eclipse Kanto.
A simple Eclipse Hono northbound business application written in Python is
provided to explore the capabilities for device bootstrapping and monitoring.

First a bootstrapping request is sent from the edge.
The custom Python application handles the request by
provisioning a new device. Upon successful provisioning it sends back 
all mandatory remote communication, identification and authentication data.
On the edge side, the response is handled by updating the connection configuration with the received data
and by executing a basic 
{{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/post_script.sh" %}}`post_script.sh`{{% /refn %}}
script to restart the Suite Connector service for the changes to take effect.

### Before you begin

To ensure that your edge device is capable to execute the steps in this guide, you need:

* Debian-based linux distribution
* If you don't have an installed and running Eclipse Kanto, follow {{% relrefn "install" %}} Install Eclipse Kanto {{% /relrefn %}}
* If you don't have a connected Eclipse Kanto to Eclipse Hono sandbox,
  follow {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}
* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/hono_commands_sb.py" %}} 
  suite bootstrapping application {{% /refn %}} and {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/post_script.sh" %}} 
  post script {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}
  guide are located and execute the following script:

  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands_sb.py && \
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/post_script.sh
  ```
* Grab the post script file and place it in the /var/tmp/suite-bootstrapping directory via executing:

  ```shell
  mkdir -p /var/tmp/suite-bootstrapping/ && sudo cp ./post_script.sh /var/tmp/suite-bootstrapping/
  ```
* Back up `/ect/suite-connector/config.json`

### Provision the Eclipse Hono tenant and device
In order to be able to provision a new device via bootstrapping, you need to have a dedicated tenant and a device instance
provisioned for it.
Fill in the required empty environmental variables definitions in the `hono_provisioning.sh`, e.g.:
{{% tip %}}
For this guide you need to provision a device different from the one you provisioned
in {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}
{{% /tip %}}
{{% tip %}}
It's nice to choose relatively unique and personalized values to avoid collisions in the public Eclipse Hono sandbox
{{% /tip %}}

```shell
# The Hono tenant to be created
export TENANT=demo.bootstrap
# The identifier of the device on the tenant
# Note! It's important for the ID to follow the convention namespace:name (e.g. demo:device)
export DEVICE_ID=demo.bootstrap:device
# The authentication identifier of the device
export AUTH_ID=demo.bootstrap_device
# A password for the device to authenticate with
export PWD=secret
```

Run the provisioning script and you will have your Eclipse Hono tenant and device ready to be connected:

```shell
./hono_provisioning.sh
```
### Configure Eclipse Kanto

Eclipse Kanto uses the `/etc/suite-bootstrapping/config.json` to acquire all the remote communication, identification and
authentication data to establish the remote connection for bootstrapping.
It is also where you need to specify the path to the post bootstrapping script and where to store received response data.
Update it with the following:


```json
{
  "logFile": "/var/log/suite-bootstrapping/suite-bootstrapping.log",
  "postBootstrapFile": "/etc/suite-connector/config.json",
  "postBootstrapScript": ["/var/tmp/suite-bootstrapping/post_script.sh"],
  "address": "hono.eclipseprojects.io:1883",
  "tenantId": "demo.bootstrap",
  "deviceId": "demo.bootstrap:device",
  "authId": "demo.bootstrap_device",
  "password": "secret"
}
```

Restart the Suite Bootstrapping service for the changes to take effect:

```shell
sudo systemctl restart suite-bootstrapping.service
```
When configured correctly the Suite Bootstrapping service automatically sends the bootstrapping request.

### Provision via bootstrapping
To explore the suite bootstrapping, we will use a Python script to provision and monitor the new device.
The location where the Python application will run does not have to be your edge device as it communicates remotely
with Eclipse Hono only.

Now we are ready to handle the bootstrapping request via executing the application
that requires the Eclipse Hono tenant (`-t`), the device identifier (`-d`) and the password you wish to use for the new device (`-p`) :

```shell
python3 hono_commands_sb.py -t demo.bootstrap -d demo.bootstrap:device -p secret
```

### Verify

You can check out that the Suite Connector is now connected to the new device via its log files
located in `/var/log/suite-connector/suite-connector.log`.

### Clean up
Revert the backed up `/ect/suite-connector/config.json` file.

Restart the Suite Connector service by executing:

```shell
sudo systemctl restart suite-connector.service
```
