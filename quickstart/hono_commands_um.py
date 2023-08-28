# Copyright (c) 2023 Contributors to the Eclipse Foundation
#
# See the NOTICE file(s) distributed with this work for additional
# information regarding copyright ownership.
#
# This program and the accompanying materials are made available under the
# terms of the Eclipse Public License 2.0 which is available at
# https://www.eclipse.org/legal/epl-2.0, or the Apache License, Version 2.0
# which is available at https://www.apache.org/licenses/LICENSE-2.0.
#
# SPDX-License-Identifier: EPL-2.0 OR Apache-2.0

import getopt
import json
import os
import signal
import sys
import threading
import time
import uuid

from string import Template
from proton import Message
from proton.handlers import MessagingHandler
from proton.reactor import Container

alpine_container_template = """
{
"id": "alpine",
"version": "latest",
"config": [
        {
            "key": "image",
            "value": "docker.io/library/alpine:latest"
        }
    ]
},"""

containers_desired_state = Template("""
{
    "topic": "$namespace/$name/things/live/messages/apply",
    "headers": {
        "content-type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": true
    },
    "path": "/features/UpdateManager/inbox/messages/apply",
    "value": {
    "activityId": "$activity_id",
    "desiredState": {
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
                                "value": "docker.io/library/influxdb:$influxdb_version"
                            }
                        ]
                    },
                    $alpine_container
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
}
""")

containers_desired_state_clean_up = Template("""
{
    "topic": "$namespace/$name/things/live/messages/apply",
    "headers": {
        "content-type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": true
    },
    "path": "/features/UpdateManager/inbox/messages/apply",
    "value": {
    "activityId": "$activity_id",
    "desiredState": {
        "domains": [
            {
                "id": "containers",
                "config": [],
                "components": []
            }
        ]
    }
}
}
""")

um_refresh_state = Template("""
{
    "topic": "$namespace/$name/things/live/messages/apply",
    "headers": {
        "content-type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": true
    },
    "path": "/features/UpdateManager/inbox/messages/refresh",
    "value": {
        "activityId": "$activity_id",
    }
}
""")


class CommandResponsesHandler(MessagingHandler):
    def __init__(self, server, address):
        super(CommandResponsesHandler, self).__init__()
        self.server = server
        self.address = address

    def on_start(self, event):
        conn = event.container.connect(self.server, user="consumer@HONO", password="verysecret")
        event.container.create_receiver(conn, self.address)
        print('[connected]')

    def on_message(self, event):
        print('[got response]')
        response = json.loads(event.message.body)
        print(json.dumps(response, indent=2))
        if response["status"] == 204:
            print('[ok]', "um")
        else:
            print('[error]')
        event.receiver.close()
        event.connection.close()

    def on_connection_closed(self, event):
        print('[connection closed]')
        os.kill(os.getpid(), signal.SIGINT)


class CommandsInvoker(MessagingHandler):
    def __init__(self, server, address):
        super(CommandsInvoker, self).__init__()
        self.server = server
        self.address = address

    def on_start(self, event):
        conn = event.container.connect(self.server, sasl_enabled=True, allowed_mechs="PLAIN", allow_insecure_mechs=True,
                                       user="consumer@HONO", password="verysecret")
        event.container.create_sender(conn, self.address)

    def on_sendable(self, event):
        print('[sending command]')
        correlation_id = str(uuid.uuid4())
        namespaced_id = device_id.split(':', 1)
        activity_id = str(uuid.uuid4())

        influxdb_version = "1.8.4"
        alpine_container = alpine_container_template
        if operation == "update":
            influxdb_version = "1.8.5"
            alpine_container = ""
        if operation == "clean":
            payload = containers_desired_state_clean_up.substitute(namespace=namespaced_id[0],
                                                                   name=namespaced_id[1],
                                                                   correlation_id=correlation_id,
                                                                   activity_id=activity_id)
        else:
            payload = containers_desired_state.substitute(namespace=namespaced_id[0], name=namespaced_id[1],
                                                          correlation_id=correlation_id,
                                                          influxdb_version=influxdb_version,
                                                          alpine_container=alpine_container,
                                                          activity_id=activity_id)
        print(json.dumps(json.loads(payload), indent=2))
        msg = Message(body=payload, address='{}/{}'.format(address, device_id),
                      content_type="application/json",
                      subject="um", reply_to=reply_to_address, correlation_id=correlation_id, id=str(uuid.uuid4()))
        event.sender.send(msg)
        event.sender.close()
        event.connection.close()
        print('[sent]')


# Parse command line args
options, reminder = getopt.getopt(sys.argv[1:], 't:d:o:')
opts_dict = dict(options)
tenant_id = os.environ.get("TENANT") or opts_dict['-t']
device_id = os.environ.get("DEVICE_ID") or opts_dict['-d']
operation = opts_dict['-o']

# AMQP global configurations
uri = 'amqps://hono.eclipseprojects.io:15671'
address = 'command/{}'.format(tenant_id)
reply_to_address = 'command_response/{}/replies'.format(tenant_id)
print('[starting] demo update manager app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))

# Create command invoker and handler
response_handler = Container(CommandResponsesHandler(uri, reply_to_address))
commands_invoker = Container(CommandsInvoker(uri, address))
thread = threading.Thread(target=lambda: response_handler.run(), daemon=True)
thread.start()
# Give it some time to link
time.sleep(2)
# Send the command
commands_invoker.run()


def handler(signum, frame):
    print('[stopping] demo update manager app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))
    response_handler.stop()
    thread.join(timeout=5)
    print('[stopped]')
    exit(0)


signal.signal(signal.SIGINT, handler)
while True:
    pass
