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
from collections import deque
from queue import Queue

from dash import Dash, dcc, html
from dash.exceptions import PreventUpdate
from dash.dependencies import Input, Output

from plotly.graph_objs import Scatter, Layout, Figure

from proton import Message
from proton.handlers import MessagingHandler
from proton.reactor import Container, AtLeastOnce

ditto_live_inbox_msg_template = Template("""
{
    "topic": "$namespace/$name/things/live/messages/metrics",
    "headers": {
        "content-type": "application/json",
        "correlation-id": "$correlation_id",
        "response-required": $response_required,
        "timeout": "$timeout"
    },
    "path": "/features/Metrics/inbox/messages/request",
    "value": {
        "frequency": "$frequency"
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
        print('[response handler connected]')

    def on_message(self, event):
        print('[got response]')
        response = json.loads(event.message.body)
        print(json.dumps(response, indent=2))
        if response["status"] == 204:
            print('[ok]', "sm")
        else:
            print('[error]')
            event.receiver.close()
            event.connection.close()

    def on_connection_closed(self, event):
        print("[connection closed]")
        os.kill(os.getpid(), signal.SIGINT)


class CommandsInvoker(MessagingHandler):
    def __init__(self, server, address, frequency, correlation_id):
        super(CommandsInvoker, self).__init__()
        self.server = server
        self.address = address
        self.frequency = frequency
        self.correlation_id = correlation_id

    def on_start(self, event):
        conn = event.container.connect(self.server, sasl_enabled=True, allowed_mechs="PLAIN", allow_insecure_mechs=True,
                                       user="consumer@HONO", password="verysecret")
        event.container.create_sender(conn, self.address)

    def on_sendable(self, event):
        print('[sending command]')
        namespaced_id = device_id.split(':', 1)
        response_required = "false" if self.frequency == "0s" else "true"
        timeout = "0s" if self.frequency == "0s" else "60s"
        payload = ditto_live_inbox_msg_template.substitute(namespace=namespaced_id[0], name=namespaced_id[1],
                                                           correlation_id=self.correlation_id,
                                                           response_required=response_required,
                                                           timeout=timeout, frequency=self.frequency)
        print(payload)
        msg = Message(body=payload, address='{}/{}'.format(self.address, device_id), content_type="application/json",
                      subject="sm", reply_to=reply_to_address, correlation_id=self.correlation_id, id=str(uuid.uuid4()))
        event.sender.send(msg)
        event.sender.close()
        event.connection.close()
        print('[sent]')


class EventsHandler(MessagingHandler):
    def __init__(self, server, receiver_address):
        super(EventsHandler, self).__init__()
        self.server = server
        self.receiver_address = receiver_address
        self.metrics_sender = queue

    def on_start(self, event):
        conn = event.container.connect(self.server, user="consumer@HONO", password="verysecret")
        event.container.create_receiver(conn, source=self.receiver_address, options=[AtLeastOnce()])
        print('[events handler connected]')

    def on_message(self, event):
        if event.message.body is not None:
            body = json.loads(event.message.body)
            if body["topic"].split("/")[-1] == "data" and \
                    body["headers"]["correlation-id"] == start_metrics_correlation_id:
                print('[metrics event received]')
                print(json.dumps(body, indent=2))
                # get timestamp from body
                timestamp = body["value"]["timestamp"]
                # get cpu utilization from body
                cpu = [d for d in body["value"]["snapshot"][0]["measurements"] if d["id"] == "cpu.utilization"]
                cpu = float(cpu[0]["value"])
                # get memory utilization from body
                mem = [d for d in body["value"]["snapshot"][0]["measurements"] if d["id"] == "memory.utilization"]
                mem = float(mem[0]["value"])
                # provide data to the application
                self.metrics_sender.put([timestamp, cpu, mem])


class DashApp:
    app = Dash(__name__, external_stylesheets=['https://codepen.io/chriddyp/pen/bWLwgP.css'])
    time_data = deque(maxlen=10)
    cpu_util_data = deque(maxlen=10)
    mem_util_data = deque(maxlen=10)

    def __init__(self):
        self.metrics_receiver = queue
        # Set up application and callbacks
        self.app.layout = html.Div([
            dcc.Loading(id="loading", parent_style={"margin-top": "20vh"}),
            dcc.Graph(id="live-graph", animate=True, style={"display": "none"}),
            dcc.Interval(id="interval-component", interval=5 * 1000, n_intervals=0)
        ], style={"display": "flex", "justify-content": "center"})
        self.register_callbacks(self.app)

    def register_callbacks(self, app):
        @app.callback([Output("live-graph", "figure"),
                       Output("loading", "parent_style"),
                       Output("live-graph", "style")],
                      [Input("interval-component", "n_intervals")])
        def update_graph_live(n):
            if n is None:
                raise PreventUpdate

            # Get new data from queue
            new_data = self.metrics_receiver.get()
            self.time_data.append(new_data[0])
            self.cpu_util_data.append(new_data[1])
            self.mem_util_data.append(new_data[2])

            # Set up scatters for graph
            cpu_trace = Scatter(
                x=list(self.time_data), y=list(self.cpu_util_data),
                name="CPU Utilization", mode="lines+markers"
            )
            mem_trace = Scatter(
                x=list(self.time_data), y=list(self.mem_util_data),
                name="Memory Utilization", mode="lines+markers"
            )
            graph_data = [cpu_trace, mem_trace]

            # Create graph layout
            layout = Layout(
                title="System Metrics Data",
                xaxis=dict(
                    range=[min(self.time_data), max(self.time_data)],
                    tickmode="array",
                    tickvals=list(self.time_data),
                    ticktext=[time.strftime("%d %b, %H:%M:%S", time.localtime(t)) for t in list(self.time_data)]
                ),
                yaxis=dict(range=[
                    min(min(self.cpu_util_data), min(self.mem_util_data)) - 1,
                    max(max(self.cpu_util_data), max(self.mem_util_data)) + 1
                ]),
                autosize=True
            )

            self.metrics_receiver.task_done()

            fig = Figure(graph_data, layout)
            loading_style = {"display": "none"}
            graph_style = {"display": "block", "width": "90vw", "height": "90vh"}
            return fig, loading_style, graph_style

    def run(self):
        # Get initial data and run application
        initial_data = self.metrics_receiver.get()
        self.time_data.append(initial_data[0])
        self.cpu_util_data.append(initial_data[1])
        self.mem_util_data.append(initial_data[2])
        self.app.run_server(debug=True, use_reloader=False)


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

print('[starting] demo system metrics app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))

# Create pipe for transferring metrics data
queue = Queue()

# Create Dash application
dash_application = DashApp()

# Create event handler and command response handlers
events_handler = Container(EventsHandler(uri, event_address))
events_thread = threading.Thread(target=lambda: events_handler.run(), daemon=True)
response_handler = Container(CommandResponsesHandler(uri, reply_to_address))
response_thread = threading.Thread(target=lambda: response_handler.run(), daemon=True)

# Create start and stop metrics messages
start_metrics_correlation_id = str(uuid.uuid4())
start_metrics_message = Container(CommandsInvoker(uri, command_address, "5s", start_metrics_correlation_id))
stop_metrics_message = Container(CommandsInvoker(uri, command_address, "0s", str(uuid.uuid4())))

# Start threads
events_thread.start()
response_thread.start()
# Give it some time to link
time.sleep(2)

# Send the command
start_metrics_message.run()


def handler(signum, frame):
    print('[stopping] demo system metrics app for tenant [{}], device [{}] at [{}]'.format(tenant_id, device_id, uri))
    # Send message to stop metrics
    stop_metrics_message.run()
    # Stop handlers
    events_handler.stop()
    response_handler.stop()
    # Wait for threads to finish execution
    events_thread.join(timeout=5)
    response_thread.join(timeout=5)
    print('[stopped]')
    exit(0)


# Start the Dash application
signal.signal(signal.SIGINT, handler)
dash_application.run()
