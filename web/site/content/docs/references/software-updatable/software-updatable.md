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
> | topic | `<name>/<namespace>/things/live/messages/install` | Information about the affected Thing and the type of operation |
> | path | `/features/SoftwareUpdatable/inbox/messages/install` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | Json presentation of software module to be installed |
> | correlationId | | Different identifier from the container UUID |
> | weight | | The weight is the priority in case of multiple, parallel instructions |
> | metadata | | The metadata is any other information which should be passed to the device |
> | forced | true/false | Forced to remove the software modules |
> | ***softwareModules*** | | An array of software modules to be installed |
> | metadata | | The metadata is any other information which should be passed to the device |
> | **softwareModule** | | An unique identifier for the software module |
> | name | | The name of the software module |
> | version | | The version of the software module |
> | **artifacts** | | An array of artifacts contained in the software module |
> | filename | | The file name of the artifact behind the provided URLs |
> | size | | The size of the file in bytes |
> | **download** | | A map with protocols and links to downloaded |
> | key | HTTP/HTTPS/FTP/SFTP | Available transport protocols |
> | url | | URL to download the artifact |
> | md5url | | MD5URL to download the MD5SUM file |
> | **checksums** | | A map with checksums to verify the proper download |
> | MD5 | | The checksum by md5 hash algorithm |
> | SHA1 | | The checksum by sha1 hash algorithm |
> | SHA256 | | The checksum by sha256 hash algorithm |

<br>

**Example** : In this example, you can install listed modules.

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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/install` | Information about the affected Thing and the type of operation |
> | path | `/features/SoftwareUpdatable/outbox/messages/install` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Status** | | Status of the operation install software modules |

<br>

**Example** : The response of the install software modules operation.

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
> | topic | `<name>/<namespace>/things/live/messages/download` | Information about the affected Thing and the type of operation |
> | path | `/features/SoftwareUpdatable/inbox/messages/download` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | Json presentation of software module to be download |
> | correlationId | | Different identifier from the container UUID |
> | weight | | The weight is the priority in case of multiple, parallel instructions |
> | metadata | | The metadata is any other information which should be passed to the device |
> | forced | true/false | Forced to remove the software modules |
> | ***softwareModules*** | | An array of software modules to be download |
> | metadata | | The metadata is any other information which should be passed to the device |
> | **softwareModule** | | An unique identifier for the software module |
> | name | | The name of the software module |
> | version | | The version of the software module |
> | **artifacts** | | An array of artifacts contained in the software module |
> | filename | | The file name of the artifact behind the provided URLs |
> | size | | The size of the file in bytes |
> | **download** | | A map with protocols and links to downloaded |
> | key | HTTP/HTTPS/FTP/SFTP | Available transport protocols |
> | url | | URL to download the artifact |
> | md5url | | MD5URL to download the MD5SUM file |
> | **checksums** | | A map with checksums to verify the proper download |
> | MD5 | | The checksum by md5 hash algorithm |
> | SHA1 | | The checksum by sha1 hash algorithm |
> | SHA256 | | The checksum by sha256 hash algorithm |

<br>

**Example** : In this example, you can download an existing software modules.

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

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/download` | Information about the affected Thing and the type of operation |
> | path | `/features/SoftwareUpdatable/outbox/messages/download` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | correlation-id | container UUID | The same correlation id as the sent request message |
> | **Status** | | Status of the operation download software modules |

<br>

**Example** : The response of the download software modules operation.

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