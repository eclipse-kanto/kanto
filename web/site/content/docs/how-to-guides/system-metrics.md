---
title: "Monitor system metrics"
type: docs
description: >
    Monitor system metrics from your edge device.
weight: 3
---

Following the steps below you will be able to monitor the system metrics from your edge device
via a publicly available Eclipse Hono sandbox using Eclipse Kanto. A simple Eclipse Hono
northbound business application written in Python is provided to explore the capabilities
for remotely monitoring the CPU and memory utilization.

### Before you begin

To ensure that all steps in this guide can be executed, you need:

* {{% refn "https://plotly.com/" %}}`Plotly`{{% /refn %}} and
  {{% refn "https://plotly.com/dash/" %}}`Dash`{{% /refn %}} installed

  {{% refn "https://plotly.com/" %}}`Plotly`{{% /refn %}} is an open-source plotting library and
  {{% refn "https://plotly.com/dash/" %}}`Dash`{{% /refn %}} is a framework for building data application in Python.
  They are used in this example to run a simple HTTP server and visualize the incoming system metrics data
  in real time, and they do not have to be running on your edge device.
  You can install them by executing:

  ```shell
  pip3 install plotly dash
  ```

* If you don't have an installed and running Eclipse Kanto on your edge device,
  follow {{% relrefn "install" %}} Install Eclipse Kanto {{% /relrefn %}}
* If you don't have a connected Eclipse Kanto to Eclipse Hono sandbox,
  follow {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}

* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/hono_commands_sm.py" %}}
  system metrics application {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono
  {{% /relrefn %}} guide are located and execute the following script:

  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands_sm.py
  ```


### Monitor System Metrics

To explore the system metrics, we will use a Python script to request and monitor the
CPU and memory utilization. The location where the Python application will run does
not have to be your edge device as it communicates remotely with Eclipse Hono only.

Now we are ready to request the system metrics from the edge via executing the application
that requires the Eclipse Hono tenant (`-t`) and the device identifier (`-d`):

```shell
python3 hono_commands_sm.py -t demo -d demo:device
```

### Verify

You can check out that the CPU and memory utilization metrics are properly received and displayed
by checking out the application dashboard (by default - {{% refn "http://127.0.0.1:8050" %}}http://127.0.0.1:8050{{% /refn %}}).
