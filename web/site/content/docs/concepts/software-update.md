---
title: "Software update"
type: docs
description: >
  Empower the edge device to handle diverse software updates.
weight: 3
---

Software update enables the deployment and management of various software artifacts, both locally and remotely via an IoT cloud ecosystem of choice. It provides the following use cases:

* **Robust download** - with a retry and resume mechanism when the network connection is interrupted
* **Artifact validation** - with an integrity validation of every downloaded artifact
* **Universal installation** - with customizable install scripts to handle any kind of software
* **Operation monitoring** - with a status reporting of the download and install operations

![Software update](/kanto/images/docs/concepts/software-update.png)

### How it works

When the install operation is received at the edge, the download process is initiated. Retrieving the artifacts will continue until they are stored at the edge or their size threshold is reached. If successful, the artifacts are validated for integrity and further processed by the configured script. It is responsible to apply the new software and finish the operation. A status report is announced on each step of the process enabling its transparent monitoring.

On start up, if there have been any ongoing operations, they will be automatically resumed as the operation state is persistently stored.

### What's next

[How to update software]({{< relref "update-software" >}})
