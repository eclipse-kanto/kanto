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
	"testing"

	"github.com/caarlos0/env/v6"

	"github.com/eclipse/ditto-clients-golang"
	"github.com/eclipse/ditto-clients-golang/model"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/stretchr/testify/require"
)

// SuiteInitializer is testify Suite initialization helper
type SuiteInitializer struct {
	Cfg *TestConfiguration

	ThingCfg *ThingConfiguration

	DittoClient *ditto.Client
	MQTTClient  MQTT.Client

	thingURL                    string
	featureURL                  string
	featureOperationURLTemplate string

	featurePropertyPathTemplate string

	eventTopicTemplate       string
	liveMessageTopicTemplate string
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
	require.NoError(t, err, "get thing configuration")
	suite.FormatURLsAndTemplates(featureID)
}

// TearDown closes all connections
func (suite *SuiteInitializer) TearDown() {
	suite.DittoClient.Disconnect()
	suite.MQTTClient.Disconnect(uint(suite.Cfg.MqttQuiesceMs))
}

// FormatURLsAndTemplates TBD
func (suite *SuiteInitializer) FormatURLsAndTemplates(featureID string) {
	suite.thingURL = GetDigitalTwinAPIURLForThingID(suite.Cfg.DigitalTwinAPIAddress, suite.ThingCfg.DeviceID)
	suite.featureURL = fmt.Sprintf("%s/features/%s", suite.thingURL, featureID)
	suite.featurePropertyPathTemplate = fmt.Sprintf("/features/%s/properties/%%s", featureID)

	namespaceID := model.NewNamespacedIDFrom(suite.ThingCfg.DeviceID)
	suite.eventTopicTemplate = fmt.Sprintf("%s/%s/things/twin/events/%%s", namespaceID.Namespace, namespaceID.Name)
	suite.liveMessageTopicTemplate = fmt.Sprintf("%s/%s/things/live/messages/%%s", namespaceID.Namespace, namespaceID.Name)
	suite.featureOperationURLTemplate = fmt.Sprintf("%s/inbox/messages/%%s", suite.featureURL)
}

// GetFeatureURL TBD
func (suite *SuiteInitializer) GetFeatureURL() string {
	return suite.featureURL
}

// GetPropertyPath TBD
func (suite *SuiteInitializer) GetPropertyPath(propertyName string) string {
	return fmt.Sprintf(suite.featurePropertyPathTemplate, propertyName)
}

// GetEventTopic TBD
func (suite *SuiteInitializer) GetEventTopic(propertyName string) string {
	return fmt.Sprintf(suite.eventTopicTemplate, propertyName)
}

// GetLiveMessageTopic TBD
func (suite *SuiteInitializer) GetLiveMessageTopic(propertyName string) string {
	return fmt.Sprintf(suite.liveMessageTopicTemplate, propertyName)
}
