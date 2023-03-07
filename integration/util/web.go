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

// SubscribeEventType is an event type description to be used with SubscribeForWSMessages.
// It specifies the type of messages which should be listened for.
type SubscribeEventType string

// UnsubscribeEventType is an event type description to be used with UnsubscribeFromWSMessages.
// It specifies the type of messages which should no longer be listened for.
type UnsubscribeEventType string

const (
	thingURLTemplate                 = "%s/api/2/things/%s"
	featureURLTemplate               = "%s/features/%s"
	featurePropertyURLTemplate       = "%s/properties/%s"
	featureOperationURLTemplate      = "%s/inbox/messages/%s"
	featurePropertyPathTemplate      = "/features/%s/properties/%s"
	featureMessageOutboxPathTemplate = "/features/%s/outbox/messages/%s"
	featureMessageInboxPathTemplate  = "/features/%s/inbox/messages/%s"

	// StartSendEvents specifies that events should be received.
	StartSendEvents SubscribeEventType = "START-SEND-EVENTS"

	// StopSendEvents specifies that events should no longer be received.
	StopSendEvents UnsubscribeEventType = "STOP-SEND-EVENTS"

	// StartSendMessages specifies that messages should be received.
	StartSendMessages SubscribeEventType = "START-SEND-MESSAGES"

	// StopSendMessages specifies that messages should no longer be received.
	StopSendMessages UnsubscribeEventType = "STOP-SEND-MESSAGES"
)

// SendDigitalTwinRequest sends a new HTTP request to the Ditto REST API
func SendDigitalTwinRequest(cfg *TestConfiguration, method string, url string, body interface{}) ([]byte, error) {
	var (
		payload []byte
		err     error
	)
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := createRequest(payload, true, method, url, cfg.DigitalTwinAPIUsername, cfg.DigitalTwinAPIPassword)
	if err != nil {
		return nil, err
	}
	return sendRequest(req, method, url)
}

// SendDeviceRegistryRequest sends a new HTTP request to the Ditto API
func SendDeviceRegistryRequest(payload []byte, method string, url string, username string, password string) ([]byte, error) {
	req, err := createRequest(payload, false, method, url, username, password)
	if err != nil {
		return nil, err
	}
	return sendRequest(req, method, url)
}

