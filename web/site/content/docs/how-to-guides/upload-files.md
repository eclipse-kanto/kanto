---
title: "Upload files"
type: docs
description: >
    Upload a log file from your edge device.
weight: 2
---

Following the steps below you will upload an example log file to your HTTP file server
via a publicly available Eclipse Hono sandbox using Eclipse Kanto.
A simple Eclipse Hono northbound business application written in Python is
provided to explore the capabilities for remotely uploading and monitoring.

### Before you begin

To ensure that all steps in this guide can be executed, you need:

* {{% refn "https://github.com/sebageek/servefile/" %}}`servefile`{{% /refn %}} installed

  This is a small Python HTTP server used in the example to serve the uploads.
  It does not have to be running on your edge device but it has to be accessible from there.
  You can install it by executing:

  ```shell
  pip3 install servefile
  ```

* If you don't have an installed and running Eclipse Kanto on your edge device,
  follow [Install Eclipse Kanto]({{< relref "install" >}})
* If you don't have a connected Eclipse Kanto to Eclipse Hono sandbox,
  follow [Explore via Eclipse Hono]({{< relref "hono" >}})

* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/hono_commands_fu.py" %}} 
  file upload application {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the [Explore via Eclipse Hono]({{< relref "hono" >}})
  guide are located and execute the following script:
  
  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands_fu.py
  ```


### Upload log file

By default, all files in `/var/tmp/file-upload/` directory can be uploaded.
For example, grab the suite connector log file and place it in the directory via executing:

```shell
mkdir -p /var/tmp/file-upload/ && sudo cp /var/log/suite-connector/suite-connector.log /var/tmp/file-upload/
```

Choose a directory where the log file will be uploaded, open a new terminal there and run `servefile`:

```shell
servefile -u .
```

To explore the file upload, we will use a Python script to request and monitor the operation.
The location where the Python application will run does not have to be your edge device as it communicates remotely
with Eclipse Hono only.

Now we are ready to request the log file upload from the edge via executing the application
that requires the Eclipse Hono tenant (`-t`) and the device identifier (`-d`):

```shell
python3 hono_commands_fu.py -t demo -d demo:device
```

### Verify

You can check out that the log file is on your HTTP file server listing the content of `servefile` working directory.

### Clean up

Stop `servefile` and clean up its working directory.
