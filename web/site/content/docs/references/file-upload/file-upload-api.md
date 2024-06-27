---
title: "File Upload API"
type: docs
description: >
  The file upload service provides the ability to start, trigger, cancel, activate or deactivate upload of the file.
weight: 4
---

## **Start**
Start a file upload operation.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//start`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/start` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/inbox/messages/start` | A path to the `AutoUploadable` Feature, it's message channel, and `start` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | |
> | correlationID | other UUID | Identifier of the uploaded file |
> | **options** | | Options are specific for each provider |
> | storage.provider | aws/azure/generic | Storage provider that will be used for uploading the files |

<br>

**Example** : Start uploading a file.

**Topic:** `command//edge:device/req//start`
```json
{
	"topic":"edge/device/things/live/messages/start",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/inbox/messages/start",
	"value":{
		"correlationID":"upload-id-1704439450#n",
		"options":{
			"aws.access.key.id":"AWSACCESSKEYID",
			"aws.region":"eu-central-1",
			"aws.s3.bucket":"blob-upload-test",
			"aws.secret.access.key":"AWSSECRETACCESSKEY",
			"storage.provider":"aws"
		}
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>/res//start`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/start` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/outbox/messages/start` | A path to the `AutoUploadable` Feature, it's message channel, and `start` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the operation start upload the file |

<br>

**Example** : Successful response of a `start` file upload operation.

**Topic:** `command//edge:device/res//start``
```json
{
	"topic":"edge/device/things/live/messages/start",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/outbox/messages/start",
	"status": 204
}
```
</details>

## **Trigger**
Trigger a file upload operation.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>/req//trigger`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/trigger` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/inbox/messages/trigger` | A path to the `AutoUploadable` Feature, it's message channel, and `trigger` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | |
> | correlationID | other UUID | Identifier of the triggered file |
> | **options** | | Options are specific for each provider |

<br>

**Example** : Trigger a file upload operation.

**Topic:** `command//edge:device/req//trigger`
```json
{
	"topic":"edge/device/things/live/messages/trigger",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/inbox/messages/trigger",
	"value":{
		"correlationID":"other <UUID>",
		"options":{}
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>/res//trigger`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/trigger` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/outbox/messages/trigger` | A path to the `AutoUploadable` Feature, it's message channel, and `trigger` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the `trigger` operation |

<br>

**Example** : Successful response of a `trigger` operation.

**Topic:** `command//edge:device/res//trigger``
```json
{
	"topic":"edge/device/things/live/messages/trigger",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/outbox/messages/trigger",
	"status":204
}
```
</details>

## **Cancel**
Cancel a file upload operation.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//cancel`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/cancel` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/inbox/messages/cancel` | A path to the `AutoUploadable` Feature, it's message channel, and `cancel` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | |
> | correlationID | other UUID | Identifier of the uploaded file |
> | statusCode | | This status code is set to update status code |
> | message | | This message is set to update status message |

<br>

**Example** : Cancel a file upload operation.

**Topic:** `command//edge:device/req//cancel`
```json
{
	"topic":"edge/device/things/live/messages/cancel",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/inbox/messages/cancel",
	"value":{
		"correlationID":"upload-id-1704439450#n",
		"statusCode": 400,
		"message":"description why the upload is canceled "
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>/res//cancel`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/cancel` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/outbox/messages/cancel` | A path to the `AutoUploadable` Feature, it's message channel, and `cancel` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the operation `cancel` file upload |

<br>

**Example** : The response of the cancel file upload operation.

**Topic:** `command//edge:device/res//cancel``
```json
{
	"topic":"edge/device/things/live/messages/cancel",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/outbox/messages/cancel",
	"status":204
}
```
</details>

## **Activate**
Activate an upload of a file.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//activate`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/activate` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/inbox/messages/activate` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | |
> | from | | A Time after which the upload must be activated |
> | to | | A Time grater than `from` marks the end of activated |

<br>

**Example** : Activate file upload.

**Topic:** `command//edge:device/req//activate`
```json
{
	"topic":"edge/device/things/live/messages/activate",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/inbox/messages/activate",
	"value":{
		"from":957139200,
		"to":959817599
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>/res//activate`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/activate` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/outbox/messages/activate` | A path to the `AutoUploadable` Feature, it's message channel, and `activate` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the operation `activate` file upload |

<br>

**Example** : Successful response of an activate file upload operation.

**Topic:** `command//edge:device/res//activate``
```json
{
	"topic":"edge/device/things/live/messages/activate",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/outbox/messages/activate",
	"status":204
}
```
</details>

## **Deactivate**
Deactivate a file upload.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>/req//deactivate`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/deactivate` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/inbox/messages/deactivate` | A path to the `AutoUploadable` Feature, it's message channel, and `deactivate` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the response message |
> | **Value** | | |

<br>

**Example** : Deactivate a file upload.

**Topic:** `command//edge:device/req//deactivate`
```json
{
	"topic":"edge/device/things/live/messages/deactivate",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/inbox/messages/deactivate",
	"value":{}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>0/res//deactivate`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/deactivate` | Information about the affected Thing and the type of operation |
> | path | `/features/AutoUploadable/outbox/messages/deactivate` | A path to the `AutoUploadable` Feature, it's message channel, and `deactivate` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the operation `deactivate` file upload |

<br>

**Example** : Successful response of the `deactivate` file upload operation.

**Topic:** `command//edge:device/res//deactivate``
```json
{
	"topic":"edge/device/things/live/messages/deactivate",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/AutoUploadable/outbox/messages/deactivate",
	"status":204
}
```
</details>