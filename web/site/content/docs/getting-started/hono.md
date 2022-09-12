---
title: "Explore via Eclipse Hono"
type: docs
description: >
  Connect and start managing your edge device via Eclipse Hono.
weight: 2
---

Following the steps below you will connect your first device to a publicly
available {{% refn "https://www.eclipse.org/hono/sandbox/" %}}Eclipse Hono sandbox{{% /refn %}} using Eclipse Kanto.
A couple of simple Eclipse Hono northbound business applications written in Python are provided to explore the capabilities for
remotely managing and monitoring your edge device.

### Before you begin

The location where the Python applications and utility shell scripts will run does not have to be your edge device as
they communicate remotely with Eclipse Hono only. To run them, you need:

* <a href="https://wiki.python.org/moin/BeginnersGuide/Download" target="_blank">Python 3</a>
  and {{% refn "https://pip.pypa.io/en/stable/installation/#" %}}pip3{{% /refn %}}
* The {{% refn "https://github.com/eclipse-kanto/kanto/tree/main/quickstart" %}} quickstart applications and
  provisioning scripts{{% /refn %}}

  You can execute the script below to download them automatically:
  ```shell
  mkdir quickstart && cd quickstart && \
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands.py && \
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_events.py && \
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/requirements.txt && \
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_provisioning.sh
  ```

* Required Python dependencies to run the scripts

  You can install them by using the downloaded `requirements.txt` file via executing:
  ```shell
  pip3 install -r requirements.txt
  ```

### Provision the Eclipse Hono tenant and device

In order to be able to connect your device to Eclipse Hono, you need to have a dedicated tenant and a device instance
provisioned for it. Fill in the required empty environmental variables definitions in the `hono_provisioning.sh`, e.g.:
{{% tip %}}
It's nice to choose relatively unique and personalized values to avoid collisions in the public Eclipse Hono sandbox
{{% /tip %}}

```shell
# The Hono tenant to be created
export TENANT=demo
# The identifier of the device on the tenant
# Note! It's important for the ID to follow the convention namespace:name (e.g. demo:device)
export DEVICE_ID=demo:device
# The authentication identifier of the device
export AUTH_ID=demo_device
# A password for the device to authenticate with
export PWD=secret
```

Run the provisioning script and you will have your Eclipse Hono tenant and device ready to be connected:

```shell
./hono_provisioning.sh
```

### Configure Eclipse Kanto

Eclipse Kanto uses the `/etc/suite-connector/config.json` to acquire all the remote communication, identification and
authentication data to establish the remote connection. Update it with the following:


```json
{
  "caCert": "/etc/suite-connector/iothub.crt",
  "provisioningFile": "/etc/suite-connector/provisioning.json",
  "logFile": "/var/log/suite-connector/suite-connector.log",
  "address":"hono.eclipseprojects.io:1883",
  "tenantId":"demo",
  "deviceId":"demo:device",
  "authId":"demo_device",
  "password":"secret"
}
```
{{% tip %}}
Instead of Eclipse Hono Sandbox a custom Eclipse Hono instance can be used by replacing the address value in accordance
with the Suite Connector configuration reference
{{% /tip %}}

Restart the Suite Connector service for the changes to take effect:

```shell
sudo systemctl restart suite-connector.service
```

### Verify

To explore remote containerized applications management, we
will use the two Python scripts to run, monitor and remove a simple {{% refn "https://www.influxdata.com/" %}}InfluxDB{{% /refn %}} container 
using its public container image available at {{% refn "https://hub.docker.com/_/influxdb" %}}Docker Hub{{% /refn %}}.

First, start the monitoring application that requires the configured Eclipse Hono tenant (`-t`) and will print all
received events triggered by the device:

```shell
python3 hono_events.py -t demo
```

In another terminal, we are ready to spin up an InfluxDB container instance at the edge via executing the second application
that requires the command to execute (`run`), the Eclipse Hono tenant (`-t`), device identifier (`-d`) and 
the full container image reference to use (`--img`):

```shell
python3 hono_commands.py run -t demo -d demo:device --img docker.io/library/influxdb:1.8.4
```

After the script exits with success, you can check out the new container running on your edge device via
executing:

```shell
sudo kanto-cm list
```

Looking at the terminal where the monitoring application is running, you will be able to see all the events triggered by
the operation. 

To remove the newly created container, execute the same application script
only this time with the `rm` command and the identifier of the container to remove (`--id`), e.g.:

```shell
python3 hono_commands.py rm -t demo -d demo:device --id e6f7fbea-0e95-433c-acc7-16ef21b9c033
```
