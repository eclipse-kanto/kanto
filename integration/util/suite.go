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
	"testing"

	"github.com/caarlos0/env/v6"

	"github.com/eclipse/ditto-clients-golang"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/stretchr/testify/require"
)

// SuiteInitializer is testify Suite initialization helper
type SuiteInitializer struct {
	Cfg *TestConfiguration

	ThingCfg *ThingConfiguration

	DittoClient *ditto.Client
	MQTTClient  MQTT.Client
}

// Setup establishes connections to the local MQTT broker and Ditto
func (suite *SuiteInitializer) Setup(t *testing.T) {
	cfg := &TestConfiguration{}

	opts := env.Options{RequiredIfNoDef: true}
	require.NoError(t, env.Parse(cfg, opts), "failed to process environment variables")

	t.Logf("%#v\n", cfg)

	mqttClient, err := NewMQTTClient(cfg)
	require.NoError(t, err, "connect to MQTT broker")

	dittoClient, err := ditto.NewClientMQTT(mqttClient, ditto.NewConfiguration())
	if err == nil {
		err = dittoClient.Connect()
	}

	if err != nil {
		mqttClient.Disconnect(uint(cfg.MQTTQuiesceMS))
		require.NoError(t, err, "initialize ditto client")
	}

	suite.Cfg = cfg
	suite.DittoClient = dittoClient
	suite.MQTTClient = mqttClient
	suite.ThingCfg, err = GetThingConfiguration(cfg, mqttClient)
	if err != nil {
		defer suite.TearDown()
		require.NoError(t, err, "cannot get thing configuration")
	}
}

// TearDown closes all connections
func (suite *SuiteInitializer) TearDown() {
	suite.DittoClient.Disconnect()
	suite.MQTTClient.Disconnect(uint(suite.Cfg.MQTTQuiesceMS))
}
