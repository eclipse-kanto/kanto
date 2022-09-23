---
title: "Backup and restore files"
type: docs
description: >
    Back up and restore a file to and from your edge device.
weight: 5
---

Following the steps below you will back up a simple text file to an HTTP file server
and then restore it back via a publicly available Eclipse Hono sandbox using Eclipse Kanto.
A simple Eclipse Hono northbound business application written in Python is
provided to explore the capabilities for backing up and restoring files.

### Before you begin

To ensure that all steps in this guide can be executed, you need:

* {{% refn "https://github.com/sebageek/servefile/" %}}`servefile`{{% /refn %}} installed

  This is a small Python HTTP server used in the example to serve the uploads.
  It does not have to be running on your edge device, but it has to be accessible from there.
  You can install it by executing:

  ```shell
  pip3 install servefile
  ```

* If you don't have an installed and running Eclipse Kanto on your edge device,
  follow {{% relrefn "install" %}} Install Eclipse Kanto {{% /relrefn %}}
* If you don't have a connected Eclipse Kanto to Eclipse Hono sandbox,
  follow {{% relrefn "hono" %}} Explore via Eclipse Hono {{% /relrefn %}}

* The {{% refn "https://github.com/eclipse-kanto/kanto/blob/main/quickstart/hono_commands_fb.py" %}} 
  backup and restore application {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono 
  {{% /relrefn %}} guide are located and execute the following script:
  
  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands_fb.py
  ```

### Back up a text file

By default, all directories in `/var/tmp/file-backup/` directory or the directory itself can be backed up.
For this example, create a file `data.txt` which will be later backed up:

```shell
sudo mkdir -p /var/tmp/file-backup && sudo echo "This is the first line in the file!" >> /var/tmp/file-backup/data.txt
```

You can verify that the file was successfully created by executing the following command:

```shell
cat /var/tmp/file-backup/data.txt
```

This should produce `This is the first line in the file!` as an output.

Create a directory `storage` where the backup file will be uploaded and run `servefile`
in a new terminal with the flag `-u` to enable uploading files to the HTTP server:

```shell
mkdir storage && servefile -u storage
```

To explore the file backup, we will use a Python script to request and monitor the operation.
The location where the Python application will run does not have to be your edge device as it communicates remotely
with Eclipse Hono only.

Now we are ready to request the log file backup from the edge via executing the application
that requires the command to execute (`backup`), Eclipse Hono tenant (`-t`) and the device identifier (`-d`):

```shell
python3 hono_commands_fb.py backup -t demo -d demo:device
```

### Verify

You can check out that the archived backup file `data.zip` is on your HTTP file 
server listing the content of the `storage` directory.

### Modify and restore a text file

To explore the restore capabilities you will first modify the `data.txt` file, and then you will restore it
by using the backup, that was created earlier.

Before you begin, modify the `data.txt` file, created in the backup process, so that you can verify
later that it is restored to its original state:

```shell
sudo echo "This is the new line in the file!" >> /var/tmp/file-backup/data.txt
```

You can verify that the file was successfully updated by executing the following command:

```shell
cat /var/tmp/file-backup/data.txt
```

This output should be: 
```text
This is the first line in the file!
This is the new line in the file!
```

Navigate to the `storage` directory where `servefile` was started and terminate it.
Start in again with the flag `-l` to enable downloading files from the HTTP server:

```shell
servefile -l storage
```

To explore the file restore, we will use a Python script to request and monitor the operation.
The location where the Python application will run does not have to be your edge device as it communicates remotely
with Eclipse Hono only.

Now we are ready to request the log file backup from the edge via executing the application
that requires the command to execute (`restore`), Eclipse Hono tenant (`-t`) and the device identifier (`-d`):

```shell
python3 hono_commands_fb.py restore -t demo -d demo:device
```

### Verify

You can check out that the original file is restored by executing the following command:

```shell
cat /var/tmp/file-backup/data.txt
```

This should produce `This is the first line in the file!` as an output.

### Clean up

Stop `servefile` and clean up its working directory.
Remove the `data.txt` file from the `/var/tmp/file-backup` directory.
