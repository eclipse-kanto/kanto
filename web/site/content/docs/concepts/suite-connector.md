---
title: "Suite connector"
type: docs
description: >
  Empower the edge device with a remote connectivity.
weight: 1
---

Suite connector enables the remote connectivity to an IoT cloud ecosystem of choice, powered by Eclipse Hono™ (e.g. <a href="https://www.eclipse.org/packages/packages/cloud2edge" target="_blank">Eclipse Cloud2Edge</a>. and Bosch IoT Suite). It provides the following use cases:

* **Enriched remote connection**
  * **Optimized** - to pass the messages via a single underlying connection
  * **Secured** - to protect the edge identity and data via TLS with basic and certificate-based authentication
  * **Maintained** - with a reconnect exponential backoff algorithm
  * **Synchronized** - on a connectivity recovering via a message buffering
* **Application protection** - suite connector is the only one component with a remote connectivity i.e. all local applications are protected from exposure to the public network
* **Offline mode** - local applications don't need to care about the status of the remote connection, they can stay fully operable in offline mode

![Suite connector](/kanto/images/docs/concepts/suite-connector.png)

## How it works

The suite connector plays a key role in two communication aspects - local and remote.

### Cloud connectivity

To initiate its connection, the edge has to be manually or automatically provisioned. The result of this operation is different parameters and identifiers. Currently, suite connector supports MQTT transport as a connection-oriented and requiring less resources in comparison to AMQP. Once established, the connection is used as a channel to pass the edge telemetry and event messages. The IoT cloud can control the edge via commands and responses.

In case of a connection interruption, the suite connector will switch to offline mode. The message buffer mechanism will be activated to ensure that there is no data loss. Reconnect exponential backoff algorithm will be started to guarantee that no excessive load will be generated to the IoT cloud. All local applications are not affected and can continue to operate as normal. Once the remote connection is restored, all buffered messages will be sent and the edge will be fully restored to online mode.

### Local communication

Ensuring that local applications are loosely coupled, Eclipse Hono™ MQTT definitions are in use. The event-driven local messages exchange is done via a MQTT message broker - Eclipse Mosquitto™. The suite connector takes the responsibility to forward these messages to the IoT cloud and vice versa.

The provisioning information used to establish the remote communication is available locally both on request via a predefined message and on update populated via an announcement. Applications that would like to extend the edge functionality can further use it in Eclipse Hono™ and Eclipse Ditto™ definitions.

Monitoring of the remote connection status is also enabled locally as well, along with details like the last known state of the connection, timestamp and a predefined connect/disconnect reason.
