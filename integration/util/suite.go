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
	"fmt"
	"net/http"
	"testing"

	"github.com/caarlos0/env/v6"

	"github.com/eclipse/ditto-clients-golang"
	"github.com/eclipse/ditto-clients-golang/model"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/stretchr/testify/require"
)

const (
	// FeatureURLTemplate TBD
	FeatureURLTemplate = "%s/features/%s"
	// FeatureOperationURLTemplate TBD
	FeatureOperationURLTemplate = "%s/inbox/messages/%s"
)

var (
	eventTopicTemplate          string
	liveMessageTopicTemplate    string
	featurePropertyPathTemplate string
)

// SuiteInitializer is testify Suite initialization helper
type SuiteInitializer struct {
	Cfg *TestConfiguration

	ThingCfg *ThingConfiguration

	DittoClient *ditto.Client
	MQTTClient  MQTT.Client

	ThingURL   string
	FeatureURL string
}

// Setup establishes connections to the local MQTT broker and Ditto
func (suite *SuiteInitializer) Setup(t *testing.T, featureID string) {
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
		mqttClient.Disconnect(uint(cfg.MqttQuiesceMs))
		require.NoError(t, err, "initialize ditto client")
	}

	suite.Cfg = cfg
	suite.DittoClient = dittoClient
	suite.MQTTClient = mqttClient
	suite.ThingCfg, err = GetThingConfiguration(cfg, mqttClient)
	if err != nil {
		defer mqttClient.Disconnect(uint(suite.Cfg.MqttQuiesceMs))
	}
	require.NoError(t, err, "get thing configuration")
	suite.InitURLsAndTemplates(suite.ThingCfg.DeviceID, featureID)
}

// TearDown closes all connections
func (suite *SuiteInitializer) TearDown() {
	suite.DittoClient.Disconnect()
	suite.MQTTClient.Disconnect(uint(suite.Cfg.MqttQuiesceMs))
}

// ExecCommand TBD
func (suite *SuiteInitializer) ExecCommand(cfg *TestConfiguration, featureURL string, command string, params interface{}) error {
	url := fmt.Sprintf(FeatureOperationURLTemplate, featureURL, command)
	_, err := SendDigitalTwinRequest(cfg, http.MethodPost, url, params)
	return err
}

// InitURLsAndTemplates TBD
func (suite *SuiteInitializer) InitURLsAndTemplates(thingID string, featureID string) {
	if len(thingID) == 0 {
		thingID = suite.ThingCfg.DeviceID
	}
	suite.ThingURL = GetDigitalTwinURLForThingID(suite.Cfg.DigitalTwinAPIAddress, thingID)
	suite.FeatureURL = suite.GetFeatureURL(featureID)

	thingIDWithNamespace := model.NewNamespacedIDFrom(thingID)
	featurePropertyPathTemplate = fmt.Sprintf("/features/%s/properties/%%s", featureID)
	eventTopicTemplate = fmt.Sprintf("%s/%s/things/twin/events/%%s", thingIDWithNamespace.Namespace, thingIDWithNamespace.Name)
	liveMessageTopicTemplate = fmt.Sprintf("%s/%s/things/live/messages/%%s", thingIDWithNamespace.Namespace, thingIDWithNamespace.Name)
}

// GetFeatureURL TBD
func (suite *SuiteInitializer) GetFeatureURL(featureID string) string {
	return fmt.Sprintf(FeatureURLTemplate, suite.ThingURL, featureID)
}

// GetPropertyPath TBD
func GetPropertyPath(featureID string, name string) string {
	return fmt.Sprintf(featurePropertyPathTemplate, name)
}

// GetEventTopic TBD
func GetEventTopic(action string) string {
	return fmt.Sprintf(eventTopicTemplate, action)
}

// GetLiveMessageTopic TBD
func GetLiveMessageTopic(action string) string {
	return fmt.Sprintf(liveMessageTopicTemplate, action)
}
