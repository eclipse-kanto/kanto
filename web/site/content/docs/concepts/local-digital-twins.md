---
title: "Local digital twins"
type: docs
description: >
  Empower the edge device with local digital twins for more advanced offline scenarios.
weight: 5
---

Local digital twins enables the digital twin state on a local level even in offline scenarios. The following use cases are provided:

* **Mirrors the applications**
* **Provides diverse services for the applications**
* **Provides own persistency** - digital twins are stored locally
* **Connectivity** - provide connectivity to the cloud (similar to [Suite connector]({{< relref "suite-connector" >}}))
* **Offline scenarios** - local applications stay fully operable as if the connection with the cloud is not interrupted
* **Synchronization** - on a connectivity recovering synchronizes the last known digital twin state

![Local Digital Twins](/kanto/images/docs/concepts/local-digital-twins.png)

### How it works

Similar to [Suite connector]({{< relref "suite-connector" >}}) Local digital twins needs to establish a connection with the cloud, for this the edge has to be manually or automatically provisioned. The connection is used as a channel to pass the edge telemetry and event messages. The IoT cloud can control the edge via commands and responses.
Local digital twins answers these commands for all local applications as the cloud and forwards the commands to the cloud without waiting for his answer. 

Local digital twins creates their persistency to store locally all changed applications. If the connection with the cloud is interrupted the applications continue to be fully operable. All their changes are stored locally.

When the connection is recovering synchronization between the cloud and application is needed.
The local digital twins stores the last digital twins state collects all reported properties from the local applications, and upon reconnection, it will notify the cloud of any changes since the last connection.
The synchronization happens in both directions:
* Cloud -> Local digital twins - desired properties get an update from the cloud as they may have changed on the cloud side
* Local digital twins -> Cloud - local digital twins' state is sent to the cloud (e.g. all current features, their reported properties values, and any removed features while there was no connection)

### What's next

[How to receive offline the structure of your edge device.]({{< relref "local-digital-twins" >}})