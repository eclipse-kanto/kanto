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
