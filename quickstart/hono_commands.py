# Copyright (c) 2022 Contributors to the Eclipse Foundation
#
# See the NOTICE file(s) distributed with this work for additional
# information regarding copyright ownership.
#
# This program and the accompanying materials are made available under the
# terms of the Eclipse Public License 2.0 which is available at
# http://www.eclipse.org/legal/epl-2.0, or the Apache License, Version 2.0
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
    "topic": "$namespace/$name:edge:containers/things/live/messages/$action",
    "headers": {
        "content-type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": true
    },
    "path": "/features/$feature/inbox/messages/$action",
    "value": $value
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
            print('[ok]', command)
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
        if command == CLI_OPT_RUN_CMD:
            action = 'create'
            feature = 'ContainerFactory'
            value = json.dumps(dict(imageRef=container_img_ref, start=True))
        else:
            action = 'remove'
            feature = 'Container:{}'.format(container_id)
            value = 'true'
        payload = ditto_live_inbox_msg_template.substitute(namespace=namespaced_id[0], name=namespaced_id[1], action=action,
                                                           correlation_id=correlation_id, feature=feature, value=value)
        print(payload)
        msg = Message(body=payload, address='{}/{}:edge:containers'.format(address, device_id), content_type="application/json",
                      subject=command, reply_to=reply_to_address, correlation_id=correlation_id, id=str(uuid.uuid4()))
        event.sender.send(msg)
        event.sender.close()
        event.connection.close()
        print('[sent]')


CLI_OPT_RUN_CMD = "run"
CLI_OPT_RM_CMD = "rm"

# Parse command line args
options, reminder = getopt.getopt(sys.argv[2:], 't:d:', ["img=", "id="])
opts_dict = dict(options)
tenant_id = os.environ.get("TENANT") or opts_dict['-t']
device_id = os.environ.get("DEVICE_ID") or opts_dict['-d']
command = sys.argv[1]
if command == CLI_OPT_RUN_CMD:
    container_img_ref = opts_dict['--img']
elif command == CLI_OPT_RM_CMD:
    container_id = opts_dict['--id']
else:
    print('[error] unsupported command', command)
    exit(1)

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
