---
title: "Local digital twins"
type: docs
description: >
  Empower the edge device with local digital twins for more advanced offline scenarios.
weight: 6
---

Local digital twins enables the digital twin state on a local level even in offline scenarios. It provides the following use cases:

* **Mirrors the applications**
* **Persistency** - digital twins are stored locally
* **Cloud connectivity** - provide connectivity to the cloud (similar to [Suite connector]({{< relref "suite-connector" >}}))
* **Offline scenarios** - local applications stay fully operable as if the connection with the cloud is not interrupted
* **Synchronization** - when connection to the cloud is established  the last known digital twin state is synchronized

![Local Digital Twins](/kanto/images/docs/concepts/local-digital-twins.png)

### How it works

Similar to the [Suite connector]({{< relref "suite-connector" >}}) service the **local digital twins** service needs to establish a connection to the cloud. To do this this the edge has to be manually or automatically provisioned. This connection is then used as a channel to pass the edge telemetry and event messages. Once a connection is established the device can be operated via commands.

To ensure that the digital twin state is available on a local level even in offline mode (no matter of the connection state) the **local digital twins** service persist all changes locally. Such capabilities were implemented to support offline scenarios and advanced edge computing involving synchronization with the cloud after disruptions or outages. The synchronization mechanisms were also designed in a way to significantly reduce data traffic, and efficiently prevent data loss due to long-lasting disruptions.

Upon reconnection the **local digital twins** it will notifies the cloud of any changes since during the offline mode to to synchronize the digital twins state.

The synchronization works in both directions:
* Cloud -> Local digital twins - desired properties are updated if such changes are requested from the cloud.
* Local digital twins -> Cloud - local digital twins state is sent to the cloud (e.g. all current features, their reported properties values, and any removed features while there was no connection).

### What's next

[How to receive offline the structure of your edge device.]({{< relref "../how-to-guides/offline-edge-device" >}})