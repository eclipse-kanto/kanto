---
title: "Update software"
type: docs
description: >
    Install a Debian package on your edge device.
weight: 1
---

Following the steps below you will install a{{% refn "https://www.gnu.org/software/hello/" %}}`hello`{{% /refn %}}
Debian package via a publicly available Eclipse Hono sandbox using Eclipse Kanto.
A couple of simple Eclipse Hono northbound business applications written in Python are provided to explore
the capabilities for remotely installing and monitoring.
On the edge side, a basic
{{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/install_hello.sh" %}}`install_hello.sh`{{% /refn %}}
script will be downloaded and executed.

### Before you begin

To ensure that your edge device is capable to execute the steps in this guide, you need:

* Debian-based linux distribution and the `apt` command line tool
* If you don't have an installed and running Eclipse Kanto, follow {{% relrefn "hono" %}} Install Eclipse Kanto {{% /relrefn %}}
* If you don't have a connected Eclipse Kanto to Eclipse Hono sandbox,
  follow {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}

* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/hono_commands_su.py" %}} 
  software update application {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}
  guide are located and execute the following script:

  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands_su.py
  ```
* Executing `hello` in the terminal will return that the command is not found

### Install Debian package

To explore the software management, we will use two Python scripts to install and monitor the `hello` Debian package.
The location where the Python applications will run does not have to be your edge device as they communicate remotely
with Eclipse Hono only.

First, start the monitoring application that requires the configured Eclipse Hono tenant (`-t`) and will print all
received events triggered by the device:

```shell
python3 hono_events.py -t demo
```

In another terminal, we are ready to spin up a `hello` Debian package at the edge via executing the second application
that requires the Eclipse Hono tenant (`-t`) and the device identifier (`-d`):

```shell
python3 hono_commands_su.py -t demo -d demo:device
```

### Verify

You can check out that the new package is installed on your edge device via executing:

```shell
hello
```

The command now displays: `Hello, world!`

### Clean up

The installed `hello` Debian package can be removed via executing:

```shell
sudo apt remove hello
```
