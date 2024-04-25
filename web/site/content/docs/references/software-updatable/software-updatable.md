---
title: "Software Updatable API"
type: docs
description: >
  The software updatable service provides the ability to install given list of software modules and to download modules.
weight: 4
---

## **Install**
Install given list of software modules.

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
> | correlation-id  | UUID | - |
> | **Value** | | |
> | correlationId | - | other UUID |
> | softwareModules | - | Array from modules, which will be installed |
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
		"softwareModules":[
			{
				"softwareModule":{
					"name":"install-hello",
					"version":"1.0.0"
				},
				"artifacts":[
					{
						"checksums":{
							"SHA256":"db954c633393c1402f145a60fd58d312f5af96ce49422fcfd6ce42a3c4cceeca",
							"MD5":"8c5a0fa2c01e218262d672bf643652fd",
							"SHA1":"7539b451d818d94bcd97d401a5467b3e1c0b8981"
						},
						"download":{
							"HTTPS":{
								"url":"https://github.com/eclipse-kanto/kanto/raw/main/quickstart/install_hello.sh"
							}
						},
						"filename":"install.sh",
						"size":544
					}
				]
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

#### Status: `Status of the operation install software modules`

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

## **Download**
Download software modules.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//download`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/download` | - |
> | path | `/features/SoftwareUpdatable/inbox/messages/download` | - |
> | **Headers** | | |
> | response-required | `true` | - |
> | content-type | `application/json` | - |
> | correlation-id  | UUID | - |
> | **Value** | | |
> | value | | Json presentation of software module to be download |
> | correlationId | - | other UUID |
> | softwareModules | - | Array from modules, which will be installed |
> | weight | - | - |
> | metadata | - | Metadata |
> | forced | true/false | Forced to install the software modules |

<br>

**Example** : In this example, the User can download an existing software modules:

**Topic:** `command//edge:device/req//download`
```json
{
	"topic":"edge/device/things/live/messages/download",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/SoftwareUpdatable/inbox/messages/download",
	"value": {
		"correlationId":"other_correlation_id",
        "softwareModules":[
			{
				"softwareModule":{
					"name":"install-hello",
					"version":"1.0.0"
				},
				"artifacts":[
					{
						"checksums":{
							"SHA256":"db954c633393c1402f145a60fd58d312f5af96ce49422fcfd6ce42a3c4cceeca",
							"MD5":"8c5a0fa2c01e218262d672bf643652fd",
							"SHA1":"7539b451d818d94bcd97d401a5467b3e1c0b8981"
						},
						"download":{
							"HTTPS":{
								"url":"https://github.com/eclipse-kanto/kanto/raw/main/quickstart/install_hello.sh"
							}
						},
						"filename":"install.sh",
						"size":544
					}
				]
			}
		],
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//download`

**Ditto Topic** : `<name>/<namespace>/things/live/messages/download`

**Ditto Path** : `/features/SoftwareUpdatable/outbox/messages/download`

#### Headers

> | Name | Value | Description |
> | - | - | - |
> | correlation-id | \<UUID\> | - |

#### Status: `Status of the operation download software modules`

**Example** :

**Topic:** `command//edge:device/res//download``
```json
{
	"topic":"edge/device/things/live/messages/download",
	"headers":{
		"correlation-id":"<UUID>"
	},
	"path":"/features/SoftwareUpdatable/outbox/messages/download",
	"status":204
}
```
</details>