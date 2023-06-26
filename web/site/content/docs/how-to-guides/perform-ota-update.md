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

* Ensure the Kanto container management & update manager services are up & running

* Ensure the Mosquitto service is up & running:

* Run Mosquitto listener to track the desired state feedback messages by executing:
  ```shell
  mosquitto_sub -t "deviceupdate/desiredstatefeedback" -v
  ```

### Publish the Desired State specification

The starting point of the OTA update process is to publish the example Desired State specification below on the MQTT topic `deviceupdate\desiredstate`:
```
{
	"activityId": "testActivitityId",
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

The update process is organized into multiple phases, which are triggered by sending specific desired state commands (`DOWNLOAD/UPDTE/ACTIVATE`).

In the example scenario, the three images for the three container components will be pulled (if not available in the cache locally), created as containers during the `UPDATING` phase and
started in the `ACTIVATING` phase.

At last, the `CLEANUP` phase will be executed to do some clean up actions after the activation of the components.

### Monitor OTA update progress

During the OTA update, the progress can be tracked and monitored by the Mosquitto listener for the desired state feedback messages, started in the prerequisite section above.

The Update Manager reports to the local Mosquitto broker at a time interval of a second the status of the active update process. For example:
```
{
	"activityId": "testActivitityId",
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
					"version": "latest"
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
7fe6b689-eb76-476d-a730-c2f422d6e8ea  |influxdb     |docker.io/library/influxdb:latest   |Running  |            |0
c36523d7-8d17-4255-ae0c-37f11003f658  |hello-world  |docker.io/library/hello-world:latest|Exited   |            |0
9b99978b-2593-4736-bb52-7a07be4a7ed1  |alpine       |docker.io/library/alpine:latest     |Exited   |            |0
```
