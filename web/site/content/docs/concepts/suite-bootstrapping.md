---
title: "Suite bootstrapping "
type: docs
description: >
  Empower the edge device to be automatically provisioned.
weight: 5
---

Suite bootstrapping provides a mechanism for automatic provisioning of devices that will be connected to an IoT cloud ecosystem of choice, powered by Eclipse Hono™ (e.g.  {{% refn "https://www.eclipse.org/packages/packages/cloud2edge" %}}Eclipse Cloud2Edge{{% /refn %}} and Bosch IoT Suite). It provides the following use cases:

* **Auto provisioning** - with zero-touch the device is configured automatically to be connected
* **Solution sufficiency** - without hardcoding the initial configuration the device will be ready to be connected
* **Load balance** - with an option to distribute the device across multiple field subscriptions
 
![Suite bootstrapping](/kanto/images/docs/concepts/suite-bootstrapping.png)

### How it works

The suite bootstrapping enables automatic provisioning of devices that will be connected without hardcoding the initial configuration. Before starting the bootstrap process, the device is specified to which IoT cloud subscription will be forwarded. Once establish the connection, the suite bootstrapping makes a request with optionally pre-configured data. The request will be sent for getting the device field subscription and the bootstrap data.

The bootstrap response is the provisioning information that is stored on the device. The optional post-bootstrapping logic will be executed in order to handle the post-provisioned data. The {{% relrefn "suite-connector" %}}suite connector{{% /relrefn %}} will use the bootstrapping provisioning data to prepare the device connectivity.

### What's next

[Bootstrap device]()