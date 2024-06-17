---
title: "Container Factory API"
type: docs
description: >
  The container factory service provides the ability to create new containers form a container image, or from a [container configuration](../container-config.md).
weight: 2
---

## **Create**
Create a container from a single container image reference with an option to start it.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//create`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/create` | Information about the affected Thing and the type of operation |
> | path | `/features/ContainerFactory/inbox/messages/create` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response required |
> | content-type | `application/json` | The content type |
> | correlation-id | container UUID | The container UUID |
> | **Value** | | |
> | imageRef | URL | Container image URL |
> | start | true/false | If the created container will be started ot only created |

<br>

**Example** : In this example, you can create and automatically start a new `Hello World` container.

**Topic:** `command//edge:device/req//create`
```json
{
	"topic":"edge/device/things/live/messages/create",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/ContainerFactory/inbox/messages/create",
	"value":{
		"imageRef":"docker.io/library/hello-world:latest",
		"start":true
	}
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//create`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/create` | Information about the affected Thing and the type of operation |
> | path | `/features/ContainerFactory/outbox/messages/create` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | content-type | application/json | The content type |
> | correlation-id | \<UUID\> | The same correlation id as the sent request message |
> | **Value** | | UUID of the created container |

<br>

**Example** : The response of the create operation.

**Topic:** `command//edge:device/res//create``
```json
{
	"topic":"edge/device/things/live/messages/create",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/ContainerFactory/outbox/messages/create",
	"value":"<Container UUID>"
}
```
</details>

## **Create  with config**
Create a container from a given container configuration.

<details>
  <summary>Request</summary>

