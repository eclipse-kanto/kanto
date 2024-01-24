# Copyright (c) 2024 Contributors to the Eclipse Foundation
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
import time
import uuid
from string import Template
import paho.mqtt.client as mqtt

retrieve_things_template = Template("""
{
    "topic": "_/_/things/twin/commands/retrieve",
    "headers": {
        "correlation-id": "$correlation_id",
        "reply-to": "command/$tenant_id",
        "response-required": true
    },
    "path": "/",
    "value": {
		"thingIds": [
			"$namespace:$name",
			"$namespace:$name:edge:containers"
		]
	}
}
""")

class MQTTClient(mqtt.Client):

    def on_connect(self, mqttc, obj, flags, rc):
        print('[client connected]')
        self.subscribe("command//+/req/#")

    def on_message(self, mqttc, obj, msg):
        print('[got response]')
        payload = str(msg.payload.decode("utf-8"))
        print(json.dumps(json.loads(payload), indent=2))

    def run(self):
        self.connect("localhost")
        self.loop_start()
        time.sleep(2)
        namespaced_id = device_id.split(':', 1)
        payload = retrieve_things_template.substitute(
            namespace=namespaced_id[0], name=namespaced_id[1],
            tenant_id=tenant_id, correlation_id=str(uuid.uuid4()))
        pub_topic = 'event/{}/{}'.format(tenant_id, device_id)
        self.publish(pub_topic, payload)
        time.sleep(1)
        self.loop_stop()

# Parse command line args
options, reminder = getopt.getopt(sys.argv[1:], 't:d:')
opts_dict = dict(options)
tenant_id = opts_dict.get('-t') or os.environ["TENANT"]
device_id = opts_dict.get('-d') or os.environ["DEVICE_ID"]

uri = 'localhost'
print('[starting] demo local digital twins app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))

def handler(signum, frame):
    print('[stopping] demo local digital twins app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))
    mqttc.loop_stop()
    print('[stopped]')
    exit(0)

signal.signal(signal.SIGINT, handler)
mqttc = MQTTClient()
mqttc.run()
print('[stopped] demo local digital twins app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))