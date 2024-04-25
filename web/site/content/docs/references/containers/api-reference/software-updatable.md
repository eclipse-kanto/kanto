---
title: "Software Updatable API"
type: docs
description: >
  The software updatable service provides the ability to install given list of software modules and to remove modules from containers.
weight: 4
---

## **Install**
Install given list of software modules to the container.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//install`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/install` | - |
> | path | `/features/SoftwareUpdatable/inbox/messages/install` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | correlationId | - | different from container UUID |
> | softwareModules | - | Array from modules, which will be installed |
> | weight | - | - |
> | metadata | - | Metadata |
> | forced | true/false | Forced to install the software modules |
<br>

**Example** : In this example, the User can install listed modules:

**Topic:** `command//edge:device/req//install`
```json
{
	"topic":"edge/device/things/live/messages/install",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/SoftwareUpdatable/inbox/messages/install",
	"value":{
		"correlationId":"other_correlation_id",
		"forced":true,
		"softwareModules":[
			{
				"artifacts":[
					{
						"filename":"valid.json",
						"download":{
							"HTTPS":{
								"url":"https://raw.githubusercontent.com/eclipse-kanto/container-management/main/containerm/pkg/testutil/config/container/valid.json",
								"md5url":"https://raw.githubusercontent.com/eclipse-kanto/container-management/main/containerm/pkg/testutil/config/container/valid.json"
							}
						},
						"checksums":{
							"MD5":"8c5a0fa2c01e218262d672bf643652fd",
							"SHA1":"7539b451d818d94bcd97d401a5467b3e1c0b8981",
							"SHA256":"be8f5def8e6a61caab078be0995826ae65f5993b1a35c18ed6045c3db37c4a3a"
						},
						"size":100
					}
				],
				"softwareModule":{
					"name":"influxdb",
					"version":"1.8.4"
				}
			}
		]
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//install`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/install`

**Ditto Path** : `/features/SoftwareUpdatable/outbox/messages/install`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | content-type | application/json | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation install software modules over the container`

**Example** :

**Topic:** `command//edge:device/res//install``
```json
{
	"topic":"edge/device/things/live/messages/install",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/SoftwareUpdatable/outbox/messages/install",
	"status": 204
}
```
</details>

## **Remove**
Remove software modules from the container.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//remove`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/remove` | - |
> | path | `/features/SoftwareUpdatable/inbox/messages/remove` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | container UUID | - |
> | **Value** | | |
> | value | - | Json presentation of software module to be removed |
<br>

**Example** : In this example, the User can remove an existing software modules container:

**Topic:** `command//edge:device/req//remove`
```json
{
	"topic":"edge/device/things/live/messages/remove",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/SoftwareUpdatable/inbox/messages/remove",
	"value": {
		"correlationId":"other_correlation_id",
		"forced":true,
		"software":[
			{
				"name":"influxdb",
				"version":""
			}
		]
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//remove`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/remove`

**Ditto Path** : `/features/SoftwareUpdatable/outbox/messages/remove`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation remove software modules from container`

**Example** :

**Topic:** `command//edge:device/res//remove``
```json
{
	"topic":"edge/device/things/live/messages/remove",
	"headers":{
		"correlation-id":"<UUID>"
	},
	"path":"/features/SoftwareUpdatable/outbox/messages/remove",
	"status":204
}
```
</details>