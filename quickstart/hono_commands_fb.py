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
from proton.reactor import Container, AtLeastOnce

ditto_live_inbox_msg_template = Template("""
{
    "topic": "$namespace/$name/things/live/messages/$action",
    "headers": {
        "content-type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": true
    },
        "path": "/features/BackupAndRestore/inbox/messages/$action",
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
            print('[ok]', "fb")
        else:
            print('[error]')
            event.receiver.close()
            event.connection.close()

    def on_connection_closed(self, event):
        print('[connection closed]')
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
        if self.action == CLI_OPT_BACKUP_CMD:
            # backup command
            value = json.dumps(dict(correlationId=self.correlation_id))
        else:
            # start and restore command
            opts = {"https.url": "http://{}:8080/data.zip".format(host)}
            value = json.dumps(dict(correlationId=self.correlation_id, options=opts))

        payload = ditto_live_inbox_msg_template.substitute(namespace=namespaced_id[0], name=namespaced_id[1],
                                                           action=self.action, correlation_id=correlation_id,
                                                           value=value)
        print(payload)
        msg = Message(body=payload, address='{}/{}'.format(self.address, device_id), content_type="application/json",
                      subject=self.action, reply_to=reply_to_address, correlation_id=correlation_id,
                      id=str(uuid.uuid4()))
        event.sender.send(msg)
        event.sender.close()
        event.connection.close()
        print('[sent]')


class EventsHandler(MessagingHandler):
    def __init__(self, server, receiver_address, correlation_id):
        super(EventsHandler, self).__init__()
        self.server = server
        self.receiver_address = receiver_address
        self.correlation_id = correlation_id

    def on_start(self, event):
        conn = event.container.connect(self.server, user="consumer@HONO", password="verysecret")
        event.container.create_receiver(conn, source=self.receiver_address, options=[AtLeastOnce()])
        print('[events handler connected]')

    def on_message(self, event):
        if event.message.body is not None:
            body = json.loads(event.message.body)
            if body["topic"].split("/")[-1] == "request":
                # received request event for uploading the file
                print('[request event received]')
                print(json.dumps(body, indent=2))
                request_correlation_id = body["value"]["correlationId"]
                Container(CommandsInvoker(uri, command_address, "start", correlation_id=request_correlation_id)).run()
            elif body["topic"].split("/")[-1] == "modify" and body["path"].split("/")[-1] == "lastUpload":
                # upload feature updated
                if body["value"]["correlationId"] == self.correlation_id:
                    print('[last upload event received]')
                    print(json.dumps(body, indent=2))
                    if body["value"]["state"] == "SUCCESS":
                        print('[successful backup]')
                        event.receiver.close()
                        event.connection.close()
                    elif body["value"]["state"] == "FAILED":
                        print('[failed backup]')
                        event.receiver.close()
                        event.connection.close()
            elif body["topic"].split("/")[-1] == "modify" and body["path"].split("/")[-1] == "lastOperation":
                # operation (backup/restore) feature updated
                if body["value"]["correlationId"] == self.correlation_id:
                    print('[last operation event received]')
                    print(json.dumps(body, indent=2))
                    if body["value"]["state"] == "RESTORE_FINISHED":
                        print('[successful restore]')
                        event.receiver.close()
                        event.connection.close()
                    elif body["value"]["state"] == "RESTORE_FAILED":
                        print('[failed restore]')
                        event.receiver.close()
                        event.connection.close()

    def on_connection_closed(self, event):
        print('[connection closed]')
        os.kill(os.getpid(), signal.SIGINT)


CLI_OPT_BACKUP_CMD = "backup"
CLI_OPT_RESTORE_CMD = "restore"

# Parse command line args
options, reminder = getopt.getopt(sys.argv[2:], 't:d:h:')
opts_dict = dict(options)
tenant_id = opts_dict.get('-t') or os.environ["TENANT"]
device_id = opts_dict.get('-d') or os.environ["DEVICE_ID"]
host = opts_dict.get('-h') or os.environ["HOST"]
command = sys.argv[1]
if command == CLI_OPT_BACKUP_CMD:
    action = "backup"
elif command == CLI_OPT_RESTORE_CMD:
    action = "restore"
else:
    print('[error] unsupported command', command)
    exit(1)

# AMQP global configurations
uri = 'amqps://hono.eclipseprojects.io:15671'
command_address = 'command/{}'.format(tenant_id)
event_address = 'event/{}'.format(tenant_id)
reply_to_address = 'command_response/{}/replies'.format(tenant_id)

print('[starting] demo file backup and restore app for tenant [{}], device [{}] at [{}]'
      .format(tenant_id, device_id, uri))

# Create command invoker and handler
upload_correlation_id = "demo.backup.and.restore" + str(uuid.uuid4())
events_handler = Container(EventsHandler(uri, event_address, correlation_id=upload_correlation_id))
response_handler = Container(CommandResponsesHandler(uri, reply_to_address))
commands_invoker = Container(CommandsInvoker(uri, command_address, action, correlation_id=upload_correlation_id))

events_thread = threading.Thread(target=lambda: events_handler.run(), daemon=True)
events_thread.start()
response_thread = threading.Thread(target=lambda: response_handler.run(), daemon=True)
response_thread.start()
# Give it some time to link
time.sleep(2)
# Send the command
commands_invoker.run()


def handler(signum, frame):
    print('[stopping] demo file backup and restore app for tenant [{}], device [{}] at [{}]'
          .format(tenant_id, device_id, uri))
    events_handler.stop()
    response_handler.stop()
    events_thread.join(timeout=5)
    response_thread.join(timeout=5)
    print('[stopped]')
    exit(0)


signal.signal(signal.SIGINT, handler)
while True:
    pass
