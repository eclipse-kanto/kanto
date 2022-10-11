---
title: "Back up and restore files"
type: docs
description: >
    Back up and restore a file from and to your edge device.
weight: 3
---

Following the steps below you will back up a simple text file to an HTTP file server
and then restore it back via a publicly available Eclipse Hono sandbox using Eclipse Kanto.
A simple Eclipse Hono northbound business application written in Python is
provided to explore the capabilities for remotely backing up and restoring files.

### Before you begin

To ensure that all steps in this guide can be executed, you need:

* {{% refn "https://github.com/sebageek/servefile/" %}}`servefile`{{% /refn %}} installed

  This is a small Python HTTP server used in the example to serve the uploads and downloads.
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
  file backup and restore application {{% /refn %}}

  Navigate to the `quickstart` folder where the resources from the {{% relrefn "hono" %}} Explore via Eclipse Hono
  {{% /relrefn %}} guide are located and execute the following script:

  ```shell
  wget https://github.com/eclipse-kanto/kanto/raw/main/quickstart/hono_commands_fb.py
  ```

### Back up

By default, all directories in `/var/tmp/file-backup/` or the directory itself can be backed up.
For this example, create a file `data.txt` which will be later backed up:

```shell
sudo mkdir -p /var/tmp/file-backup && sudo echo "This is the first line in the file!" >> /var/tmp/file-backup/data.txt
```

You can verify that the file was successfully created by executing the following command:

```shell
cat /var/tmp/file-backup/data.txt
```

This should produce `This is the first line in the file!` as an output.

Choose a directory where the text file will be uploaded, open a new terminal there and run `servefile`
with the flag `-u` to enable a file upload:

```shell
servefile -u .
```

To explore the file backup, we will use a Python script to request and monitor the operation.
The location where the Python application will run does not have to be your edge device as it communicates remotely
with Eclipse Hono only.

Now we are ready to request the text file backup from the edge via executing the application that requires the command
to execute (`backup`), Eclipse Hono tenant (`-t`), the device identifier (`-d`) and the host where the backup will
be uploaded to:

```shell
python3 hono_commands_fb.py backup -t demo -d demo:device -h localhost
```

You can check out that the backup file `data.zip` is on your HTTP file server by
listing the content of the `servefile` working directory.

### Restore

To explore the restore capabilities you will first modify the `data.txt` file, and then you will restore it to
the version before the changes by using the backup, that was created earlier.

You can modify the `data.txt` file with the following command:

```shell
sudo echo "This is the second line in the file!" >> /var/tmp/file-backup/data.txt
```

You can verify that the file was successfully updated by executing the following command:

```shell
cat /var/tmp/file-backup/data.txt
```

This output should be:
```text
This is the first line in the file!
This is the second line in the file!
```

Navigate to the terminal where `servefile` was started and terminate it.
Start it again with the flag `-l` to enable a file download:

```shell
servefile -l .
```

To explore the file restore, we will use a Python script to request and monitor the operation.
The location where the Python application will run does not have to be your edge device as it communicates remotely
with Eclipse Hono only.

Now we are ready to request the text file restore from the edge via executing the application that requires the command
to execute (`restore`), Eclipse Hono tenant (`-t`), the device identifier (`-d`) and the host where the backup file
will be downloaded from:

```shell
python3 hono_commands_fb.py restore -t demo -d demo:device -h localhost
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
