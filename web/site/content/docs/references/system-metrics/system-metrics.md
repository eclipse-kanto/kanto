---
title: "System Metrics API"
type: docs
description: >
  The system metrics service provides the ability to make requests and receive the data for some originators.
weight: 4
---

## **Request**
Request to receive data.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//request`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/request` | Information about the affected Thing and the type of operation |
> | path | `/features/Metrics/inbox/messages/request` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | |
> | frequency | | Duration of how often the metrics data to be published |
> | ***filter*** | | An array of the type of metric data to be reported |
> | id | | An array of identifiers whose metric data to be reported, supported are: cpu.utilization, memory.utilization, memory.total, memory.used, io.readBytes, io.writeBytes |
> | originator | | The originator for whose metric data to be reported |

<br>

**Example** : In this example, you can request metrics data with some specified filter and frequency.

**Topic:** `command//edge:device/req//request`
```json
{
	"topic":"edge/device/things/live/messages/request",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Metrics/inbox/messages/request",
	"value":{
		"filter":[
			{
				"id":["io.*","cpu.*","memory.*"],
				"originator":"SYSTEM"
			}
		],
		"frequency":"5s"
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//request`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/request` | Information about the affected Thing and the type of operation |
> | path | `/features/Metrics/outbox/messages/request` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the operation request the metrics data |

<br>

**Example** : The response of the request the metrics data.

**Topic:** `command//edge:device/res//request``
```json
{
	"topic":"edge/device/things/live/messages/request",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Metrics/outbox/messages/request",
	"status": 204
}
```
</details>

## **Data**
Receive the data from container.

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//data`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/data` | Information about the affected Thing and the type of operation |
> | path | `/features/Metrics/outbox/messages/data` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | **Value** | | The value of the received data from the device in json format |
> | timestamp | | The timestamp in ms when this measure data is published |
> | ***shapshot*** | | All the measurements collected at a concrete time per originator
> | originator | | The originator for whose metric data to be reported |
> | **measurements** | | An array of measurements identifier and value for originator |
> | id | | The identifier whose metric data to be reported, supported are: cpu.utilization, cpu.load1, cpu.load5, cpu.load15, memory.utilization, memory.total, memory.available, memory.used, io.readBytes, io.writeBytes |
> | value | | The measured value per metric ID |

<br>

**Example** : The received data from the device.

**Topic:** `command//edge:device/res//data``
```json
{
	"topic":"edge/device/things/live/messages/data",
	"headers":{
		"content-type":"application/json",
	},
	"path":"/features/Metrics/outbox/messages/data",
	"value":{
		"snapshot":[
			{
				"originator":"SYSTEM",
				"measurements":[
					{
						"id":"cpu.utilization",
						"value":1.1555555555484411
					},
					{
						"id":"cpu.load1",
						"value":0.17
					},
					{
						"id":"cpu.load5",
						"value":0.27
					},
					{
						"id":"cpu.load15",
						"value":0.24
					},
					{
						"id":"memory.total",
						"value":10371616768
					},
					{
						"id":"memory.available",
						"value":5281644544
					},
					{
						"id":"memory.used",
						"value":4563206144
					},
					{
						"id":"memory.utilization",
						"value":43.99705702662538
					}
				]
			}
		],
		"timestamp":1234567890
	}
}
```
</details>
