---
title: "Perform OTA update"
type: docs
description: >
    Perform an OTA update on your edge device.
weight: 3
---

By following the steps below you will publish a simple desired state specification to a MQTT topic
and then the specification will be forwarded to the Update Manager, which will trigger an OTA update on
the edge device.

A simple listener for MQTT messages will track the progress and the status of the update process.
### Before you begin

To ensure that all steps in this guide can be executed, you need:

* Debian-based linux distribution and the `apt` command line tool

* If you don't have an installed and running Eclipse Kanto, follow {{% relrefn "install" %}} Install Eclipse Kanto {{% /relrefn %}}

* Run Mosquitto listener to track the desired state feedback messages by executing:
  ```shell
  mosquitto_sub -t "deviceupdate/desiredstatefeedback" -v
  ```

### Publish the Desired State specification

The starting point of the OTA update process is to publish the example Desired State specification below on the MQTT topic `deviceupdate\desiredstate`:
```
{
	"activityId": "testActivityId",
	"payload": {
		"domains": [
			{
				"id": "containers",
				"config": [],
				"components": [
					{
						"id": "influxdb",
						"version": "2.7.1",
						"config": [
							{
								"key": "image",
								"value": "docker.io/library/influxdb:2.7.1"
							}
						]
					},
					{
						"id": "alpine",
						"version": "latest",
						"config": [
							{
								"key": "image",
								"value": "docker.io/library/alpine:latest"
							}
						]
					},
					{
						"id": "hello-world",
						"version": "latest",
						"config": [
							{
								"key": "image",
								"value": "docker.io/library/hello-world:latest"
							}
						]
					}
				]
			}
		]
	}
}
```
The desired state specification in this case consists of single domain section definition for the containers domain and a three container components - `influxdb`, `hello-world` and `alpine` image.

### Apply Desired State specification

The Update Manager receives the published desired state to the local Mosquitto broker, splits the specification (in this case into single domain) and then
distributes the processed specification to the domain agents which initiates the actual update process logic on the domain agents side.

The update process is organized into multiple phases, which are triggered by sending specific desired state commands (`DOWNLOAD/UPDATE/ACTIVATE`).

In the example scenario, the three images for the three container components will be pulled (if not available in the cache locally), created as containers during the `UPDATING` phase and
started in the `ACTIVATING` phase.

### Monitor OTA update progress

During the OTA update, the progress can be tracked and monitored by the Mosquitto listener for the desired state feedback messages, started in the prerequisite section above.

The Update Manager reports to the local Mosquitto broker at a time interval of a second the status of the active update process. For example:
```
{
	"activityId": "testActivityId",
	"payload": {
		"status": "RUNNING",
		"actions": [
			{
				"component": {
					"id": "containers:hello-world",
					"version": "latest"
				},
				"status": "UPDATE_SUCCESS",
				"message": "New container instance is started."
			},
			{
				"component": {
					"id": "containers:influxdb",
					"version": "2.7.1"
				},
				"status": "UPDATE_SUCCESS",
				"message": "New container instance is started."
			},
			{
				"component": {
					"id": "containers:alpine",
					"version": "latest"
				},
				"status": "UPDATING",
				"message": "New container created."
			}
		]
	}
}
```

### List containers

After the update process is completed, list the installed containers by executing the command `kanto-cm list` to the verify the desired state is applied correctly.

The output of the command should display the info about the three containers, described in the desired state specification. The `influxdb` is expected to be in `RUNNING` state and
the other containers in status `EXITED`. For example :
```
ID                                    |Name         |Image                               |Status   |Finished At |Exit Code
|-------------------------------------|-------------|------------------------------------|----------------------|---------
7fe6b689-eb76-476d-a730-c2f422d6e8ea  |influxdb     |docker.io/library/influxdb:2.7.1    |Running  |            |0
c36523d7-8d17-4255-ae0c-37f11003f658  |hello-world  |docker.io/library/hello-world:latest|Exited   |            |0
9b99978b-2593-4736-bb52-7a07be4a7ed1  |alpine       |docker.io/library/alpine:latest     |Exited   |            |0
```

### Update Desired State specification

As described in the sections above, apply the Desired State specification below to update the existing desired state. The update changes affect two containers - `alpine` and `influxdb`. Being not present in the updated Desired State specification, the `alpine` container will be removed from the system. The `influxdb` will be updated to the latest version. The last container - `hello-world` is not affected and any events will be not reported from the container update agent for this particular container.

```
{
	"activityId": "testActivityId",
	"payload": {
		"domains": [
			{
				"id": "containers",
				"config": [],
				"components": [
					{
						"id": "influxdb",
						"version": "latest",
						"config": [
							{
								"key": "image",
								"value": "docker.io/library/influxdb:latest"
							}
						]
					},
					{
						"id": "hello-world",
						"version": "latest",
						"config": [
							{
								"key": "image",
								"value": "docker.io/library/hello-world:latest"
							}
						]
					}
				]
			}
		]
	}
}
```

### List updated containers

After the update process of the existing desired state is completed, list again the available containers to the verify the desired state is updated correctly.

The output of the command should display the info about the two containers, described in the Desired State specification. The `influxdb` is expected to be updated with the latest version and in `RUNNING` state and `hello-world` container to be status `EXITED` with version unchanged. The `alpine` container must be removed and not displayed.
```
ID                                    |Name         |Image                               |Status   |Finished At |Exit Code
|-------------------------------------|-------------|------------------------------------|----------------------|---------
7fe6b689-eb76-476d-a730-c2f422d6e8ea  |influxdb     |docker.io/library/influxdb:latest   |Running  |            |0
c36523d7-8d17-4255-ae0c-37f11003f658  |hello-world  |docker.io/library/hello-world:latest|Exited   |            |0
```

### Remove all containers

To remove all containers, publish an empty Desired State specification (with empty `components` section):
```json
{
	"activityId": "testActivityId",
	"payload": {
		"domains": [
			{
				"id": "containers",
				"config": [],
				"components": []
			}
		]
	}
}
```

As a final step, execute the command `kanto-cm list` to verify that the containers are actually removed from the Kanto container management.
The expected output is `No found containers.`.