func createRequest(payload []byte, rspRequired bool, method, url, username, password string) (*http.Request, error) {
	var reqBody io.Reader

	if payload != nil {
		reqBody = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if payload != nil {
		req.Header.Add("Content-Type", "application/json")
		if rspRequired {
			correlationID := uuid.New().String()
			req.Header.Add("correlation-id", correlationID)
			req.Header.Add("response-required", "true")
		}
	}

	req.SetBasicAuth(username, password)
	return req, nil
}

func sendRequest(req *http.Request, method string, url string) ([]byte, error) {
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

// NewDigitalTwinWSConnection creates a new WebSocket connection
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

func sendMessageAndAwaitAck(cfg *TestConfiguration, conn *websocket.Conn, msg string, eventType string) error {
	err := websocket.Message.Send(conn, msg)
	if err != nil {
		return err
	}
	return WaitForWSMessage(cfg, conn, eventType+":ACK")
}

// SubscribeForWSMessages subscribes for the messages that are sent from a WebSocket session and awaits confirmation response.
func SubscribeForWSMessages(cfg *TestConfiguration, conn *websocket.Conn, eventType SubscribeEventType, filter string) error {
	var msg string
	if len(filter) > 0 {
		msg = fmt.Sprintf("%s?filter=%s", eventType, filter)
	} else {
		msg = string(eventType)
	}
	return sendMessageAndAwaitAck(cfg, conn, msg, string(eventType))
}

// UnsubscribeFromWSMessages unsubscribes from the messages that are sent from a WebSocket session
// and awaits confirmation response.
func UnsubscribeFromWSMessages(cfg *TestConfiguration, ws *websocket.Conn, eventType UnsubscribeEventType) error {
	return sendMessageAndAwaitAck(cfg, ws, string(eventType), string(eventType))
}

// WaitForWSMessage waits for received a specific message from a WebSocket session or timeout expires
func WaitForWSMessage(cfg *TestConfiguration, ws *websocket.Conn, expectedMessage string) error {
	deadline := time.Now().Add(MillisToDuration(cfg.WSEventTimeoutMS))
	if err := ws.SetDeadline(deadline); err != nil {
		return fmt.Errorf("unable to set deadline to websocket: %v", err)
	}
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

	return errors.New("timeout waiting for websocket message")
}

// ProcessWSMessages processes messages for the satisfied condition from the WebSocket session or timeout expires
func ProcessWSMessages(cfg *TestConfiguration, ws *websocket.Conn, process func(*protocol.Envelope) (bool, error)) error {
	timeout := MillisToDuration(cfg.WSEventTimeoutMS)
	deadline := time.Now().Add(timeout)
	if err := ws.SetDeadline(deadline); err != nil {
		return fmt.Errorf("unable to set deadline to websocket: %v", err)
	}

	var err error
	finished := false

	for !finished && time.Now().Before(deadline) {
		var payload []byte
		wsErr := websocket.Message.Receive(ws, &payload)
		if wsErr != nil {
			return fmt.Errorf("error reading from websocket: %v", wsErr)
		}

		envelope := &protocol.Envelope{}
		if unmarshalErr := json.Unmarshal(payload, envelope); unmarshalErr == nil {
			finished, err = process(envelope)
		} else {
			// Unmarshalling error, the payload is not a JSON of protocol.Envelope
			// Ignore the error
		}
	}

	if !finished {
		return fmt.Errorf("not finished, expected websocket response not received in %v, last error: %v", timeout, err)
	}

	return err
}

// ExecuteOperation executes an operation of a feature
func ExecuteOperation(cfg *TestConfiguration, featureURL string, operation string, params interface{}) ([]byte, error) {
	url := fmt.Sprintf(featureOperationURLTemplate, featureURL, operation)
	return SendDigitalTwinRequest(cfg, http.MethodPost, url, params)
}

// GetFeaturePropertyValue gets the value of a feature's property
func GetFeaturePropertyValue(cfg *TestConfiguration, featureURL string, property string) ([]byte, error) {
	url := fmt.Sprintf(featurePropertyURLTemplate, featureURL, property)
	return SendDigitalTwinRequest(cfg, http.MethodGet, url, nil)
}

// GetThingURL returns the url of a thing
func GetThingURL(digitalTwinAPIAddress string, thingID string) string {
	return fmt.Sprintf(thingURLTemplate, strings.TrimSuffix(digitalTwinAPIAddress, "/"), thingID)
}

// GetFeatureURL returns the url of a feature
func GetFeatureURL(thingURL string, featureID string) string {
	return fmt.Sprintf(featureURLTemplate, thingURL, featureID)
}

// GetFeaturePropertyPath returns the path to a property on a feature
func GetFeaturePropertyPath(featureID string, name string) string {
	return fmt.Sprintf(featurePropertyPathTemplate, featureID, name)
}

// GetFeatureOutboxMessagePath returns the path to an outbox message of a feature
func GetFeatureOutboxMessagePath(featureID string, name string) string {
	return fmt.Sprintf(featureMessageOutboxPathTemplate, featureID, name)
}

// GetFeatureInboxMessagePath returns the path to an inbox message of a feature
func GetFeatureInboxMessagePath(featureID string, name string) string {
	return fmt.Sprintf(featureMessageInboxPathTemplate, featureID, name)
}

// GetTwinEventTopic returns the twin event topic
func GetTwinEventTopic(fullThingID string, action protocol.TopicAction) string {
	thingID := model.NewNamespacedIDFrom(fullThingID)
	t := (&protocol.Topic{}).
		WithNamespace(thingID.Namespace).
		WithEntityName(thingID.Name).
		WithGroup(protocol.GroupThings).
		WithChannel(protocol.ChannelTwin).
		WithCriterion(protocol.CriterionEvents).
		WithAction(action)
	return t.String()
}

// GetLiveMessageTopic returns the live message topic
func GetLiveMessageTopic(fullThingID string, action protocol.TopicAction) string {
	thingID := model.NewNamespacedIDFrom(fullThingID)
	t := (&protocol.Topic{}).
		WithNamespace(thingID.Namespace).
		WithEntityName(thingID.Name).
		WithGroup(protocol.GroupThings).
		WithChannel(protocol.ChannelLive).
		WithCriterion(protocol.CriterionMessages).
		WithAction(action)
	return t.String()
}
