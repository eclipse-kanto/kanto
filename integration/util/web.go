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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eclipse/ditto-clients-golang/protocol"
	"golang.org/x/net/websocket"
)

// SendDigitalTwinRequest sends new HTTP request to ditto REST API
func SendDigitalTwinRequest(cfg *TestConfiguration, method string, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
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

// NewWSConnection creates new web socket connection
func NewWSConnection(cfg *TestConfiguration) (*websocket.Conn, error) {
	wsAddress, err := asWSAddress(cfg.DigitalTwinAPIAddress)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/ws/2", wsAddress)
	wscfg, err := websocket.NewConfig(url, cfg.DigitalTwinAPIAddress)
	if err != nil {
		return nil, err
	}

	auth := fmt.Sprintf("%s:%s", cfg.DigitalTwinAPIUsername, cfg.DigitalTwinAPIPassword)
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	wscfg.Header = http.Header{
		"Authorization": {"Basic " + enc},
	}

	return websocket.DialConfig(wscfg)
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

// WaitForWSMessage polls messages from the web socket connection until specific message is received or timeout expires
func WaitForWSMessage(timeout time.Duration, ws *websocket.Conn, expectedMessage string) error {
	var payload []byte
	deadline := time.Now().Add(timeout)
	ws.SetDeadline(deadline)

	for time.Now().Before(deadline) {
		err := websocket.Message.Receive(ws, &payload)
		if err != nil {
			return fmt.Errorf("error reading from websocket: %s", err)
		}
		message := strings.TrimSpace(string(payload))
		if message == expectedMessage {
			return nil
		}
	}

	return errors.New("timeout")
}

// ProcessWSMessageResult is the result of messages processing
type ProcessWSMessageResult struct {
	Finished bool
	Err      error
}

// ProcessWSMessages polls messages from the web socket connection until specific condition is satisfied or timeout expires
func ProcessWSMessages(
	timeout time.Duration,
	ws *websocket.Conn,
	process func(*protocol.Envelope) ProcessWSMessageResult) ProcessWSMessageResult {

	deadline := time.Now().Add(timeout)
	ws.SetDeadline(deadline)

	result := ProcessWSMessageResult{}
	for !result.Finished && time.Now().Before(deadline) {
		var payload []byte
		webSocketErr := websocket.Message.Receive(ws, &payload)
		if webSocketErr != nil {
			result.Err = fmt.Errorf("error reading from websocket: %s", webSocketErr)
			return result
		}
		envelope := &protocol.Envelope{}
		unmarshalErr := json.Unmarshal(payload, envelope)
		if unmarshalErr == nil {
			result = process(envelope)
		} else {
			// Unmarshalling error, the payload is not a JSON of protocol.Envelope
			// Ignore the error
		}
	}

	if !result.Finished {
		result.Err = fmt.Errorf("not finished, expected WS response not received in %v, last error: %v", timeout, result.Err)
	}

	return result
}
