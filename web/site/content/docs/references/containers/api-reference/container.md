---
title: "Container API"
type: docs
description: >
  The container service provides the ability to start, pause, resume, stop, stop with specified options, rename, update or remove existing containers.
weight: 3
---

## **Start**
Start an existing container.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//start`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/start` | - |
> | path | `/features/Container:<UUID>/inbox/messages/start` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
<br>

**Example** : In this example, the User can start an existing container:

**Topic:** `command//edge:device/req//start`
```json
{
	"topic":"edge/device/things/live/messages/start",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/inbox/messages/start",
	"value":{}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//start`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/start`

**Ditto Path** : `/features/Container:<UUID>/outbox/messages/start`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation start over the container`

**Example** :

**Topic:** `command//edge:device/res//start``
```json
{
	"topic":"edge/device/things/live/messages/start",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/outbox/messages/start",
	"status": 204
}
```
</details>

## **Stop**
Stop an existing and running container.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//stop`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/stop` | - |
> | path | `/features/Container:<UUID>/inbox/messages/stop` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
<br>

**Example** : In this example, the User can stop an existing and running container:

**Topic:** `command//edge:device/req//stop`
```json
{
	"topic":"edge/device/things/live/messages/stop",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/inbox/messages/stop",
	"value":{}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//stop`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/stop`

**Ditto Path** : `/features/Container:<UUID>/outbox/messages/stop`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation stop over the container`

**Example** :

**Topic:** `command//edge:device/res//stop``
```json
{
	"topic":"edge/device/things/live/messages/stop",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/outbox/messages/stop",
	"status":204
}
```
</details>

## **Stop with options**
Stop an existing and running container with given options.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//stopWithOptions`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/stopWithOptions` | - |
> | path | `/features/Container:<UUID>/inbox/messages/stopWithOptions` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | signal | `SIGTERM` | Stop a container using a specific signal. Signals could be specified by using their names or numbers, e.g. SIGINT or 2 |
> | timeout | -1 << 63 // -9223372036854775808 | Sets the timeout period in seconds to gracefully stop the container. When timeout expires the container process would be forcibly killed |
> | name | - | Stop a container with a specific name |
> | force | true/false | Whether to send a SIGKILL signal to the container's process if it does not finish within the timeout specified |
<br>

**Example** : In this example, the User can stop an existing and running container with specified options:

**Topic:** `command//edge:device/req//stopWithOptions`
```json
{
	"topic":"edge/device/things/live/messages/stopWithOptions",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/inbox/messages/stopWithOptions",
	"value":{
		"signal":"SIGINT"
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//stopWithOptions`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/stopWithOptions`

**Ditto Path** : `/features/Container:<UUID>/outbox/messages/stopWithOptions`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation stop with options over the container`

**Example** :

**Topic:** `command//edge:device/res//stopWithOptions``
```json
{
	"topic":"edge/device/things/live/messages/stopWithOptions",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/outbox/messages/stopWithOptions",
	"status":204
}
```
</details>

## **Rename**
Rename an existing container with given new name.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//rename`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/rename` | - |
> | path | `/features/Container:<UUID>/inbox/messages/rename` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | value | container new name | - |
<br>

**Example** : In this example, the User can rename an existing container with specified name:

**Topic:** `command//edge:device/req//rename`
```json
{
	"topic":"edge/device/things/live/messages/rename",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/inbox/messages/rename",
	"value":"new_container_name"
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//rename`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/rename`

**Ditto Path** : `/features/Container:<UUID>/outbox/messages/rename`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation rename container`

**Example** :

**Topic:** `command//edge:device/res//rename``
```json
{
	"topic":"edge/device/things/live/messages/rename",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/outbox/messages/rename",
	"status":204
}
```
</details>

## **Update**
Update an existing container without recreating it. The provided configurations will be merged with the current one.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//update`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/update` | - |
> | path | `/features/Container:<UUID>/inbox/messages/update` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | restartPolicy | - | Updates the restart policy for the container. The policy will be applied when the container exits |
> | type | no/always/unless-stopped/on-failure | Supported restart policies |
> | timeout | -1 << 63 // -9223372036854775808 | Updates the time out period in seconds for each retry that will be made to restart the container on exit if the policy is set to on-failure |
> | maxRetryCount | -1 << 31 // -2147483648 | Updates the number of retries that will be made to restart the container on exit if the policy is on-failure |

<br>

**Example** : In this example, the User can update an existing container with specified options:

**Topic:** `command//edge:device/req//update`
```json
{
	"topic":"edge/device/things/live/messages/update",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/inbox/messages/update",
	"value":{
		"restartPolicy":{
			"type":"always"
		}
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//update`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/update`

**Ditto Path** : `/features/Container:<UUID>/outbox/messages/update`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation update over the container`

**Example** :

**Topic:** `command//edge:device/res//update``
```json
{
	"topic":"edge/device/things/live/messages/update",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/outbox/messages/update",
	"status":204
}
```
</details>

## **Remove**
Remove a container and frees the associated resources.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//remove`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/remove` | - |
> | path | `/features/Container:<UUID>/inbox/messages/remove` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | value | true/false | Force stopping before removing a container |
<br>

**Example** : In this example, the User can remove an existing container:

**Topic:** `command//edge:device/req//remove`
```json
{
	"topic":"edge/device/things/live/messages/remove",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/inbox/messages/remove",
	"value":true
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//remove`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/remove`

**Ditto Path** : `/features/Container:<UUID>/outbox/messages/remove`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation remove container`

**Example** :

**Topic:** `command//edge:device/res//remove``
```json
{
	"topic":"edge/device/things/live/messages/remove",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/Container:<UUID>/outbox/messages/remove",
	"status":204
}
```
</details>