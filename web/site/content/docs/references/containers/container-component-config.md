---
title: "Container configuration as Desired State component"
type: docs
description: >
  Customize the deployment of a container instance as a Desired State component.
weight: 3
---

### Domain Identifier
The deafult domain identifier for the Containers Update Agent is `containers`.
This can be modified within the update agent section in the [container management JSON config file](./container-manager-config).

### Containers Update Agent Properties 
To control the container update agent behavior through desired state specification. As defined in the Desired State Specification, all properties are of type string.

| Key | Required | Default | Description |
| - | - | - | - |
| systemContainers | No | | Comma-separated list of container names that shall not be processed by the update agent during the application of the given desired state. This configuration option can be used to temporarily override the general `systemContainers` setting from the update agent section in the container management JSON config file. The setting is valid only for the given desired state where it is present. |

### Container Properties
To control all aspects of the container instance behavior. As defined in the Desired State Specification, all properties are of type string.

| Key | Required | Default | Description |
| - | - | - | - |
| **General config** | | | |
| image | Yes | | Fully qualified image reference, that follows the {{% refn "https://github.com/opencontainers/image-spec" %}}OCI Image Specification{{% /refn %}}, the format is: `host[:port]/[namespace/]name:tag`. |
| env | No | | Sets the provided environment variable in the root container's process environment.Example: `VAR1=2`. If `VAR1`= is used, the environment variable would be set to empty. If `VAR1` is used, the environment variable would be removed from the container environment inherited from the image. The property can be included multiple times, each one specifying another environment variable. |
| cmd | No | | Command with arguments that is executed upon the container’s start. The property can be included multiple times (order is important), each one specifying another command argument. |
| **Host config** | | | |
| device | No | | Device to be made available in the container and optional cgroups permissions configuration. Both path on host and in container must be set. Possible cgroup permissions options are "r" (read), "w" (write), "m" (mknod) and all combinations of the three are possible. If not set, "rwm" is default device configuration. Example: `/dev/ttyACM0:/dev/ttyUSB0[:rwm]`. The property can be included multiple times, each one specifying a separate device. |
| port | No | | Port to be mapped from the host to the container instance. Format: `[<host-ip>:]<host-port>:<container-port>[-<range>][/<proto>]`. Most common use-case: `80:80`. Mapping the container's 80 port to a host port in the 5000-6000 range: `5000-6000:80/udp`. Specifying port protocol (default is tcp): `80:80/udp`. By default the port mapping will set on all network interfaces, but this is also manageable: `0.0.0.0:80-100:80/udp`. The property can be included multiple times, each one specifying another port mapping. |
| network | No | bridge | Sets the networking mode for the container. Possible options are: `bridge` - the container is connected to the default bridge network interface of the engine and is assigned an IP. `host` - the container shares the network stack of the host (use with caution as this breaks the network's isolation!) |
| host | No | | Extra host to be added in the current container's `/etc/hosts` file. Example: hostname1:<IP1>. If the IP of the host machine is to be added to the container's hosts file the reserved `host_ip[_<network-interface>]` must be provided. Example: `local.host.machine.ip.custom.if:host_ip_myNetIf0` - this will automatically resolve the host's IP on the `myNetIf0` network interface and add it to the container's hosts file. `local.host.machine.ip.default.bridge:host_ip` - this will automatically resolve the host's IP on the default bridge network interface for container management and add it to the container's hosts file if the container is configured to use it. The property can be included multiple times, each one specifying another extra host |
| mount | No | | Sets mount points so a source directory on the host can be accessed via a destination directory in the container. Format: `source:destination[:propagation_mode]`. If the propagation mode parameter is omitted, `rprivate` will be set by default. Available propagation modes are: `rprivate`, `private`, `rshared`, `shared`, `rslave`, `slave`. The property can be included multiple times, each one specifying another mount point. |
| **IO config** | | | |
| terminal | No | false | Boolean flag. Enables terminal for the current container, e.g. attach standard streams to a TTY. |
| interactive | No | false | Boolean flag. Enables interaction with the container, e.g. open the terminal's standard input for an interaction with the container. |
| priviledged | No | false | Boolean flag. Creates the container as privileged, grants root capabilities to all devices on the host system |
| **Restart policy config** | | | |
| restartPolicy | string | unless-stopped | The container's restart policy, the supported values are: `always` - an attempt to restart the container will be me made each time the container exits regardless of the exit code, `no` - no attempts to restart the container for any reason will be made, `on-failure` - restart attempts will be made if the container exits with an exit code != 0, `unless-stopped` - restart attempts will be made only if the container has not been stopped by the user. |
| restartMaxRetries | No | 1 | Integer value. Maximum number of retries that are made to restart the container on exit with fail, valid only if the `restartPolicy` is `on-failure`. |
| restartTimeout | No | 30 | Integer value. Timeout period in seconds for each retry that is made to restart the container on exit with fail, valid only if the `restartPolicy` is `on-failure`. |
| **Logging config** | | | |
| logDriver | No | json-file | Sets the type of the log driver to be used for the container - `json-file`, `none`. |
| logMaxFiles | No | 2 | Integer value. Sets the max number of log files to be rotated - applicable for `json-file` log driver only. |
| logMaxSize | No | 100M | Sets the max size of the logs files for rotation in the form of 1, 1.2m, 1g, etc. - applicable for `json-file` log driver only. |
| logPath | No | | Sets the path to the directory where the log files will be stored - applicable for `json-file` log driver only. |
| logMode | No | blocking | Sets the mode of the logger - `blocking`, `non-blocking`. |
| logMaxBufferSize | No | 1M | Sets the max size of the logger buffer in the form of 1, 1.2m - applicable for `non-blocking` mode only. |
| **Resources config** | | | |
| memory | No | | Sets the max amount of memory the container can use in the form of 200m, 1.2g. The minimum allowed value is 3m. By default, a container has no memory constraints. |
| memorySwap | No | | Sets the total amount of memory + swap that the container can use in the form of 200m, 1.2g. If set must not be smaller than `memory`. If equal to `memory`, than the container will not have access to swap. If not set and `memory` is set, than the container can use as much swap as the `memory` setting. If set to -1, the container can use unlimited swap, up to the amount available on the host. |
| memoryReservation | No | | Sets a soft memory limitation in the form of 200m, 1.2g. Must be smaller than `memory`. When the system detects memory contention or low memory, control groups are pushed back to their soft limits. There is no guarantee that the container memory usage will not exceed the soft limit. |

#### Desired State Containers Domain Example

```json
{
	"domains": [
		{
			"id": "containers",
			"config": [
				{
					"key": "systemContainers",
					"value": "self-update-agent"
				}
			],
			"components": [
				{
					"id": "hello-world",
					"version": "latest",
					"config": [
						{
							"key": "image",
							"value": "docker.io/library/hello-world:latest"
						},
						{
							"key": "env",
							"value": "x=y"
						},
						{
							"key": "env",
							"value": "a=b"
						},
						{
							"key": "cmd",
							"value": "arg1"
						},
						{
							"key": "cmd",
							"value": "arg2"
						},
						{
							"key": "device",
							"value": "/dev/tty:/dev/tty:rw"
						},
						{
							"key": "port",
							"value": "80:80/tcp"
						},
						{
							"key": "network",
							"value": "host"
						},
						{
							"key": "host",
							"value": "host_name"
						},
						{
							"key": "mount",
							"value": "/data:/data:private"
						},
						{
							"key": "terminal",
							"value": "true"
						},
						{
							"key": "interactive",
							"value": "true"
						},
						{
							"key": "priviledged",
							"value": "true"
						},
						{
							"key": "restartPolicy",
							"value": "always"
						},
						{
							"key": "restartMaxRetries",
							"value": "3"
						},
						{
							"key": "restartTimeout",
							"value": "1000"
						},
						{
							"key": "logDriver",
							"value": "json-file"
						},
						{
							"key": "logMaxFiles",
							"value": "3"
						},
						{
							"key": "logMaxSize",
							"value": "5M"
						},
						{
							"key": "logPath",
							"value": "/var/log"
						},
						{
							"key": "logMode",
							"value": "blocking"
						},
						{
							"key": "logMaxBufferSize",
							"value": "1M"
						},
						{
							"key": "memory",
							"value": "200M"
						},
						{
							"key": "memorySwap",
							"value": "300M"
						},
						{
							"key": "memoryReservation",
							"value": "100M"
						}
					]
				}
			]
		}
	]
}
```