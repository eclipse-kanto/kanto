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
from proton._reactor import AtLeastOnce
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
    "path": "/features/AutoUploadable/inbox/messages/$action",
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
        print('[response handler connected]')

    def on_message(self, event):
        print('[got response]')
        response = json.loads(event.message.body)
        print(json.dumps(response, indent=2))
        if response["status"] == 204:
            print('[ok]')
            if response["topic"].split("/")[-1] == "start":
                os.kill(os.getpid(), signal.SIGINT)
        else:
            print('[error]')
            os.kill(os.getpid(), signal.SIGINT)


class CommandsInvoker(MessagingHandler):
    def __init__(self, server, address, action, correlation_id=None):
        super(CommandsInvoker, self).__init__()
        self.server = server
        self.address = address
        self.action = action
        self.correlation_id = correlation_id

    def on_start(self, event):
        conn = event.container.connect(self.server, sasl_enabled=True, allowed_mechs="PLAIN", allow_insecure_mechs=True,
                                       user="consumer@HONO", password="verysecret")
        event.container.create_sender(conn, self.address)

    def on_sendable(self, event):
        print('[sending command]')
        correlation_id = str(uuid.uuid4())
        namespaced_id = device_id.split(':', 1)
        if self.action == "trigger":
            value = json.dumps(dict(correlation_id=self.correlation_id))
        else:
            upload_options = {"https.url": "http://localhost:8080/suite-connector.log"}
            value = json.dumps(dict(correlationId=self.correlation_id, options=upload_options))

        payload = ditto_live_inbox_msg_template.substitute(namespace=namespaced_id[0], name=namespaced_id[1],
                                                           action=self.action, correlation_id=correlation_id,
                                                           value=value)
        print(payload)
        msg = Message(body=payload, address='{}/{}'.format(self.address, device_id), content_type="application/json",
                      subject="fu", reply_to=reply_to_address, correlation_id=correlation_id, id=str(uuid.uuid4()))
        event.sender.send(msg)
        event.sender.close()
        event.connection.close()
        print('[sent]')


class EventsHandler(MessagingHandler):
    def __init__(self, server, receiver_address):
        super(EventsHandler, self).__init__()
        self.server = server
        self.receiver_address = receiver_address

    def on_start(self, event):
        conn = event.container.connect(self.server, user="consumer@HONO", password="verysecret")
        event.container.create_receiver(conn, source=self.receiver_address, options=[AtLeastOnce()])
        print('[events handler connected]')

    def on_message(self, event):
        print('[event received]')
        if event.message.body is not None:
            body = json.loads(event.message.body)
            print(json.dumps(body, indent=2))
            if body["topic"].split("/")[-1] == "request":
                correlation_id = body["value"]["correlationId"]
                Container(CommandsInvoker(uri, command_address, "start", correlation_id=correlation_id)).run()
        else:
            print('<empty>')


# Parse command line args
options, reminder = getopt.getopt(sys.argv[1:], 't:d:')
opts_dict = dict(options)
tenant_id = os.environ.get("TENANT") or opts_dict['-t']
device_id = os.environ.get("DEVICE_ID") or opts_dict['-d']

# AMQP global configurations
uri = 'amqp://hono.eclipseprojects.io:15672'
command_address = 'command/{}'.format(tenant_id)
event_address = 'event/{}'.format(tenant_id)
reply_to_address = 'command_response/{}/replies'.format(tenant_id)

print('[starting] file upload commands app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))

# Create command invoker and handler
events_handler = Container(EventsHandler(uri, event_address))
response_handler = Container(CommandResponsesHandler(uri, reply_to_address))
commands_invoker = Container(CommandsInvoker(uri, command_address, "trigger", correlation_id="demo.upload"))

events_thread = threading.Thread(target=lambda: events_handler.run(), daemon=True)
events_thread.start()
response_thread = threading.Thread(target=lambda: response_handler.run(), daemon=True)
response_thread.start()
# Give it some time to link
time.sleep(2)
# Send the command
commands_invoker.run()


def handler(signum, frame):
    print('[stopping] file upload commands app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))
    events_handler.stop()
    response_handler.stop()
    events_thread.join(timeout=5)
    response_thread.join(timeout=5)
    print('[stopped]')
    exit(0)


signal.signal(signal.SIGINT, handler)
while True:
    pass
