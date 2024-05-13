---
title: "File Backup API"
type: docs
description: >
  The file backup service allows you to backup and restore data to and from a backend storage.
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
> | topic | `<name>/<namespace>/things/live/messages/backup` | Information about the affected Thing and the type of operation |
> | path | `/features/BackupAndRestore/inbox/messages/backup` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | |
> | correlationID | UUID | Identifier of the backup file |
> | providers | | The providers of the restore command |
> | ***options*** | | |
> | backup.dir | | A local directory, which to be backup and upload it, using HTTP upload or Azure/AWS temporary credentials |
> | https.url | | The URL for restore the backup directory |

<br>

**Example** : In this example, you can create the backup.

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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/backup` | Information about the affected Thing and the type of operation |
> | path | `/features/BackupAndRestore/outbox/messages/backup` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the operation backup |

<br>

**Example** : The response of the backup operation.

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
> | topic | `<name>/<namespace>/things/live/messages/restore` | Information about the affected Thing and the type of operation |
> | path | `/features/BackupAndRestore/inbox/messages/restore` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | |
> | correlationID | other UUID | Identifier of the restore file |
> | providers | | The providers of the restore command |
> | ***options*** | | Options are specific for each provider |
> | backup.dir | | A local directory, which to be backup and upload it, using HTTP upload or Azure/AWS temporary credentials |
> | https.url | | The URL for restore the backup directory |

<br>

**Example** : In this example, you can restore from backend.

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
		"providers":{},
		"options":{
			"backup.dir":"/var/tmp/backup",
			"https.url":"https://raw.githubusercontent.com/eclipse-kanto/container-management/main/containerm/pkg/testutil/config/"
		}
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//restore`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/restore` | Information about the affected Thing and the type of operation |
> | path | `/features/BackupAndRestore/outbox/messages/restore` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the operation restore |

<br>

**Example** : The response of the restore operation.

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