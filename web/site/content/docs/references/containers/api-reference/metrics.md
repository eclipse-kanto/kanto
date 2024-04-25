---
title: "Metrics API"
type: docs
description: >
  The metrics service provides the ability to make requests and receive the data for some originators or containers.
weight: 5
---

## **Request**
Request to receive data from the container.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//request`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/request` | - |
> | path | `/features/Metrics/inbox/messages/request` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | filter | - | Filter defines the type of metric data to be reported |
> | frequency | - | Duration of how often the metrics data to be published |
<br>

**Example** : In this example, the User can request metrics data with some specified filter and frequency:

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
				"id":"Container",
				"originator":"test.process"
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

**Ditto Topic** : `<name>/<namespace>/things/live/messages/request`

**Ditto Path** : `/features/Metrics/outbox/messages/request`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation start over the container`

**Example** :

**Topic:** `command//edge:device/res//request``
```json
{
	"topic":"edge/device/things/live/messages/request",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/outbox/messages/request",
	"status": 204
}
```
</details>

## **Data**
Receive the data from container.

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//data`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/data`

**Ditto Path** : `/features/Metrics/outbox/messages/data`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |

#### Value: `The value of the received data from the container in json format`

**Example** :

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
