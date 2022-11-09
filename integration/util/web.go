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
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

// SendDittoRequest sends new HTTP request to ditto REST API
func SendDittoRequest(cfg *TestConfiguration, method string, url string) ([]byte, error) {
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

// WaitResult waits for the result or error to be received over the channel up to a timeout
func WaitResult(timeout time.Duration, resultCh chan error, closer func()) chan error {
	ch := make(chan error)

	go func() {
		select {
		case result := <-resultCh:
			ch <- result
		case <-time.After(timeout):
			closer()
			ch <- errors.New("timeout")
		}
	}()

	return ch
}

// BeginWSWait waits for a message to be received via websocket
func BeginWSWait(cfg *TestConfiguration, ws *websocket.Conn, check func(payload []byte) error) chan error {
	timeout := time.Duration(cfg.EventTimeoutMs * int(time.Millisecond))
	resultCh := make(chan error)

	go func() {
		var payload []byte
		threshold := time.Now().Add(timeout)
		var err error
		for time.Now().Before(threshold) {
			err = websocket.Message.Receive(ws, &payload)
			if err == nil {
				err = check(payload)
			}
			if err == nil {
				resultCh <- nil
				return
			}
		}
		resultCh <- fmt.Errorf("WS response not received in %v, last error: %v", timeout, err)
	}()

	closer := func() {
		ws.Close()
	}

	return WaitResult(timeout, resultCh, closer)
}
