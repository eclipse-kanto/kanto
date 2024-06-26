---
title: "Metrics API"
type: docs
description: >
  With the metrics service, you can request and receive metrics data for specific containers.
weight: 4
---

## **Request**
Request to receive data from the container.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//request`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/request` | Information about the affected Thing and the type of operation |
> | path | `/features/Metrics/inbox/messages/request` | A path to the `Metrics` Feature, it's message channel, and `request` command|
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | |
> | frequency | | Time interval of how often the metrics data will be published as duration string (e.g. 5s) |
> | **filter** | | Filter defines the type of metric data to be reported |
> | id | | An array of identifiers whose metric data to be reported, supported are: `cpu.utilization`, `memory.utilization`, `memory.total`, `memory.used`, `io.readBytes`, `io.writeBytes`, `net.readBytes`, `net.writeBytes`, `pids` |
> | originator | | Metrics data originator |

<br>

**Example** : Request metrics data with a specified filter and frequency.

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
				"id":null,
				"originator":"SYSTEM"
			}
		],
		"frequency":"2s"
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
> | path | `/features/Metrics/outbox/messages/request` | A path to the `Metrics` Feature, it's message channel, and `request` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the `request` metrics operation |

<br>

**Example** : The response of the request metrics data operation.

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
Metrics data from a container based on the frequency specified in the request.

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//data`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/data` | Information about the affected Thing and the type of operation |
> | path | `/features/Metrics/outbox/messages/data` | A path to the `Metrics` Feature and it's message channel. |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | **Value** | | The value of the received data from the container in json format |
> | timestamp | | The timestamp in ms when this measure data is published |
> | **shapshot** | | All the measurements collected per originator |
> | originator | | The originator for whose metric data to be reported |
> | **measurements** | | An array of measurements identifier and value for originator |
> | id | | The identifier whose metric data to be reported, supported are: `cpu.utilization`, `memory.utilization`, `memory.total`, `memory.used`, `io.readBytes`, `io.writeBytes`, `net.readBytes`, `net.writeBytes`, `pids` |
> | value | | The measured value per metric ID |

<br>

**Example** : Metrics data from a container.

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
				"originator":"Container:test",
				"measurements":[
					{
						"id":"memory.total",
						"value":10371616768
					},
					{
						"id":"memory.used",
						"value":1396736
					},
					{
						"id":"memory.utilization",
						"value":0.01346690714903206
					},
					{
						"id":"net.readBytes",
						"value":180
					},
					{
						"id":"net.writeBytes",
						"value":0
					},
					{
						"id":"pids",
						"value":6
					}
				]
			},
			{
				"originator":"Container:test2",
				"measurements":[
					{
						"id":"cpu.utilization",
						"value":8.751566666666667
					},
					{
						"id":"memory.total",
						"value":10371616768
					},
					{
						"id":"memory.used",
						"value":4759552
					},
					{
						"id":"memory.utilization",
						"value":0.04589016453717083
					},
					{
						"id":"io.readBytes",
						"value":0
					},
					{
						"id":"io.writeBytes",
						"value":4096
					},
					{
						"id":"net.readBytes",
						"value":610
					},
					{
						"id":"net.writeBytes",
						"value":202
					},
					{
						"id":"pids",
						"value":14
					}
				]
			}
		],
		"timestamp":1234567890
	}
}
```
</details>
