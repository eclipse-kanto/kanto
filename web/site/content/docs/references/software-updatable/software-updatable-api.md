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
> | path | `/features/SoftwareUpdatable/inbox/messages/install` | A path to the `SoftwareUpdatable` Feature, it's message channel, and `install` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the response message |
> | **Value** | | JSON presentation of software module that will be installed |
> | correlationId | | Unique identifier that is used to associate and track the series of messages |
> | weight | | The weight is the priority in case of multiple, parallel instructions |
> | metadata | | The metadata is any other information which should be passed to the device |
> | forced | true/false | Forced to remove the software modules |
> | ***softwareModules*** | | An array of software modules to be installed |
> | metadata | | The metadata is any other information which should be passed to the device |
> | **softwareModule** | | An unique identifier for the software module |
> | name | | The name of the software module |
> | version | | The version of the software module |
> | **artifacts** | | An array of artifacts contained in the software module |
> | filename | | The file name of the artifact |
> | size | | Artifact file size in bytes |
> | **download** | | A map with protocols and links for downloading the artifacts |
> | key | HTTP/HTTPS/FTP/SFTP | Available transport protocols |
> | url | | URL to download the artifact |
> | md5url | | MD5URL to download the MD5SUM file |
> | **checksums** | | A map with checksums to verify the proper download |
> | MD5 | | MD5 checksum |
> | SHA1 | | MD5 checksum |
> | SHA256 | | MD5 checksum |

<br>

**Example** : Install a `hello-world` software module.

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
> | path | `/features/SoftwareUpdatable/outbox/messages/install` | A path to the `SoftwareUpdatable` Feature, it's message channel, and `install` command |
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
> | path | `/features/SoftwareUpdatable/inbox/messages/download` | A path to the `SoftwareUpdatable` Feature, it's message channel, and `download` command |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response is required |
> | content-type | `application/json` | The content type |
> | correlation-id | UUID | Used for correlating protocol messages, the same correlation-id as the sent back response message |
> | **Value** | | JSON representation of the software modules that will be downloaded |
> | correlationId | | Unique identifier that is used to associate and track the series of messages |
> | weight | | The weight is the priority in case of multiple, parallel instructions |
> | metadata | | The metadata is any other information which should be passed to the device |
> | forced | true/false | Remove the software modules forcefully |
> | ***softwareModules*** | | An array of software modules that will be downloaded|
> | metadata | | The metadata is any other information which should be passed to the device |
> | **softwareModule** | | A unique identifier for the software module |
> | name | | The name of the software module |
> | version | | The version of the software module |
> | **artifacts** | | An array of artifacts contained in the software module |
> | filename | | The file name of the artifact |
> | size | | Artifact file size in bytes |
> | **download** | | A map with protocols and links for downloading the artifacts |
> | key | HTTP/HTTPS/FTP/SFTP | Available transport protocols |
> | url | | URL to download the artifact |
> | md5url | | MD5URL to download the MD5SUM file |
> | **checksums** | | A map with checksums to verify the proper download |
> | MD5 | | MD5 checksum |
> | SHA1 | | MD5 checksum |
> | SHA256 | | MD5 checksum |

<br>

**Example** : Download a hello-world software module.

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
> | path | `/features/SoftwareUpdatable/outbox/messages/download` | A path to the `SoftwareUpdatable` Feature, it's message channel, and `download` command |
> | **Headers** | | Additional headers |
> | correlation-id | container UUID | The same correlation id as the request message |
> | **Status** | | Status of the  `download` operation |

<br>

**Example** : Successful download response.

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