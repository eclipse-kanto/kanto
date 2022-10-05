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
# import signal
import sys
# import threading
import time
import uuid
from string import Template
import paho.mqtt.client as mqtt

ditto_receive_thing_template = Template("""
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

class MQTTClientClass(mqtt.Client):

    def on_connect(self, mqttc, obj, flags, rc):
        print('[client connected]')
        self.subscribe("command/#")

    def on_message(self, mqttc, obj, msg):
        print('[got response]')
        payload = str(msg.payload.decode("utf-8"))
        print(json.dumps(json.loads(payload), indent=2))

    def run(self):
        self.connect("localhost")

        rc = 0
        while rc == 0:
            rc = self.loop_start()
            namespaced_id = device_id.split(':', 1)
            payload = ditto_receive_thing_template.substitute(
                namespace=namespaced_id[0], name=namespaced_id[1],
                tenant_id=tenant_id, correlation_id=str(uuid.uuid4()))
            pub_topic = 'e/{}/{}'.format(tenant_id, device_id)
            self.publish(pub_topic, payload)
            time.sleep(1)
            rc = self.loop_stop()
            return rc
            
# Parse command line args
options, reminder = getopt.getopt(sys.argv[1:], 't:d:')
opts_dict = dict(options)
tenant_id = os.environ.get("TENANT") or opts_dict['-t']
device_id = os.environ.get("DEVICE_ID") or opts_dict['-d']

uri = 'localhost'
print('[starting] local-digital-twins commands app for tenant [{}], device [{}] at [{}]'.format(
    tenant_id, device_id, uri))


mqttc = MQTTClientClass()
mqttc.run()
# time.sleep(1)
print('[stopped] local-digital-twins commands app for tenant [{}], device [{}] at [{}]'.format(
        tenant_id, device_id, uri))
exit(0)
