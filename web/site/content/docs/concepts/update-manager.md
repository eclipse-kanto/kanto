---
title: "Update manager"
type: docs
description: >
  Empower the edge device for OTA updates.
weight: 2
---

Update manager enables a lightweight core component which is capable to easily perform complex OTA update scenarios on a target device. The following capabilities are provided:

* **Lightweight** - consists of a single core component which orchestrates the update process
* **Flexible deployment** - supports different deployment models - natively, as an executable or container
* **Unified API** - all domain agents utilize a unified Update Agent API for interacting with the Update Manager
* **MQTT Connectivity** - the connectivity and communication between the Update Manager and domain agent is MQTT-based
* **Multi-domain integration** - easily integrates, scales and performs complex update operations across multiple domains
* **Default update agents** - integrates with the Kanto provided out-of-the box domain update agent implementation for deployment of container into the Kanto container management
* **Pluggable architecture** - provides an extensible model for plug-in custom orchestration logic
* **Asynchronous updates** - asynchronous and independent update process across the different domains
* **Multi-staged updates** - the update process is organized into different stages
* **Configurable** - offers a variety of configuration options for connectivity, supported domains, message reporting and etc

![Update manager](/kanto/images/docs/concepts/update-manager.png)

### How it works

The update process is initiated by sending the desired state specification as an MQTT message towards the device, which is handled by the Update Manager component.

The desired state specification in the scope of the Update Manager is a JSON-based document, which consists of multiple component definitions per domain, representing the desired state to be applied on the target device.
A component in the same context means a single, atomic and updatable unit, for example, OCI-compliant container, software application or firmware image.

Each domain agent is a separate and independent software component, which implements the Update Agent API for interaction with the Update Manager and manages the update logic for concrete domain. For example - container management.

The Update Manager, operating as a coordinator, is responsible for processing the desired state specification, distributing the split specification across the different domain agents, orchestrating the update process via MQTT-based commands, collecting and consolidating the feedback responses from the domain update agents, and reporting the final result of the update campaign to the backend.

As extra features and capabilities, the Update Manager enables reboot of the host after the update process is completed and reporting of the current state of the system to the backend.