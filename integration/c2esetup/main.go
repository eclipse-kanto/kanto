// Copyright (c) 2021 Contributors to the Eclipse Foundation
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

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type c2eConfig struct {
	RegistryHost string
	RegistryPort int

	DittoHost string
	DittoPort int

	DittoUser     string `def:"ditto"`
	DittoPassword string `def:"ditto"`

	MqttAdapterHost string
	MqttAdapterPort int
}

const (
	devices     = "devices/"
	credentials = "credentials/"

	policies = "policies/"
	things   = "things/"

	indent = " "
)

var (
	c2eCfg c2eConfig

	connectionID string
	tenantID     string

	deviceID   string
	devicePass string

	authID string

	configFile string
	certFile   string
	logFile    string

	envVariablesPrefix string
)

type resource struct {
	base string
	path string

	method string
	body   string

	delete bool
}

func main() {
	flag.StringVar(&tenantID, "tenant", "", "Hono tenant ID")
	flag.StringVar(&connectionID, "conn", "", "Hono connection ID")

	flag.StringVar(&deviceID, "deviceID", "test:device", "Test device ID")
	flag.StringVar(&devicePass, "devicePass", "123456", "Test device password")

	flag.StringVar(&configFile, "configFile", "/etc/suite-connector/config.json", "Path to Suite Connector config file")
	flag.StringVar(&certFile, "certFile", "/etc/suite-connector/iothub.crt", "Path to Suite Connector CA certificates file")
	flag.StringVar(&logFile, "logFile", "/var/log/suite-connector/suite-connector.log", "Path for Suite Connector log file")

	flag.StringVar(&envVariablesPrefix, "envprefix", "SCT", "Test environmental variables prefix")

	cleanup := flag.Bool("cleanup", false, "Clean up test resources")

	flag.Parse()

	assertFlag(tenantID, "Hono tenant ID")
	assertFlag(connectionID, "Hono connection ID")

	if tenantID == "" {
		fmt.Println("tenant can not be empty")

		os.Exit(1)
	}

	err := initConfigFromEnv(&c2eCfg, envVariablesPrefix)
	if err != nil {
		fmt.Println(err)
		fmt.Println(getConfigHelp(c2eCfg, envVariablesPrefix))
		os.Exit(1)
	}

	registryAPI := fmt.Sprintf("http://%s:%d/v1", c2eCfg.RegistryHost, c2eCfg.RegistryPort)
	devicePath := fmt.Sprintf("%s/%s", tenantID, deviceID)

	resources := make([]*resource, 0, 4)
	resources = append(resources, &resource{base: registryAPI, path: devices + devicePath, method: http.MethodPost, body: deviceJSON, delete: true})

	authID = strings.ReplaceAll(deviceID, ":", "_")
	auth := fmt.Sprintf(authJSON, authID, devicePass)
	resources = append(resources, &resource{base: registryAPI, path: credentials + devicePath, method: http.MethodPut, body: auth, delete: false})

	dittoAPI := fmt.Sprintf("http://%s:%d/api/2", c2eCfg.DittoHost, c2eCfg.DittoPort)
	policy := fmt.Sprintf(policyJSON, connectionID)
	resources = append(resources, &resource{base: dittoAPI, path: policies + deviceID, method: http.MethodPut, body: policy, delete: true})

	thing := fmt.Sprintf(thingJSON, deviceID)
	resources = append(resources, &resource{base: dittoAPI, path: things + deviceID, method: http.MethodPut, body: thing, delete: true})

	if *cleanup {
		performCleanUp(resources)
	} else {
		code := performSetUp(resources)

		os.Exit(code)
	}
}

func assertFlag(value string, name string) {
	if value == "" {
		fmt.Printf("'%s' not specified\n", name)
		flag.Usage()
		os.Exit(1)
	}
}

func performSetUp(resources []*resource) int {
	fmt.Println("performing setup...")
	for i, r := range resources {
		url := fmt.Sprintf("%s/%s", r.base, r.path)

		if b, err := doRequest(r.method, url, ([]byte)(r.body)); err != nil {
			fmt.Println(err)
			if b != nil {
				fmt.Println(b)
				fmt.Println()
			}

			if i > 0 {
				performCleanUp(resources[:i])
			}

			return 1
		}

		fmt.Printf("%s '%s' created\n", indent, r.path)
	}

	if err := writeConfig(configFile); err != nil {
		fmt.Println(err)

		performCleanUp(resources)
		return 1
	}

	fmt.Printf("%s config file '%s' written\n", indent, configFile)

	fmt.Println("setup successful")

	return 0
}

func performCleanUp(resources []*resource) {
	fmt.Println("performing cleanup...")
	for i := len(resources) - 1; i >= 0; i-- {
		r := resources[i]

		if !r.delete {
			continue
		}

		url := fmt.Sprintf("%s/%s", r.base, r.path)

		if _, err := doRequest(http.MethodDelete, url, nil); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s '%s' deleted\n", indent, r.path)
		}
	}

	fmt.Println("cleanup successful")
}

func writeConfig(path string) error {
	type connectorConfig struct {
		CaCert   string `json:"cacert"`
		LogFile  string `json:"logFile"`
		Address  string `json:"address"`
		TenantID string `json:"tenantId"`
		DeviceID string `json:"deviceId"`
		AuthID   string `json:"authId"`
		Password string `json:"password"`
	}

	cfg := &connectorConfig{}
	cfg.CaCert = certFile
	cfg.LogFile = logFile
	cfg.Address = fmt.Sprintf("tcp://%s:%d", c2eCfg.MqttAdapterHost, c2eCfg.MqttAdapterPort)
	cfg.DeviceID = deviceID
	cfg.TenantID = tenantID
	cfg.AuthID = authID
	cfg.Password = devicePass

	b, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 666)
}

func doRequest(method string, url string, payload []byte) ([]byte, error) {
	var r io.Reader

	if payload != nil {
		r = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c2eCfg.DittoUser, c2eCfg.DittoPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, fmt.Errorf("%s %s request failed: %s", method, url, resp.Status)
	}

	return body, err
}
