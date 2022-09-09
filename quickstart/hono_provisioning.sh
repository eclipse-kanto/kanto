#!/bin/sh
#
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

# The Hono AMQP endpoint to connect to
export HONO_EP=hono.eclipseprojects.io
# The Hono tenant to be created
export TENANT=
# The identifier of the device on the tenant
# Note! It's important for the ID to follow the convention namespace:name (e.g. edge:device)
export DEVICE_ID=
# The authentication identifier of the device
export AUTH_ID=
# A password for the device to authenticate with
export PWD=

curl -i -X POST http://$HONO_EP:28080/v1/tenants/$TENANT -H  "content-type: application/json"  --data-binary '{"ext": {"messaging-type": "amqp"}}'
curl -i -X POST http://$HONO_EP:28080/v1/devices/$TENANT/$DEVICE_ID -H  "content-type: application/json" --data-binary '{"authorities":["auto-provisioning-enabled"]}'
curl -i -X PUT -H "content-type: application/json" --data-binary '[{
  "type": "hashed-password",
  "auth-id": "'$AUTH_ID'",
  "secrets": [{
      "pwd-plain": "'$PWD'"
  }]
}]' http://$HONO_EP:28080/v1/credentials/$TENANT/$DEVICE_ID
