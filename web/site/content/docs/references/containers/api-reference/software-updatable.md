---
title: "Software Updatable API"
type: docs
description: >
  The software updatable service utilizes the Eclipse hawkBit message format to install a specified list of containers (software modules) and remove already installed modules.
weight: 4
---

## **Install**
You can install a specified list of containers (software modules).

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//install`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>:edge:containers/things/live/messages/install` | Information about the affected Thing and the type of operation |
> | path | `/features/SoftwareUpdatable/inbox/messages/install` | A path to the `SoftwareUpdatable` Feature, it's message channel, and `install` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | |
> | correlationId | | Unique identifier that is used to associate and track the series of messages |
> | weight | | The weight is the priority in case of multiple, parallel instructions |
> | metadata | | The metadata is any other information which should be passed to the device |
> | forced | true/false | Forced to install the software modules |
> | **softwareModules** | | An array of modules that will be installed |
> | metadata | | The metadata is any other information which should be passed to the device |
> | **softwareModule** | | An unique identifier for the software module |
> | name | | The name of the software module |
> | version | | The version of the software module |
> | **artifacts** | | An array of artifacts contained in the software module |
> | filename | | The file name of the artifact behind the provided URLs |
> | size | | The size of the file in bytes |
> | **download** | | A map with protocols and links for artifact download |
> | key | HTTP/HTTPS/FTP/SFTP | Available transport protocols |
> | url | | URL to download the artifact |
> | md5url | | MD5URL to download the MD5SUM file |
> | **checksums** | | A map with checksums to verify the proper download |
> | MD5 | | MD5 checksum of the downloaded file |
> | SHA1 | | SHA1 checksum of the downloaded file |
> | SHA256 | | SHA256 checksum of the downloaded file |

<br>

**Example** : In this example, you can install the listed modules.

**Topic:** `command//edge:device:edge:containers/req//install`
```json
{
	"topic":"edge/device:edge:containers/things/live/messages/install",
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
				"softwareModule":{
					"name":"influxdb",
					"version":"1.8.4"
				},
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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>:edge:containers/things/live/messages/install` | Information about the affected Thing and the type of operation |
> | path | `/features/SoftwareUpdatable/outbox/messages/install` | A path to the `SoftwareUpdatable` Feature, it's message channel, and `install` command |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the `install` operation`|

<br>

**Example** : Response of a successful install of the software modules.

**Topic:** `command//edge:device/res//install``
```json
{
	"topic":"edge/device:edge:containers/things/live/messages/install",
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
Remove of an installed software module.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//remove`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>:edge:containers/things/live/messages/remove` | Information about the affected Thing and the type of operation |
> | path | `/features/SoftwareUpdatable/inbox/messages/remove` | A path to the `SoftwareUpdatable` Feature, it's message channel, and `remove` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | Json presentation of the software module to be removed |
> | correlationId | | Unique identifier that is used to associate and track the series of messages |
> | weight | | The weight is the priority in case of multiple, parallel instructions |
> | metadata | | The metadata is any other information which should be passed to the device |
> | forced | true/false | Force remove the software modules |
> | **software** | | An array of software modules to be removed |
> | group | | An identifier which groups the dependency into a certain category |
> | name | | The dependency name |
> | version | | The dependency version |
> | type | | The "category" classifier of the dependency |

<br>

**Example** : In this example, you can remove an existing software modules.

**Topic:** `command//edge:device:edge:containers/req//remove`
```json
{
	"topic":"edge/device/things:edge:containers/live/messages/remove",
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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>:edge:containers/things/live/messages/remove` | Information about the affected Thing and the type of operation |
> | path | `/features/SoftwareUpdatable/outbox/messages/remove` | A path to the `SoftwareUpdatable` Feature, it's message channel, and `remove` command |
> | **Headers** | | Additional headers |
> | correlation-id | container UUID | The container UUID |
> | **Status** | | Status of the operation remove software modules from container |

<br>

**Example** : The response of successful removal of software modules.

**Topic:** `command//edge:device:edge:containers/res//remove``
```json
{
	"topic":"edge/device:edge:containers/things/live/messages/remove",
	"headers":{
		"correlation-id":"<UUID>"
	},
	"path":"/features/SoftwareUpdatable/outbox/messages/remove",
	"status":204
}
```
</details>