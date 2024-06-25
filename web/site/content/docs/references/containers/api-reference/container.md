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
> | topic | `<name>/<namespace>/things/live/messages/start` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/inbox/messages/start` | A path to the `Container` Feature, it's message channel, and `start` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | |

<br>

**Example** : Start an existing container.

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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/start` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/outbox/messages/start` | A path to the `Container` Feature, it's message channel, and `start` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the request message |
> | **Status** | | Status of the operation start over the container |

<br>


**Example** : Response of a successful `start` operation.

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
> | topic | `<name>/<namespace>/things/live/messages/stop` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/inbox/messages/stop` | A path to the `Container` Feature, it's message channel, and `stop` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | |

<br>

**Example** : Stop an existing and running container.

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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/stop` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/outbox/messages/stop` | A path to the `Container` Feature, it's message channel, and `stop` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the request message |
> | **Status** | | Status of the operation stop over the container |

<br>

**Example** : Response of a successful stop operation.

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
> | topic | `<name>/<namespace>/things/live/messages/stopWithOptions` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/inbox/messages/stopWithOptions` | A path to the `Container` Feature, it's message channel, and `stopWithOptions` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | |
> | signal | `SIGTERM` | Stop a container using a specific signal. Signals could be specified by using their names or numbers, e.g. `SIGINT` or 2 |
> | timeout | -1 << 63 // -9223372036854775808 | Sets the timeout period in seconds to gracefully stop the container. When timeout expires the container process would be forcibly killed |
> | force | true/false | Whether to send a SIGKILL signal to the container's process if it does not finish within the timeout specified |

<br>

**Example** : Stop an existing and running container with specified options.

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
		"signal":"SIGINT",
		"timeout": 30,
		"force": true
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//stopWithOptions`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/stopWithOptions` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/outbox/messages/stopWithOptions` | A path to the `Container` Feature, it's message channel, and `stopWithOptions` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the request message |
> | **Status** | | Status of the operation stop with options over the container |

<br>


**Example** : Response of a successful the `stopWithOptions` operation.

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
Change the name of an existing container to the specified new name.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//rename`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/rename` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/inbox/messages/rename` | A path to the `Container` Feature, it's message channel, and `rename` command  |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | The new name of the container |

<br>

**Example** : Change the name of an existing container to the specified new name.

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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/rename` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/outbox/messages/rename` | A path to the `Container` Feature, it's message channel, and `rename` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the request message |
> | **Status** | | Status of the operation rename container |

<br>

**Example** : The response of the rename operation.

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
> | topic | `<name>/<namespace>/things/live/messages/update` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/inbox/messages/update` | A path to the `Container` Feature, it's message channel, and `update` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | |
> | **restartPolicy** | | Updates the restart policy for the container. The policy will be applied when the container exits |
> | type | no/always/unless-stopped/on-failure | The container's restart policy, the supported types are: always, no, on-failure and unless-stopped |
> | maxRetryCount | -1 << 31 // -2147483648 | Maximum number of retries that are made to restart the container on exit with fail, if the `type` is on-failure |
> | timeout | -1 << 63 // -9223372036854775808 | Timeout period in seconds for each retry that is made to restart the container on exit with fail, if the `type` is on-failure  |
> | **resources** | | |
> | memory | | Hard memory limitation of the container as a number with a unit suffix of B, K, M and G, the minimum allowed value is 3M |
> | memoryReservation | | Soft memory limitation of the container as a number with a unit suffix of B, K, M and G, if `memory` is specified, the `memoryReservation` must be smaller than it |
> | memorySwap | | Total amount of memory and swap that the container can use as a number with a unit suffix of B, K, M and G, use -1 to allow the container to use unlimited swap |

<br>

**Example** : Update an existing container resources and restart policy.

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
			"type":"on-failure",
			"maxRetryCount":3,
			"timeout":10
		},
		"resources":{
			"memory":"500M",
			"memoryReservation":"300M",
			"memorySwap":"1G",
		}
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//update`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/update` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/outbox/messages/update` | A path to the `Container` Feature, it's message channel, and `update` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the request message |
> | **Status** | | Status of the `update` operation over the container |

<br>

**Example** : Successful response of an `update` operation.

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
> | topic | `<name>/<namespace>/things/live/messages/remove` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/inbox/messages/remove` | A path to the `Container` Feature, it's message channel, and `remove` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | true/false | Force stopping before removing a container |

<br>

**Example** : Remove an existing container.

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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/remove` | Information about the affected Thing and the type of operation |
> | path | `/features/Container:<UUID>/outbox/messages/remove` | A path to the `Container` Feature, it's message channel, and `remove` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the request message |
> | **Status** | | Status of the operation remove container |

<br>

**Example** : Successful response of an `remove` operation.

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