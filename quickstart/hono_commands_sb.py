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

import base64
import getopt
import json
import os
import signal
import sys
import threading
import time
import uuid
import requests
import hashlib

from string import Template
from proton import Message
from proton._reactor import AtLeastOnce
from proton.handlers import MessagingHandler
from proton.reactor import Container

ditto_live_inbox_msg_template = Template("""
{
    "topic": "$namespace/$name/things/live/messages/response",
    "headers": {
        "Content-Type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": true
    },
    "path": "/features/Bootstrapping/inbox/messages/response",
    "value":
    {
        "requestId": "$request_id",
        "hash": "$hash",
        "chunk": "$chunk"
    }
}
""")

config_json_template = Template("""{
    "logFile": "/var/log/suite-connector/suite-connector.log",
    "address": "mqtts://hono.eclipseprojects.io:8883",
    "caCert": "/etc/suite-connector/iothub.crt",
    "tenantId": "$tenant_id",
    "deviceId": "$device_id",
    "authId": "$auth_id",
    "password": "$password"
}""")

update_credentials_template = Template("""
[{
  "type": "hashed-password",
  "auth-id": "$auth_id",
  "secrets": [{
      "pwd-plain": "$pwd"
  }]
}]
""")


class CommandResponsesHandler(MessagingHandler):
    def __init__(self, server, address):
        super(CommandResponsesHandler, self).__init__()
        self.server = server
        self.address = address

    def on_start(self, event):
        conn = event.container.connect(self.server, user="consumer@HONO", password="verysecret")
        event.container.create_receiver(conn, self.address)
        print('[command handler connected]')

    def on_message(self, event):
        print('[got response]')
        response = json.loads(event.message.body)
        print(json.dumps(response, indent=2))
        if response["status"] == 204:
            print('[ok]', "sb")
        else:
            print('[error]')
            delete_res()
            event.receiver.close()
            event.connection.close()
            os.kill(os.getpid(), signal.SIGINT)
            exit(1)
        event.receiver.close()
        event.connection.close()

    def on_connection_closed(self, event):
        print("[command handler connection closed]")


class CommandsInvoker(MessagingHandler):
    def __init__(self, server, address, request_id):
        super(CommandsInvoker, self).__init__()
        self.server = server
        self.address = address
        self.request_id = request_id

    def on_start(self, event):
        conn = event.container.connect(self.server, sasl_enabled=True, allowed_mechs="PLAIN", allow_insecure_mechs=True,
                                       user="consumer@HONO", password="verysecret")
        event.container.create_sender(conn, self.address)

    def on_sendable(self, event):
        print('[provisioning tenant]')
        headers = {
            'Content-Type': 'application/json',
        }
        tenant_headers = {
            'accept': 'application/json',
            'Content-Type': 'application/json',
        }
        response = requests.post(
            '{}/v1/tenants/{}'.format(hono_ep, new_tenant_id), headers=tenant_headers,
            data='{"ext":{"messaging-type": "amqp"}}')

        status_code = response.status_code
        print("Creating tenant [{}]".format(status_code))
        if status_code < 200 or status_code > 300:
            print("Error creating tenant [{}] with code [{}]".format(new_tenant_id, status_code))
            event.sender.close()
            event.connection.close()
            os.kill(os.getpid(), signal.SIGINT)
            exit(1)

        response = requests.post(
            '{}/v1/devices/{}/{}'.format(hono_ep, new_tenant_id, new_device_id),
            headers=headers,
            data='{"authorities":["auto-provisioning-enabled"]}')
        status_code = response.status_code
        print("Creating device [{}]".format(status_code))
        if status_code < 200 or status_code > 300:
            print("Error creating device [{}] with code [{}]".format(new_device_id, status_code))
            event.sender.close()
            event.connection.close()
            os.kill(os.getpid(), signal.SIGINT)
            exit(1)

        response = requests.put(
            '{}/v1/credentials/{}/{}'.format(hono_ep, new_tenant_id, new_device_id),
            headers=headers,
            data=update_credentials_template.substitute(auth_id=new_authentication_id, pwd=new_password))
        status_code = response.status_code
        print("Update credentials [{}]".format(status_code))
        if status_code < 200 or status_code > 300:
            print("Error update credentials for device [{}] with code [{}]".format(new_device_id, status_code))
            event.sender.close()
            event.connection.close()
            os.kill(os.getpid(), signal.SIGINT)
            exit(1)

        print('[sending command]')
        correlation_id = str(uuid.uuid4())
        namespaced_id = device_id.split(':', 1)
        chunk = config_json_template.substitute(tenant_id=new_tenant_id, device_id=new_device_id,
                                                auth_id=new_authentication_id, password=new_password)
        chunk_hash = hashlib.sha256(chunk.encode())
        chunk = base64.b64encode(chunk.encode()).decode('ascii')
        payload = ditto_live_inbox_msg_template.substitute(namespace=namespaced_id[0], name=namespaced_id[1],
                                                           chunk=chunk, hash=chunk_hash.hexdigest(),
                                                           request_id=self.request_id, correlation_id=correlation_id)

        print(payload)
        msg = Message(body=payload, address='{}/{}'.format(address, device_id),
                      content_type="application/json",
                      subject="sb", reply_to=reply_to_address, correlation_id=correlation_id, id=str(uuid.uuid4()))

        event.sender.send(msg)
        event.sender.close()
        event.connection.close()
        print('[sent]')


