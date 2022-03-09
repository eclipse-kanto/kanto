---
title: "Explore via Bosch IoT Device Management"
type: docs
description: >
  Connect and start managing your edge device via Bosch IoT Device Management.
weight: 3
---

{{% refn "https://bosch-iot-suite.com/service/bosch-iot-device-management/" %}}Bosch IoT Device
Management{{% /refn %}}enables remote data, management and configuration of devices and fleets of devices in a unified and
flexible manner. Eclipse Kanto is Bosch IoT Device Management integrated out-of-the-box.

### Before you begin

To ensure you have prepared the backend where your device will be connected you will need to perform the following

* Create a free {{% refn "https://docs.bosch-iot-suite.com/device-management/Subscribe-a-service-instance.html" %}}Bosch
  IoT Device Management subscription{{% /refn %}} if you do not already have one
* {{% refn "https://docs.bosch-iot-suite.com/device-management/Register-a-gateway-device-via-the-Bosch-IoT-Suite-Console.html" %}}
Provision a gateway device{{% /refn %}} from scratch with automatic provisioning of edge devices enabled 
so that you are able to explore all capabilities provided by Eclipse Kanto

### Configure Eclipse Kanto

Eclipse Kanto uses the generated `provisioning.json` to acquire all the remote communication, identification and
authentication data to establish the remote connection. Once you have it available on the device
in `/etc/suite-connector/provisioning.json`, Eclipse Kanto will connect your device automatically.

### Verify

If the connection has been established successfully, you should be able to see two _Things_ with a green plugged-in icon
on the right-hand side

* **\<your-namespace\>:\<gateway-name\>** - the gateway you have provisioned with _Features_, e.g.
    * ConnectionStatus 
    * Autouploadable
    * SoftwareUpdatable
* **\<your-namespace\>:\<gateway-name\>:edge:containers** - a virtual one that provides the edge container management _Features_, e.g.
    * ConnectionStatus
    * ContainerFactory
    * SoftwareUpdatable

### What's next

{{% refn "https://docs.bosch-iot-suite.com/device-management/How-to-guides.html" %}}Learn how to manage your device via Bosch
IoT Device Management{{% /refn %}}