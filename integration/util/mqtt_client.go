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
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// NewMQTTClient creates a new MQTT client and connects it to the broker from the test config
func NewMQTTClient(cfg *TestConfiguration) (MQTT.Client, error) {
	opts := MQTT.NewClientOptions().
		AddBroker(cfg.LocalBroker).
		SetClientID(uuid.New().String()).
		SetKeepAlive(30 * time.Second).
		SetCleanSession(true).
		SetAutoReconnect(true)

	mqttClient := MQTT.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return mqttClient, nil
}

// SendMQTTMessage sends a message to a topic using specified client. The message is serialized to JSON Format.
func SendMQTTMessage(cfg *TestConfiguration, client MQTT.Client, topic string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	token := client.Publish(topic, 1, false, payload)
	timeout := time.Duration(cfg.MqttAcknowledgeTimeoutMs * int(time.Millisecond))
	if !token.WaitTimeout(timeout) {
		return errors.New("timeout")
	}
	return token.Error()
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

	defer mqttClient.Unsubscribe("edge/thing/response")

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
