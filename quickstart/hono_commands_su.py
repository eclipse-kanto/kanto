# Copyright (c) 2022 Contributors to the Eclipse Foundation
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

ditto_live_inbox_msg_template = Template("""
{
    "topic": "$namespace/$name/things/live/messages/install",
    "headers": {
        "content-type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": true
    },
    "path": "/features/SoftwareUpdatable/inbox/messages/install",
    "value": {
        "softwareModules": [{
            "softwareModule": {
                "name": "install-hello",
                "version": "1.0.0"
            },
            "artifacts": [{
                "checksums": {
                    "SHA256": "db954c633393c1402f145a60fd58d312f5af96ce49422fcfd6ce42a3c4cceeca"
                },
                "download": {
                    "HTTPS": {
                        "url": "https://github.com/eclipse-kanto/kanto/raw/main/quickstart/install_hello.sh"
                    }
                },
                "filename": "install.sh",
                "size": 544
            }]
        }],
        "correlationId": "$correlation_id"
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
            print('[ok]', "su")
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

        payload = ditto_live_inbox_msg_template.substitute(namespace=namespaced_id[0], name=namespaced_id[1],
                                                           correlation_id=correlation_id)
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
uri = 'amqps://hono.eclipseprojects.io:15671'
address = 'command/{}'.format(tenant_id)
reply_to_address = 'command_response/{}/replies'.format(tenant_id)

print('[starting] demo software update app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))

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
    print('[stopping] demo software update app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))
    response_handler.stop()
    thread.join(timeout=5)
    print('[stopped]')
    exit(0)


signal.signal(signal.SIGINT, handler)
while True:
    pass
