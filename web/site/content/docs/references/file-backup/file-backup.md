---
title: "File Backup API"
type: docs
description: >
  The file backup service provides the ability to backup and restore to and from a backend storage.
weight: 4
---

## **Backup**
Create a backup to backend storage.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//backup`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/backup` | - |
> | path | `/features/BackupAndRestore/inbox/messages/backup` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | UUID | - |
> | **Value** | | |
> | correlationID | UUID | Identifier of the backup file |
> | providers | - | - |
> | options | - | - |
<br>

**Example** : In this example, the User can create the backup:

**Topic:** `command//edge:device/req//backup`
```json
{
	"topic":"edge/device/things/live/messages/backup",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/BackupAndRestore/inbox/messages/backup",
	"value":{
		"correlationID":"upload-id-1704439450#n",
		"providers":{},
		"options":{
			"backup.dir":"/var/tmp/backup",
			"https.url":""
		}
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//backup`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/backup`

**Ditto Path** : `/features/BackupAndRestore/outbox/messages/backup`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation backup`

**Example** :

**Topic:** `command//edge:device/res//backup``
```json
{
	"topic":"edge/device/things/live/messages/backup",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/BackupAndRestore/outbox/messages/backup",
	"status":204
}
```
</details>

## **Restore**
Restore the backup from backend.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//restore`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/restore` | - |
> | path | `/features/BackupAndRestore/inbox/messages/restore` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | UUID | - |
> | **Value** | | |
> | correlationID | UUID | - |
> | providers | - | - |
> | options | - | - |
<br>

**Example** : In this example, the User can restore from backend:

**Topic:** `command//edge:device/req//restore`
```json
{
	"topic":"edge/device/things/live/messages/restore",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/BackupAndRestore/inbox/messages/restore",
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

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//restore`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/restore`

**Ditto Path** : `/features/BackupAndRestore/outbox/messages/restore`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the restore operation`

**Example** :

**Topic:** `command//edge:device/res//restore``
```json
{
	"topic":"edge/device/things/live/messages/restore",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/BackupAndRestore/outbox/messages/restore",
	"status": 204
}
```
</details>