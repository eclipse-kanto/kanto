---
title: "Update Manager API"
type: docs
description: >
  The kanto update manager service provides orchestration of OTA Updates towards a target device in a smart way.
weight: 5
---

## **Apply**
Applies a desired state to the device.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//apply`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/apply` | Information about the affected Thing and the type of operation |
> | path | `/features/UpdateManager/inbox/messages/apply` | A path to the `UpdateManager` Feature, it's message channel, and `apply` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | |
> | activityId | | The activity id of the new desired state |
> | ***desiredState*** | | The desired state to be applied on a device |
> | **baselines** | | An array of domain or cross-domain dependencies between components |
> | title | | The title of the dependency |
> | description | | The description of the dependency |
> | preconditions | | The preconditions of the dependency |
> | components | | An array of the components of the dependency |
> | **domains** | | An array of desired state for a single domain |
> | id | | The id of this domain|
> | ***config*** | | An array of key/value string pair|
> | key | | The key string |
> | value | | The value of the key string |
> | ***components*** | | An array of desired state component with additional key-value configuration pairs |
> | id | | The id of the component |
> | version | | The version of the component |
> | key | | The key string |
> | value | | The value of the key string |

<br>

**Example** : Apply a desired state to the device.

**Topic:** `command//edge:device/req//apply`
```json
{
	"topic":"edge/device/things/live/messages/apply",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/UpdateManager/inbox/messages/apply",
	"value":{
		"activityId": "d91ad6fe-9b0c-4549-bf31-17d0a71b61de",
		"desiredState": {
			"baselines": [
				{
					"title": "simple-baseline",
					"description": "",
					"precondition": "",
					"components": [
						"domain:component1",
						"domain:component2"
					]
				}
			],
			"domains": [
				{
					"id": "containers",
					"config": [
						{
							"key": "source",
							"value": "value"
						}
					],
					"components": [
						{
							"id": "containers:influxdb",
							"version": "2.7.1",
							"config": [
								{
									"key": "image",
									"value": "docker.io/library/influxdb:$influxdb_version"
								}
							]
						}
					]
				}
			]
		}
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//apply`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/apply` | Information about the affected Thing and the type of operation |
> | path | `/features/UpdateManager/outbox/messages/apply` | A path to the `UpdateManager` Feature, it's message channel, and `apply` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the request message |
> | **Status** | | Status of the operation apply |

<br>


**Example** : Successful response of an `apply` desired state operation.

**Topic:** `command//edge:device/res//apply`
```json
{
	"topic": "edge/device/things/live/messages/apply",
	"headers": {
		"content-type": "application/json",
		"correlation-id": "<UUID>"
	},
	"path": "/features/UpdateManager/outbox/messages/apply",
	"status": 204
}
```
</details>

## **Refresh**
Reads the current state from the device and updates the status of the `UpdateManager` feature.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//refresh`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/refresh` | Information about the affected Thing and the type of operation |
> | path | `/features/UpdateManager/inbox/messages/refresh` | A path to the `UpdateManager` Feature, it's message channel, and `refresh` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | |
> | activityId | | The activity id of the `refresh` operation |

<br>

**Example** : Update the status of the `UpdateManager` feature.

**Topic:** `command//edge:device/req//refresh`
```json
{
	"topic": "edge/device/things/live/messages/refresh",
	"headers": {
		"response-required": true,
		"content-type":" application/json",
		"correlation-id": "<UUID>"
	},
	"path": "/features/UpdateManager/inbox/messages/refresh",
	"value": {
		"activityId": "e08b071c-c19e-41de-8da0-e2843113161f"
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//refresh`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/refresh` | Information about the affected Thing and the type of operation |
> | path | `/features/UpdateManager/outbox/messages/refresh` | A path to the `UpdateManager` Feature, it's message channel, and `refresh` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the request message |
> | **Status** | | Status of the `refresh` operation |

<br>

**Example** : Successful response of a `refresh` operation.

**Topic:** `command//edge:device/res//refresh`
```json
{
	"topic": "edge/device/things/live/messages/refresh",
	"headers": {
		"content-type": "application/json",
		"correlation-id": "<UUID>"
	},
	"path": "/features/UpdateManager/outbox/messages/refresh",
	"status": 204
}
```
</details>