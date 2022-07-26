# Copyright (c) 2022 Contributors to the Eclipse Foundation
#
# See the NOTICE file(s) distributed with this work for additional
# information regarding copyright ownership.
#
# This program and the accompanying materials are made available under the
# terms of the Eclipse Public License 2.0 which is available at
# http://www.eclipse.org/legal/epl-2.0
#
# SPDX-License-Identifier: EPL-2.0

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

ditto_live_inbox_msg_template = Template("""
{
    "topic": "$namespace/$name/things/live/messages/$action",
    "headers": {
        "content-type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": true
    },
    "path": "/features/$feature/inbox/messages/$action",
    "value": $value
}
""")

software_update_action_template = Template("""
{
    "softwareModules": [{
          "softwareModule": {
            "name": "$module_name",
            "version": "$version"
          },
          "artifacts": [{
            "checksums": {
              "SHA256": "$sha256"
            },
            "download": {
              "HTTPS": {
                "url": "$url"
              }
            },
            "filename": "install.sh",
            "size": 	$size
          }]
        }],
    "correlationId": "$correlation_id"
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
            print('[ok]', "su")
        else:
            print('[error]')
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
        action = "install"
        feature = "SoftwareUpdatable"

        value = json.dumps(json.loads(software_update_action_template.substitute(
            module_name="hello", version="2.10", size=30,
            url="https://github.com/eclipse-kanto/kanto/raw/main/quickstart/install_hello.sh",
            sha256="03a105509663680d6d5db0e2b5939a77fab68429fca4dd5d181924a73e82b5d9",
            correlation_id=str(uuid.uuid4()))))

        payload = ditto_live_inbox_msg_template.substitute(namespace=namespaced_id[0], name=namespaced_id[1],
                                                           action=action,
                                                           correlation_id=correlation_id, feature=feature, value=value)
        print(json.dumps(json.loads(payload), indent=2))
        msg = Message(body=payload, address='{}/{}'.format(address, device_id),
                      content_type="application/json",
                      subject="su", reply_to=reply_to_address, correlation_id=correlation_id, id=str(uuid.uuid4()))
        event.sender.send(msg)
        event.sender.close()
        event.connection.close()
        print('[sent]')


# Parse command line args
options, reminder = getopt.getopt(sys.argv[1:], 't:d:')
opts_dict = dict(options)
tenant_id = os.environ.get("TENANT") or opts_dict['-t']
device_id = os.environ.get("DEVICE_ID") or opts_dict['-d']


# AMQP global configurations
uri = 'amqp://hono.eclipseprojects.io:15672'
address = 'command/{}'.format(tenant_id)
reply_to_address = 'command_response/{}/replies'.format(tenant_id)

print('[starting] demo commands app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))

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
    print('[stopping] demo commands app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))
    response_handler.stop()
    thread.join(timeout=5)
    print('[stopped]')
    exit(0)


signal.signal(signal.SIGINT, handler)
while True:
    pass

