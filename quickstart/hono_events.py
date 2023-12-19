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

from proton.handlers import MessagingHandler
from proton.reactor import Container, AtLeastOnce

class EventsHandler(MessagingHandler):
    def __init__(self, server, address):
        super(EventsHandler, self).__init__()
        self.server = server
        self.address = address

    def on_start(self, event):
        conn = event.container.connect(self.server, user="consumer@HONO", password="verysecret")
        event.container.create_receiver(conn, source=self.address, options=[AtLeastOnce()])
        print('[connected]')

    def on_message(self, event):
        print('[event received]')
        if event.message.body is not None:
            print(json.dumps(json.loads(event.message.body), indent=2))
        else:
            print('<empty>')


# Parse command line args
options, reminder = getopt.getopt(sys.argv[1:], 't:')
tenant_id = dict(options).get('-t') or os.environ["TENANT"]

uri = 'amqps://hono.eclipseprojects.io:15671'
address = 'event/{}'.format(tenant_id)

print('[starting] demo events handler for tenant [{}] at [{}]'.format(tenant_id, uri))

events_handler = Container(EventsHandler(uri, address))
thread = threading.Thread(target=lambda: events_handler.run(), daemon=True)
thread.start()


def handler(signum, frame):
    print('[stopping] demo events handler for tenant [{}] at [{}]'.format(tenant_id, uri))
    events_handler.stop()
    thread.join(timeout=5)
    print('[stopped]')
    exit(0)


signal.signal(signal.SIGINT, handler)
while True:
    pass
