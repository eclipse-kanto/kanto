---
title: "File upload"
type: docs
description: >
  Empower the edge device to upload files to various storage providers.
weight: 4
---

File upload enables sending of files to a backend storage of choice. It can be used both locally and remotely via a desired IoT cloud ecosystem. The following use cases are provided:

* **Storage diversity** - with ready to use integrations with Azure Blob Storage, Amazon S3 and standard HTTP upload
* **Automatic uploads** - with periodically triggered uploads at a specified interval in a given time frame
* **Data integrity** - with an option to calculate and send the integrity check required information
* **Operation monitoring** - with a status reporting of the upload operation

![File upload](/kanto/images/docs/concepts/file-upload.png)

### How it works

It's not always possible to inline all the data into exchanged messages. For example, large log files or large diagnostic files cannot be sent as a telemetry message. In such scenarios, file upload can assist enabling massive amount of data to be stored to the backend storage.

There are different triggers which can initiate the upload operation: periodic or explicit. Once initiated, the request will be sent to the IoT cloud for confirmation or cancellation transferred back to the edge. If starting is confirmed, the files 
to upload will be selected according to the specified configuration, their integrity check information can be calculated and the transfer of the binary content will begin. A status report is announced on each step of the upload process 
enabling its transparent monitoring.

### What's next

[How to upload files]({{< relref "upload-files" >}})