**Hono Command:** `command//<name>:<namespace>:edge:containers/req//createWithConfig`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/createWithConfig` | Information about the affected Thing and the type of operation |
> | path | `/features/ContainerFactory/inbox/messages/createWithConfig` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | response-required | true/false | If response required |
> | content-type | `application/json` | The content type |
> | correlation-id  | container UUID | The container UUID |
> | **Value** | | |
> | imageRef | URL | Fully qualified image reference, that follows the {{% refn "https://github.com/opencontainers/image-spec" %}}OCI Image Specification{{% /refn %}}, the format is: `host[:port]/[namespace/]name:tag` |
> | start | true/false | Force to start created container |
> | ***config*** | | json presentation of the configuration |
> | domainName | | Domain name inside the container, if omitted the container's domain name will be set to a system-defined value |
> | hostName | | Host name for the container, if omitted the container's hostname will be set to a system-defined value |
> | env | | An array of environment variables that are set into the container |
> | cmd | | An array of command with arguments that is executed upon the container's start |
> | privileged | false | Grant root capabilities to all devices on the host system|
> | extraHosts | | An array of additional extra host names to IP address mappings added to the container network configuration, the format is: hostname:ip. If the IP of the host machine is to be added to the container’s hosts file the reserved host_ip[_<network-interface>] must be provided. If only host_ip (the network-interface part is skipped) is used, by default it will be resolved to the host’s IP on the default bridge network interface for containerm (the default configuration is kanto-cm0) and add it to the container’s hosts file. If the IP of a container in the same bridge network is to be added to the hosts file the reserved container_<container-host_name> must be provided. |
> | extraCapabilities | | An array of additional capabilities for a container |
> | networkMode | | The container’s networking capabilities type based on the desired communication mode, the possible options are: bridge or host |
> | openStdin | true/false | Open the terminal's standard input for an interaction with the current container |
> | tty | true/false | Attach standard streams to a TTY|
> | **mountPoints** | | An array of the mount points |
> | source | | Path to the file or directory on the host that is referred from within the container |
> | destination | | Path to the file or directory that is mounted inside the container |
> | propagationMode | | Bind propagation for the mount, supported are: rprivate, private, rshared, shared, rslave or slave |
> | **decryption** | | |
> | keys | | A string array of private keys (GPG private key ring, JWE or PKCS7) used for decrypting the container's image, the format is: `filepath_private_key[:password]` |
> | recipients | | A string array of recipients (only for PKCS7 and must be an x509) used for decrypting the container's image, the format is: `pkcs7:filepath_x509_certificate` |
> | **devices** | | An array of accessible devices from the host |
> | pathOnHost | | Path to the device on the host |
> | pathInContainer | | Path to the device in the container |
> | cgroupPermissions | rwm | Cgroup permissions for the device access, possible options are: r(read), w(write), m(mknod) and all combinations are possible |
> | **restartPolicy** | | The container restart policy |
> | type | unless-stopped | The container's restart policy, the supported types are: always, no, on-failure and unless-stopped |
> | maxRetryCount | | Maximum number of retries that are made to restart the container on exit with fail, if the `type` is on-failure |
> | retryTimeout | | Timeout period in seconds for each retry that is made to restart the container on exit with fail, if the `type` is on-failure |
> | **portMappings** | | An array of port mappings from the host to a container |
> | proto | tcp | Protocol used for the port mapping from the container to the host, the possible options are: tcp and udp |
> | containerPort | | Port number on the container that is mapped to the host port |
> | hostIP | 0.0.0.0 | Host IP address |
> | hostPort | | Beginning of the host ports range |
> | hostPortEnd | <host_port> | Ending of the host ports range |
> | **log** | | |
> | type | json-file | Type in which the logs are produced, the possible options are: json-file or none |
> | maxFiles | 2 | Maximum log files before getting rotated |
> | maxSize | 100M | Maximum log file size before getting rotated as a number with a unit suffix of B, K, M and G |
> | rootDir | <meta_path>/containers/<container_id> | Root directory where the container's log messages are stored |
> | mode | blocking | Messaging delivery mode from the container to the log driver, the supported modes are: blocking and non-blocking |
> | maxBufferSize | 1M | Maximum size of the buffered container's log messages in a non-blocking mode as a number with a unit suffix of B, K, M and G |
> | **resources** | | |
> | memory | | Hard memory limitation of the container as a number with a unit suffix of B, K, M and G, the minimum allowed value is 3M |
> | memoryReservation | | Soft memory limitation of the container as a number with a unit suffix of B, K, M and G, if `memory` is specified, the `memoryReservation` must be smaller than it |
> | memorySwap | | Total amount of memory and swap that the container can use as a number with a unit suffix of B, K, M and G, use -1 to allow the container to use unlimited swap |

<br>

**Example** : In this example, you can create and automatically start a new `Hello World` container.

**Topic:** `command//edge:device/req//createWithConfig`
```json
{
	"topic":"edge/device/things/live/messages/createWithConfig",
	"headers":{
		"response-required":true,
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/ContainerFactory/inbox/messages/createWithConfig",
	"value":{
		"imageRef":"docker.io/library/influxdb:1.8.4",
		"start":true,
		"config":{
			"domainName": "",
			"hostName": "",
			"env": [],
			"cmd": [],
			"privileged": false,
			"extraHosts": ["ctrhost:host_ip"],
			"extraCapabilities": [],
			"networkMode": "bridge",
			"openStdin": false,
			"tty": false,
			"mountPoints": [
				{
					"source": "",
					"destination": "",
					"propagationMode": "rprivate"
				}
			],
			"decryption": {
				"keys": [],
				"recipients": []
			},
			"devices": [
				{
					"pathOnHost": "",
					"pathInContainer": "",
					"cgroupPermissions": "rwm"
				}
			],
			"restartPolicy": {
				"type": "unless-stopped",
				"maxRetryCount": 0,
				"retryTimeout": 0
			},
			"portMappings":[
				{
					"proto": "tcp",
					"containerPort": 80,
					"hostIP": "0.0.0.0",
					"hostPort": 5000,
					"hostPortEnd": 5005,
				}
			],
			"log": {
				"type": "json-file",
				"maxFiles": 2,
				"maxSize": "100M",
				"rootDir": "",
				"mode": "blocking",
				"maxBufferSize": "1M"
			},
			"resources": {
				"memory": "",
				"memoryReservation": "",
				"memorySwap": ""
			},
		}
    }
}
```
</details>

<details>
  <summary>Response</summary>

**Hono Command** : `command//<name>:<namespace>:edge:containers/res//createWithConfig`

**Ditto Message:**

> | Name | Value | Description |
> | - | - | - |
> | topic | `<name>/<namespace>/things/live/messages/createWithConfig` | Information about the affected Thing and the type of operation |
> | path | `/features/ContainerFactory/outbox/messages/createWithConfig` | A path that references a part of a Thing which is affected by this message |
> | **Headers** | | Additional headers |
> | content-type | `application/json` | The content type |
> | correlation-id | \<UUID\> | |
> | **Value** | | UUID of the created container |

<br>

**Example** : The response of the create with config operation.

**Topic:** `command//edge:device/res//createWithConfig``
```json
{
	"topic":"edge/device/things/live/messages/createWithConfig",
	"headers":{
		"content-type":"application/json",
		"correlation-id":"<UUID>"
	},
	"path":"/features/ContainerFactory/outbox/messages/createWithConfig",
	"value":"<Container UUID>"
}
```
</details>