---
title: "File Upload API"
type: docs
description: >
  The file upload service provides the ability to start, trigger, cancel, activate or deactivate upload of the file.
weight: 4
---

## **Start**
Start to upload file.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//start`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/start` | - |
> | path | `/features/AutoUploadable/inbox/messages/start` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | UUID | - |
> | **Value** | | |
> | correlationID  | other UUID | - |
> | options | - | - |
<br>

**Example** : In this example, the User can start uploading a file:

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

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//start`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/start`

**Ditto Path** : `/features/AutoUploadable/outbox/messages/start`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation start to upload the file`

**Example** :

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
Trigger operation is invoked from the backend.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//trigger`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/trigger` | - |
> | path | `/features/AutoUploadable/inbox/messages/trigger` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | UUID | - |
> | **Value** | | |
> | correlationID  | other UUID | - |
> | options | - | - |
<br>

**Example** : In this example, the User can pause an existing and running container:

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

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//trigger`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/trigger`

**Ditto Path** : `/features/AutoUploadable/outbox/messages/trigger`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation trigger`

**Example** :

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
Cancel upload the file.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//cancel`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/cancel` | - |
> | path | `/features/AutoUploadable/inbox/messages/cancel` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | UUID | - |
> | **Value** | | |
> | correlationID  | other UUID | - |
> | options  | - | - |
<br>

**Example** : In this example, the User can resume cancel operation upload of the file:

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

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//cancel`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/cancel`

**Ditto Path** : `/features/AutoUploadable/outbox/messages/cancel`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation cancel upload of file`

**Example** :

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
Activate upload of the file.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//activate`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/activate` | - |
> | path | `/features/AutoUploadable/inbox/messages/activate` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | UUID | - |
> | **Value** | | |
> | correlationID | other UUID | - |
> | options | - | - |
<br>

**Example** : In this example, the User can activate upload of the file:

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

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//activate`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/activate`

**Ditto Path** : `/features/AutoUploadable/outbox/messages/activate`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation activate upload of file`

**Example** :

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
Deactivate upload file.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//deactivate`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/deactivate` | - |
> | path | `/features/AutoUploadable/inbox/messages/deactivate` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | UUID | - |
> | **Value** | | |
> | correlationID | other UUID | - |
> | options | - | - |
<br>

**Example** : In this example, the User can deactivate upload of file:

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

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//deactivate`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/deactivate`

**Ditto Path** : `/features/AutoUploadable/outbox/messages/deactivate`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation deactivate`

**Example** :

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