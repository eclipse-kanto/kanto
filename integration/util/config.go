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
	"encoding/json"
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// TestConfiguration is common IT configuration
type TestConfiguration struct {
	LocalBroker              string `env:"LOCAL_BROKER" envDefault:"tcp://localhost:1883"`
	MqttQuiesceMs            int    `env:"MQTT_QUIESCE_MS" envDefault:"500"`
	MqttAcknowledgeTimeoutMs int    `env:"MQTT_QUIESCE_MS" envDefault:"3000"`
	MqttStatusTimeoutMs      int    `env:"MQTT_STATUS_TIMEOUT_MS" envDefault:"10000"`

	DigitalTwinAPIAddress  string `env:"DIGITAL_TWIN_API_ADDRESS"`
	DigitalTwinAPIUsername string `env:"DIGITAL_TWIN_API_USERNAME" envDefault:"ditto"`
	DigitalTwinAPIPassword string `env:"DIGITAL_TWIN_API_PASSWORD" envDefault:"ditto"`

	WsEventTimeoutMs int `env:"WS_EVENT_TIMEOUT_MS" envDefault:"30000"`
}

// ThingConfiguration is thing configuration info for Edge
type ThingConfiguration struct {
	DeviceID string `json:"deviceId"`
	TenantID string `json:"tenantId"`
	PolicyID string `json:"policyId"`
}

// GetThingConfiguration retrieves ThingConfig using specified client
func GetThingConfiguration(mqttClient MQTT.Client) (*ThingConfiguration, error) {
	type result struct {
		cfg *ThingConfiguration
		err error
	}

	ch := make(chan result)

	if token := mqttClient.Subscribe("edge/thing/response", 1, func(client MQTT.Client, message MQTT.Message) {
		var cfg ThingConfiguration
		if err := json.Unmarshal(message.Payload(), &cfg); err != nil {
			ch <- result{nil, err}
		}
		ch <- result{&cfg, nil}
	}); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := mqttClient.Publish("edge/thing/request", 1, false, ""); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	timeout := 5 * time.Second
	select {
	case result := <-ch:
		return result.cfg, result.err
	case <-time.After(timeout):
		return nil, fmt.Errorf("thing config not received in %v", timeout)
	}
}
