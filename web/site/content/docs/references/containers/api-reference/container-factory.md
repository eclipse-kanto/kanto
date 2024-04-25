---
title: "Container Factory API"
type: docs
description: >
  The container factory service provides the ability to create new containers form a container image, or from a [container configuration](../container-config.md).
weight: 2
---

## **Create**
Create a container from a single container image reference with an option to start it.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//create`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/create` | - |
> | path | `/features/ContainerFactory/inbox/messages/create` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | imageRef | URL | Container image URL |
> | start | true/false | - |
<br>

**Example** : In this example, the User can create and automatically start a new `Hello World` container:

**Topic:** `command//edge:device/req//create`
```json
{
	"topic":"edge/device/things/live/messages/create",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/ContainerFactory/inbox/messages/create",
	"value":{
		"imageRef":"docker.io/library/hello-world:latest",
		"start":true
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//create`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/create`

**Ditto Path** : `/features/ContainerFactory/outbox/messages/create`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Value: `UUID of the container`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/create` | - |
> | path | `/features/ContainerFactory/outbox/messages/create` | - |
> | **Headers** | | |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |
> | **Value** | | UUID of the container |
<br>

**Example** :

**Topic:** `command//edge:device/res//create``
```json
{
	"topic":"edge/device/things/live/messages/create",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/ContainerFactory/outbox/messages/create",
	"value":"<Container UUID>"
}
```
</details>

## **Create  with config**
Create a container from a given container configuration.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//createWithConfig`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/createWithConfig` | - |
> | path | `/features/ContainerFactory/inbox/messages/createWithConfig` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | config | | json presentation of the configuration|
> | imageRef | URL | Container image URL |
> | start | true/false | - |
<br>

**Example** : In this example, the User can create and automatically start a new `Hello World` container:

**Topic:** `command//edge:device/req//createWithConfig`
```json
{
	"topic":"edge/device/things/live/messages/createWithConfig",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/ContainerFactory/inbox/messages/createWithConfig",
	"value":{
		"config":{
			"extraHosts":["ctrhost:host_ip"],
			"portMappings":[
				{
					"containerPort":80,
					"hostPort":5000
				}
			]
		},
		"imageRef":"docker.io/library/hello-world:latest",
		"start":true
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//createWithConfig`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/createWithConfig`

**Ditto Path** : `/features/ContainerFactory/outbox/messages/createWithConfig`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Value: `UUID of the container`

**Example** :

**Topic:** `command//edge:device/res//createWithConfig``
```json
{
	"topic":"edge/device/things/live/messages/createWithConfig",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/ContainerFactory/outbox/messages/createWithConfig",
	"value":"<Container UUID>"
}
```
</details>