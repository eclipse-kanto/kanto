---
title: "AWS Connector"
type: docs
description: >
  Empower the edge device with a remote connectivity.
weight: 1
---

AWS Connector enables the remote connectivity to an AWS IoT cloud ecosystem. It provides the following use cases:

* **Enriched remote connection**
  * **Optimized** - to pass the messages via a single underlying connection
  * **Secured** - to protect the edge identity and data via TLS with basic and certificate-based authentication
  * **Maintained** - with a reconnect exponential backoff algorithm
  * **Synchronized** - on a connectivity recovering via a message buffering
* **Application protection** - AWS Connector is the only one component with a remote connectivity i.e. all local applications are protected from exposure to the public network
* **Offline mode** - local applications don't need to care about the status of the remote connection, they can stay fully operable in offline mode
* **Device Shadow** - messages sent to the [Twin Channel](https://eclipse.dev/ditto/protocol-twinlive.html#twin) are converted to messages more suitable for [AWS Device Shadow](https://docs.aws.amazon.com/iot/latest/developerguide/device-shadow-document.html) service and sent to it.

![AWS Connector](/kanto/images/docs/concepts/aws-connector.png)

### How it works

The AWS Connector plays a key role in two communication aspects - local and remote.

#### Cloud connectivity

To initiate its connection, the edge has to be manually or automatically provisioned. The result of this operation is different parameters and identifiers. Currently, AWS Connector supports MQTT transport as a connection-oriented and requiring less resources in comparison to AMQP. Once established, the connection is used as a channel to pass the edge telemetry and event messages. The IoT cloud can control the edge via commands and responses.

In case of a connection interruption, the AWS Connector will switch to offline mode. The message buffer mechanism will be activated to ensure that there is no data loss. Reconnect exponential backoff algorithm will be started to guarantee that no excessive load will be generated to the IoT cloud. All local applications are not affected and can continue to operate as normal. Once the remote connection is restored, all buffered messages will be sent and the edge will be fully restored to online mode.

#### Local communication

Ensuring that local applications are loosely coupled, Eclipse Hono™ MQTT definitions are in use. The event-driven local messages exchange is done via a MQTT message broker - Eclipse Mosquitto™. The AWS Connector takes the responsibility to forward these messages to the IoT cloud and vice versa.

Monitoring of the remote connection status is also enabled locally as well, along with details like the last known state of the connection, timestamp and a predefined connect/disconnect reason.
