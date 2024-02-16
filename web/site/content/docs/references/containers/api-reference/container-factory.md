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

| Name | Value | Description |
| - | - | - | 
| topic | `<name>/<namespace>/things/live/messages/create` | - |
| path | `/features/ContainerFactory/inbox/messages/create` | - |
| **Headers** | | |
| response-required | `true` | - |
| content-type | `application/json` | - |
| correlation-id  | container UUID | - |
| **Value** | | |
| imageRef | URL | Container image URL |
| start | true/false | - |
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
  <summary>Request</summary>
### Create - Response

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//create`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/create`

**Ditto Path** : `/features/ContainerFactory/outbox/messages/create`

#### Headers

> | name      |  value     | description  |
> |-----------|-----------|-------------------------|
> | content-type      |  application/json |  -  |
> | correlation-id      |  \<UUID\> |  -  |

#### Value: `UUID of the container`

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