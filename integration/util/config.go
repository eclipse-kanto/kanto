// Copyright (c) 2022 Contributors to the Eclipse Foundation
//
// See the NOTICE file(s) distributed with this work for additional
// information regarding copyright ownership.
//
// This program and the accompanying materials are made available under the
// terms of the Eclipse Public License 2.0 which is available at
// https://www.eclipse.org/legal/epl-2.0, or the Apache License, Version 2.0
// which is available at https://www.apache.org/licenses/LICENSE-2.0.
//
// SPDX-License-Identifier: EPL-2.0 OR Apache-2.0

package util

import (
	"time"
)

// TestConfiguration is a common IT configuration
type TestConfiguration struct {
	LocalBroker              string `env:"LOCAL_BROKER" envDefault:"tcp://localhost:1883"`
	MqttQuiesceMs            int    `env:"MQTT_QUIESCE_MS" envDefault:"500"`
	MqttAcknowledgeTimeoutMs int    `env:"MQTT_ACKNOWLEDGE_TIMEOUT_MS" envDefault:"3000"`
	MqttConnectMs            int    `env:"MQTT_CONNECT_TIMEOUT_MS" envDefault:"30000"`

	DigitalTwinAPIAddress  string `env:"DIGITAL_TWIN_API_ADDRESS"`
	DigitalTwinAPIUsername string `env:"DIGITAL_TWIN_API_USERNAME" envDefault:"ditto"`
	DigitalTwinAPIPassword string `env:"DIGITAL_TWIN_API_PASSWORD" envDefault:"ditto"`

	WsEventTimeoutMs int `env:"WS_EVENT_TIMEOUT_MS" envDefault:"30000"`
}

// MillisToDuration converts milliseconds to Duration
func MillisToDuration(millis int) time.Duration {
	return time.Duration(millis) * time.Millisecond
}
