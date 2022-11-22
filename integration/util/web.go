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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eclipse/ditto-clients-golang/model"
	"github.com/eclipse/ditto-clients-golang/protocol"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

const (
	featureURLTemplate          = "%s/features/%s"
	featureOperationURLTemplate = "%s/inbox/messages/%s"
	thingURLTemplate            = "%s/api/2/things/%s"
	featurePropertyPathTemplate = "/features/%s/properties/%s"
	messagePathTemplate         = "/features/%s/outbox/messages/%s"
	eventTopicTemplate          = "%s/%s/things/twin/events/%s"
	liveMessageTopicTemplate    = "%s/%s/things/live/messages/%s"
)

// SendHTTPFeatureRequest sends а new HTTP request to Ditto REST API
func SendHTTPFeatureRequest(cfg *TestConfiguration, method string, url string, body interface{}) ([]byte, error) {
	var reqBody io.Reader

	if body != nil {
		jsonValue, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonValue)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		correlationID := uuid.New().String()
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("correlation-id", correlationID)
		req.Header.Add("response-required", "true")
	}

	req.SetBasicAuth(cfg.DigitalTwinAPIUsername, cfg.DigitalTwinAPIPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s request failed: %s", method, url, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// NewDigitalTwinWSConnection creates new web socket connection
func NewDigitalTwinWSConnection(cfg *TestConfiguration) (*websocket.Conn, error) {
	wsAddress, err := asWSAddress(cfg.DigitalTwinAPIAddress)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/ws/2", wsAddress)
	wsCfg, err := websocket.NewConfig(url, cfg.DigitalTwinAPIAddress)
	if err != nil {
		return nil, err
	}

	auth := fmt.Sprintf("%s:%s", cfg.DigitalTwinAPIUsername, cfg.DigitalTwinAPIPassword)
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	wsCfg.Header = http.Header{
		"Authorization": {"Basic " + enc},
	}

	return websocket.DialConfig(wsCfg)
}

func getPortOrDefault(url *url.URL, defaultPort string) string {
	port := url.Port()
	if port == "" {
		return defaultPort
	}
	return port
}

func asWSAddress(address string) (string, error) {
	url, err := url.Parse(address)
	if err != nil {
		return "", err
	}

	if url.Scheme == "https" {
		return fmt.Sprintf("wss://%s:%s", url.Hostname(), getPortOrDefault(url, "443")), nil
	}

	return fmt.Sprintf("ws://%s:%s", url.Hostname(), getPortOrDefault(url, "80")), nil
}

// StartListening sends a request message for listening to a given event type and awaits confirmation response
func StartListening(cfg *TestConfiguration, conn *websocket.Conn, eventType string, filter string) error {
	var msg string
	if len(filter) > 0 {
		msg = fmt.Sprintf("%s?filter=%s", eventType, filter)
	} else {
		msg = eventType
	}
	err := websocket.Message.Send(conn, msg)
	if err != nil {
		return err
	}
	return WaitForWSMessage(cfg, conn, fmt.Sprintf("%s:ACK", eventType))
}

// WaitForWSMessage polls messages from the web socket connection until specific message is received or timeout expires
func WaitForWSMessage(cfg *TestConfiguration, ws *websocket.Conn, expectedMessage string) error {
	deadline := time.Now().Add(MillisToDuration(cfg.WsEventTimeoutMs))
	ws.SetDeadline(deadline)

	var payload []byte
	for time.Now().Before(deadline) {
		err := websocket.Message.Receive(ws, &payload)
		if err != nil {
			return fmt.Errorf("error reading from websocket: %v", err)
		}
		message := strings.TrimSpace(string(payload))
		if message == expectedMessage {
			return nil
		}
	}

	return errors.New("timeout waiting for web socket message")
}

// ProcessWSMessages polls messages from the web socket connection until specific condition is satisfied or timeout expires
func ProcessWSMessages(cfg *TestConfiguration, ws *websocket.Conn, process func(*protocol.Envelope) (bool, error)) error {
	timeout := MillisToDuration(cfg.WsEventTimeoutMs)
	deadline := time.Now().Add(timeout)
	ws.SetDeadline(deadline)

	var err error
	finished := false

	for !finished && time.Now().Before(deadline) {
		var payload []byte
		webSocketErr := websocket.Message.Receive(ws, &payload)
		if webSocketErr != nil {
			return fmt.Errorf("error reading from websocket: %v", webSocketErr)
		}

		envelope := &protocol.Envelope{}
		unmarshalErr := json.Unmarshal(payload, envelope)
		if unmarshalErr == nil {
			finished, err = process(envelope)
		} else {
			// Unmarshalling error, the payload is not a JSON of protocol.Envelope
			// Ignore the error
		}
	}

	if !finished {
		return fmt.Errorf("not finished, expected WS response not received in %v, last error: %v", timeout, err)
	}

	return err
}

// ExecuteOperation executes an operation of a feature using http POST method
func ExecuteOperation(cfg *TestConfiguration, featureURL string, operation string, params interface{}) ([]byte, error) {
	url := fmt.Sprintf(featureOperationURLTemplate, featureURL, operation)
	return SendHTTPFeatureRequest(cfg, http.MethodPost, url, params)
}

// GetThingURL returns the url for executing operations on a thing
func GetThingURL(digitalTwinAPIAddress string, thingID string) string {
	return fmt.Sprintf(thingURLTemplate, strings.TrimSuffix(digitalTwinAPIAddress, "/"), thingID)
}

// GetFeatureURL returns the url of a feature
func GetFeatureURL(thingURL string, featureID string) string {
	return fmt.Sprintf(featureURLTemplate, thingURL, featureID)
}

// GetPropertyPath returns the path to a property on a feature
func GetPropertyPath(featureID string, name string) string {
	return fmt.Sprintf(featurePropertyPathTemplate, featureID, name)
}

// GetMessagePath returns the path to a message on a feature
func GetMessagePath(featureID string, name string) string {
	return fmt.Sprintf(messagePathTemplate, featureID, name)
}

// GetEventTopic returns the topic of an event on a thing
func GetEventTopic(thingID string, action protocol.TopicAction) string {
	thingIDWithNamespace := model.NewNamespacedIDFrom(thingID)
	t := protocol.Topic{
		Namespace:  thingIDWithNamespace.Namespace,
		EntityName: thingIDWithNamespace.Name,
		Group:      protocol.GroupThings,
		Channel:    protocol.ChannelTwin,
		Criterion:  protocol.CriterionEvents,
		Action:     action,
	}
	return t.String()
}

// GetLiveMessageTopic returns the topic of a live message on a thing
func GetLiveMessageTopic(thingID string, action protocol.TopicAction) string {
	thingIDWithNamespace := model.NewNamespacedIDFrom(thingID)
	t := protocol.Topic{
		Namespace:  thingIDWithNamespace.Namespace,
		EntityName: thingIDWithNamespace.Name,
		Group:      protocol.GroupThings,
		Channel:    protocol.ChannelLive,
		Criterion:  protocol.CriterionMessages,
		Action:     action,
	}
	return t.String()
}
