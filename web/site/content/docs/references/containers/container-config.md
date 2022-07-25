---
title: "Container configuration"
type: docs
description: >
  Customize the deployment of container instance.
weight: 2
---

### Properties

To control all aspects of the container instance behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| container_name | string | \<container ID\> | User-defined name of the container |
| domain_name | string | \<container name\>-domain | Domain name inside the container |
| host_name | string | \<container name\>-host | Host name for the container  |
| **Image** | | | |
| name | string | | The full name of the container image |
| **Encrypted images** | | | |
| keys | string[] | | Filepath to secret keys (GPG private key ring, JWE or PKCS7) and an optional password separated by colon |
| recipients | string[] | | Protocol and certificate filepath separated by colon. Recipient certificate of the image, used only for PKCS7 and must be an x509 |
| **Mount points** | | | |
| destination | string | | Path where the file or directory is mounted in the container |
| source | string | | Path to the file or directory on the host |
| propagation_mode | string | | Allows mounting volumes in propagation mode on a container, the supported modes are: `rprivate`, `private`, `rshared`, `shared`, `rslave`, `slave` |
| **Container config** | | | |
| env | string[] | | Sets environment variables in the root container's process environment |
| **Host config** | | | |
| privileged | bool | false | Sets a container as privileged if the value is true |
| network_mode | string | bridge | Configure the networking mode for the container, the possible options are `bridge` and `host` |
| **Devices** | | | |
| path_on_host | string | | Path on the host for mapped device |
| path_in_container | string | | Path in the container for mapped device. |
| cgroup_permissions | string | rwm | Sets cgroup permissions for the device access, possible options are: r (read), w (write), m (mknod) and all combinations of the three are possible |
| **Restart policy** | | | |
| type | string | unless-stopped | Strategy how to restart a container automatically, the supported types are: `always`, `no`, `on-failure` and `unless-stopped` |
| maximum_retry_count | int | 1 | The number of retries that are made to restart the container on exit only if the policy is set to `on-failure` |
| retry_timeout | int | 100 | The time out period in seconds for each retry that is made to restart the container on exit only if the policy is set to `on-failure` |
| **Extra hosts** | | | |
| extra_hosts | string[] | | Specify which hosts to be added in the current container's /etc/hosts file |
| **Port mapping** | | | |
| proto | string | tcp | The protocol for port mapping from the host to a container |
| container_port | int | | Port of the mapped container |
| host_ip | string | | IP of the mapped host |
| host_port | int | | Beginning of the range of host port |
| host_port_end | int | | End of the range of host port |
| **Driver config** | | | |
| type | string | json-file | The type of the full log driver configuration, the possible options are: `json-file` or `none` |
| max_files | int | 2 | The maximum number of files the logs can be stored in. After the max number is reached, the files are rotated - the oldest is removed and the index is updated. This is applicable for json-file log driver only. |
| max_size | string | 100M | The max size of the logs. This is applicable for json-file log driver only. |
| root_dir | string | *Default container-management's meta path* | The target directory to store the logs per container. A sub-directory is created for each container’s ID. This is applicable for json-file log driver only. |
| **Mode config** | | | |
| mode | string | blocking | The supported logging modes are:<ul><li>blocking - instructs the logger not to use buffering</li><li>non-blocking - instructs the logger to use buffering</li></ul> |
| max_buffer_size | string | 1M | Sets the max size of the logger buffer in the from 1, 1.2m. This option is applicable for non-blocking mode only. |
| **Resources** | | | |
| memory | string | | Hard memory limitation of a container, the minimum allowed value is 3M |
| memory_reservation | string | | Soft memory limitation of a container. If `memory` is specified, the `memory_reservation` must be smaller than it |
| memory_swap | string | | The total amount of memory and swap that the container can use |

#### Example
The minimal required information for the container instance is the receptive container image.

```json
{
  "image": {
    "name": "docker.io/nginxdemos/hello:plain-text"
  }
}
```

### Template

The configuration can be further adjusted according to the use case. The following template illustrates all possible properties with their default values.

```json
{
   "container_name":"test-name",
   "image":{
      "name":"docker.io/nginxdemos/hello:plain-text",
      "decrypt_config":{
         "keys":[
            "/home/user/pkcs7/key.pem"
         ],
         "recipients":[
            "pkcs7:/home/user/pkcs7/cert.pem"
         ]
      }
   },
   "domain_name":"test-name-domain",
   "host_name":"test-name-host",
   "mount_points":[
      {
         "destination":"/etc/test-config",
         "source":"/etc/test-source-config",
         "propagation_mode":"rprivate"
      }
   ],
   "config":{
      "env":[
         "VAR_1=1",
         "VAR_2=my var"
      ]
   },
   "host_config":{
      "devices":[
         {
            "path_on_host":"/dev/ttyACM0",
            "path_in_container":"/dev/ttyACM1",
            "cgroup_permissions":"rwm"
         }
      ],
      "privileged":false,
      "network_mode":"bridge",
      "restart_policy":{
         "maximum_retry_count":0,
         "retry_timeout":0,
         "type":"unless-stopped"
      },
      "extra_hosts":[
         "ctrhost:host_ip"
      ],
      "port_mappings":[
         {
            "proto":"tcp",
            "container_port":80,
            "host_ip":"192.168.1.101",
            "host_port":81,
            "host_port_end":82
         }
      ],
      "log_config":{
         "driver_config":{
            "type":"json-file",
            "max_files":2,
            "max_size":"100M",
            "root_dir":"/my/secured/dir"
         },
         "mode_config":{
            "mode":"non-blocking",
            "max_buffer_size":"5M"
         }
      },
      "resources":{
         "memory":"200M",
         "memory_reservation":"100M",
         "memory_swap":"500M"
      }
   }
}
```