---
title: "Suite Bootstrapping API"
type: docs
description: >
  The suite bootstrapping service provides the ability to customize the automatic provisioning [suite bootstrapping configuration](../suite-bootstrapping-config.md).
weight: 2
---

## **Response**
Generate response for the given corelation identifier.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//response`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/response` | - |
> | path | `/features/Bootstrapping/inbox/messages/response` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | chunk | - | - |
> | hash | - | - |
> | requestId | other UUID | - |
<br>

**Example** : In this example, the User received the example response of the automated provisioning.

**Topic:** `command//edge:device/req//response`
```json
{
	"topic":"edge/device/things/live/messages/response",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Bootstrapping/inbox/messages/response",
	"value":{
		"chunk":"ewoJImNhQ2VydCI6ICIvZXRjL3N1aXRlLWNvbm5lY3Rvci9pb3RodWIuYgIjEyMzQ1NiIKfQ==",
		"hash":"8768b29d6130433673ce3a1bc9005fc03c213306dc328f1ea33d240c5a116ee1",
		"requestId":"a41a0d55-dba6-4e8b-b1ad-a3c414fff42b"
	},
	"timestamp":1234567890
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//response`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/response`

**Ditto Path** : `/features/Bootstrapping/outbox/messages/response`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Value: `other UUID `

**Example** :

**Topic:** `command//edge:device/res//create``
```json
{
	"topic":"edge/device/things/live/messages/response",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Bootstrapping/outbox/messages/response",
	"value":"<other UUID>"
}
```
</details>
