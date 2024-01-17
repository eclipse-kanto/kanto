---
title: "Manager configuration"
type: docs
description: >
  Customize the container manager components.
weight: 2
---

### Properties

To control all aspects of the container manager behavior.

| Property | Type | Default | Description |
| - | - | - | - |
| home_dir | string | /var/lib/container-management | Home directory for the container manager data |
| exec_root_dir | string | /var/run/container-management | Root directory for the container manager's executable artifacts |
| container_client_sid | string | container-management.service.local.v1.service-containerd-client | Unique identifier that is used for an interaction with the runtime |
| network_manager_sid | string | container-management.service.local.v1.service-libnetwork-manager | Unique identifier that is used for networking |
| default_ctrs_stop_timeout | int | 30 | (Deprecated since v0.1.0-M4, type changed to string) Timeout in seconds for a container to stop gracefully, otherwise its root process will be forcefully stopped |
| default_ctrs_stop_timeout | string | 30s | Timeout for a container to stop gracefully in duration string format (e.g. 1h2m3s5ms), otherwise its root process will be forcefully stopped |
| **Runtime** | | | |
| default_ns | string | kanto-cm | Namespace that is used by the runtime for isolation |
| address_path | string | /run/containerd/containerd.sock | Path to the runtime's communication endpoint |
| home_dir | string | /var/lib/container-management | Home directory for the runtime data |
| exec_root_dir | string | /var/run/container-management | Root directory for the runtime's executable artifacts |
| image_dec_keys | string[] | | Private keys (GPG private key ring, JWE or PKCS7) used for decrypting container images, the format is: `filepath_private_key[:password]` |
| image_dec_recipients | string[] | | Recipients (only for PKCS7 and must be an x509) used for decrypting container images, the format is: `pkcs7:filepath_x509_certificate` |
| runc_runtime | string | io.containerd.runc.v2 | Runc communication mode, the possible values are: io.containerd.runtime.v1.linux, io.containerd.runc.v1 and io.containerd.runc.v2 |
| image_expiry | string | 744h | Time period for the cached images and content to be kept in the form of e.g. 72h3m0.5s |
| image_expiry_disable | bool | false | Disable expiry management of cached images and content, must be used with caution as it may lead to large memory volumes being persistently allocated |
| lease_id | string | kanto-cm.lease | Lease identifier to be used for container resources persistence |
| image_verifier_type | string | none | The image verifier type - possible values are none and notation, when set to none image signatures wil not be verified |
| image_verifier_config | string |  | The configuration of the image verifier, as comma separated {key}={value} pairs - possible keys for notation verifier are `configDir` and `libexecDir`, for more info check [notation documentation](https://notaryproject.dev/docs/user-guides/how-to/directory-structure/#user-level) |
| **Registry access - secure** | | | |
| user_id | string | | User unique identifier to authenticate to the image registry |
| password | string | | Password to authenticate to the image registry |
| root_ca | string | | PEM encoded CA certificates file |
| client_cert | string | | PEM encoded certificate file to authenticate to the image registry |
| client_key | string | | PEM encoded unencrypted private key file to authenticate to the image registry |
| **Registry access - insecure** | | | |
| insecure_registries | string[] | localhost | Image registries that do not use valid certificates or do not require a HTTPS connection, the format is: `host[:port]` |
| **Networking** | | | |
| home_dir | string | /var/lib/container-management | Home directory for the network manager data |
| exec_root_dir | string | /var/run/container-management | Root directory for the network manager's executable artifacts |
| **Networking - bridge** | | | |
| name | string | kanto-cm0 | Bridge name |
| ip4 | string | | Bridge IPv4 address |
| fcidr4 | string | | IPv4 address range for the bridge, using the standard CIDR notation |
| gwip4 | string | | Bridge gateway IPv4 address |
| enable_ip6 | bool | false | Permit the bridge IPv6 support |
| mtu | int | 1500 | Bridge maximum transmission unit in bytes |
| icc | bool | true | Permit the inter-container communication |
| ip_tables | bool | true | Permit the IP tables rules |
| ip_forward | bool | true | Permit the IP forwarding |
| ip_masq | bool | true | Permit the IP masquerading |
| userland_proxy | bool | false | Forbid the userland proxy for the loopback traffic |
| **Local communication** | | | |
| protocol | string | unix | Communication protocol used for accessing the gRPC server, the possible values are: tcp, tcp4, tcp6, unix or unixpacket |
| address_path | string | /run/container-management/container-management.sock | Path to the gRPC server's communication endpoint |
| **Digital twin** | | | |
| enable | bool | true | Permit the container manager digital twin representation |
| home_dir | string | /var/lib/container-management | Home directory for the digital twin data |
| features | string[] | ContainerFactory, SoftwareUpdatable, Metrics | Features that will be registered for the container manager digital twin, the possible values are: ContainerFactory, SoftwareUpdatable and Metrics |
| **Digital twin - connectivity** | | | |
| broker_url | string | tcp://localhost:1883 | Address of the MQTT server/broker that the container manager will connect for the local communication, the format is: `scheme://host:port` |
| keep_alive | int | 20000 | Keep alive duration in milliseconds for the MQTT requests |
| disconnect_timeout | int | 250 | Disconnect timeout in milliseconds for the MQTT server/broker |
| client_username | string | | Username that is a part of the credentials |
| client_password | string | | Password that is a part of the credentials |
| connect_timeout | int | 30000 | Connect timeout in milliseconds for the MQTT server/broker |
| acknowledge_timeout | int | 15000 | Acknowledge timeout in milliseconds for the MQTT requests |
| subscribe_timeout | int | 15000 | Subscribe timeout in milliseconds for the MQTT requests |
| unsubscribe_timeout | int | 5000 | Unsubscribe timeout in milliseconds for the MQTT requests |
| **Digital twin - connectivity - TLS** | | | |
| root_ca | string | | PEM encoded CA certificates file |
| client_cert | string | | PEM encoded certificate file to authenticate to the MQTT server/broker |
| client_key | string | | PEM encoded unencrypted private key file to authenticate to the MQTT server/broker |
| **Logging** | | | |
| log_file | string | log/container-management.log | Path to the file where the container manager's log messages are written |
| log_level | string | INFO | All log messages at this or a higher level will be logged, the log levels in descending order are: ERROR, WARN, INFO, DEBUG and TRACE |
| log_file_count | int | 5 | Log file maximum rotations count |
| log_file_max_age | int | 28 | Log file rotations maximum age in days, use 0 to not remove old log files |
| log_file_size | int | 2 | Log file size in MB before it gets rotated |
| syslog | bool | false | Route logs to the local syslog |
| **Deployment** | | | |
| enable | bool | true | Permit the deployment manager service providing installation/update of containers via the container descriptor files |
| mode | string | update | Deployment manager mode, the possible values are: init (container descriptors are processed only on first start, new containers are deployed and started), update (container descriptors are processed on each restart, new containers can be deployed and started, existing containers may be updated, no container removals) |
| home_dir | string | /var/lib/container-management | Home directory for the deployment manager data |
| ctr_dir | string | /etc/container-management/containers | Directory containing descriptors of containers that will be automatically deployed on first start or updated on restart |

### Example

The minimal required configuration that sets a timeout period of 5 seconds for the managed containers to stop gracefully.

```json
{
    "manager": {
        "default_ctrs_stop_timeout": 5
    },
    "log": {
        "log_file": "/var/log/container-management/container-management.log"
    }
}
```

### Template

The configuration can be further adjusted according to the use case.
The following template illustrates all possible properties with their default values.

{{% warn %}}
Be aware that in the registry configuration the host (used as a key) has to be set instead of the default empty string, the format is: host[:port]
{{% /warn %}}

```json
{
    "manager": {
        "home_dir": "/var/lib/container-management",
        "exec_root_dir": "/var/run/container-management",
        "container_client_sid": "container-management.service.local.v1.service-containerd-client",
        "network_manager_sid": "container-management.service.local.v1.service-libnetwork-manager",
        "default_ctrs_stop_timeout": 30
    },
    "containers": {
        "default_ns": "kanto-cm",
        "address_path": "/run/containerd/containerd.sock",
        "exec_root_dir": "/var/run/container-management",
        "home_dir": "/var/lib/container-management",
        "image_dec_keys": [],
        "image_dec_recipients": [],
        "runc_runtime": "io.containerd.runc.v2",
        "image_expiry": "744h",
        "image_expiry_disable": false,
        "image_verifier_type": "none",
        "image_verifier_config": "configDir=/home/user/.config/notation",
        "lease_id": "kanto-cm.lease",
        "registry_configurations": {
            "": {
                "credentials": {
                    "user_id": "",
                    "password": ""
                },
                "transport": {
                    "root_ca": "",
                    "client_cert": "",
                    "client_key": ""
                }
            }
        },
        "insecure_registries": [
            "localhost"
        ]
    },
    "network": {
        "home_dir": "/var/lib/container-management",
        "exec_root_dir": "/var/run/container-management",
        "default_bridge": {
            "name": "kanto-cm0",
            "ip4": "",
            "fcidr4": "",
            "enable_ip6": false,
            "mtu": 1500,
            "icc": true,
            "ip_tables": true,
            "ip_forward": true,
            "ip_masq": true,
            "userland_proxy": false
        }
    },
    "grpc_server": {
        "protocol": "unix",
        "address_path": "/run/container-management/container-management.sock"
    },
    "things": {
        "enable": true,
        "home_dir": "/var/lib/container-management",
        "features": [
            "ContainerFactory",
            "SoftwareUpdatable",
            "Metrics"
        ],
        "connection": {
            "broker_url": "tcp://localhost:1883",
            "keep_alive": 20000,
            "disconnect_timeout": 250,
            "client_username": "",
            "client_password": "",
            "connect_timeout": 30000,
            "acknowledge_timeout": 15000,
            "subscribe_timeout": 15000,
            "unsubscribe_timeout": 5000,
            "transport": {
                "root_ca": "",
                "client_cert": "",
                "client_key": ""
            }
        }
    },
    "log": {
        "log_file": "log/container-management.log",
        "log_level": "INFO",
        "log_file_count": 5,
        "log_file_size": 2,
        "log_file_max_age": 28,
        "syslog": false
    },
    "deployment": {
        "enable": true,
        "mode": "update",
        "home_dir": "/var/lib/container-management",
        "ctr_dir": "/etc/container-management/containers"
    }
}
```