class EventsHandler(MessagingHandler):
    def __init__(self, server, receiver_address, tenant_id):
        super(EventsHandler, self).__init__()
        self.server = server
        self.receiver_address = receiver_address
        self.tenant_id = tenant_id

    def on_start(self, event):
        conn = event.container.connect(self.server, user="consumer@HONO", password="verysecret")
        event.container.create_receiver(conn, source=self.receiver_address, options=[AtLeastOnce()])
        print('[events handler connected for tenant {}]'.format(self.tenant_id))

    def on_message(self, event):
        print('[events handler on message for tenant {}]'.format(self.tenant_id))
        if event.message.body is not None:
            body = json.loads(event.message.body)
            if body["topic"].split("/")[-1] == "request":
                print('[request event received]')
                print(json.dumps(body, indent=2))
                request_id = body["value"]["requestId"]
                Container(CommandsInvoker(uri, address, request_id=request_id)).run()
                event.receiver.close()
                event.connection.close()

            elif body["topic"].split("/")[0] == new_tenant_id:
                print('[event received from new tenant [{}]]'.format(new_tenant_id))
                delete_res()

                event.receiver.close()
                event.connection.close()

    def on_connection_closed(self, event):
        print("[events handler connection for tenant {} closed]".format(self.tenant_id))
        if self.tenant_id == new_tenant_id:
            os.kill(os.getpid(), signal.SIGINT)


# Parse command line args
options, reminder = getopt.getopt(sys.argv[1:], 't:d:p:')
opts_dict = dict(options)
tenant_id = opts_dict['-t'] or os.environ.get("TENANT")
device_id = opts_dict['-d'] or os.environ.get("DEVICE_ID")
new_password = opts_dict['-p'] or os.environ.get("PASSWORD")

hono_ep = 'https://hono.eclipseprojects.io:28443'

pref = "FromBootstrapping"
new_tenant_id = pref + tenant_id
new_device = device_id.removeprefix("{}:".format(tenant_id)) + pref
new_device_id = '{}:{}'.format(new_tenant_id, new_device)
new_authentication_id = new_device_id.replace(":","_")

# AMQP global configurations
uri = 'amqps://hono.eclipseprojects.io:15671'
address = 'command/{}'.format(tenant_id)
reply_to_address = 'command_response/{}/replies'.format(tenant_id)
event_address = 'event/{}'.format(tenant_id)
new_tenant_event_address = 'event/{}'.format(new_tenant_id)

print('[starting] demo suite bootstrapping app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))

# Create command invoker and handler
response_handler = Container(CommandResponsesHandler(uri, reply_to_address))
events_handler = Container(EventsHandler(uri, event_address, tenant_id))
new_tenant_events_handler = Container(EventsHandler(uri, new_tenant_event_address, new_tenant_id))

response_thread = threading.Thread(target=lambda: response_handler.run(), daemon=True)
response_thread.start()
events_thread = threading.Thread(target=lambda: events_handler.run(), daemon=True)
events_thread.start()
new_tenant_events_thread = threading.Thread(target=lambda: new_tenant_events_handler.run(), daemon=True)
# Give it some time to link
time.sleep(2)
new_tenant_events_thread.start()

def delete_res():
    print('[remove new created resources]')
    headers = {
        'Content-Type': 'application/json',
    }
    response = requests.delete(
        '{}/v1/devices/{}/{}'.format(hono_ep, new_tenant_id, new_device_id), headers=headers)

    print("Deleting device [{}] with code [{}]".format(new_device_id, response.status_code))
    response = requests.delete(
        '{}/v1/credentials/{}/{}'.format(hono_ep, new_tenant_id, new_device_id), headers=headers)
    print("Delete credentials [{}]".format(response.status_code))
    response = requests.delete(
        '{}/v1/tenants/{}'.format(hono_ep, new_tenant_id), headers=headers)
    print("Deleting tenant [{}] with code[{}]".format(new_tenant_id, response.status_code))

def handler(signum, frame):
    print('[stopping] demo suite bootstrapping app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id,
                                                                                                uri))
    response_thread.join(timeout=5)
    events_thread.join(timeout=5)
    new_tenant_events_thread.join(timeout=5)
    print('[stopped]')
    exit(0)


signal.signal(signal.SIGINT, handler)
while True:
    pass
