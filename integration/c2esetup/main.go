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
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

const (
	devices     = "devices/"
	credentials = "credentials/"

	policies = "policies/"
	things   = "things/"

	indent = " "

	systemctl             = "systemctl"
	restart               = "restart"
	suiteConnectorService = "suite-connector.service"

	fileFlags = 0666
)

var (
	c2eCfg c2eConfig

	tenantID string
	policyID string

	deviceID   string
	devicePass string

	authID string

	configFile       string
	configFileBackup string
	certFile         string
	logFile          string

	envVariablesPrefix string
)

type resource struct {
	base string
	path string

	method string
	body   string

	user string
	pass string

	delete bool
}

type c2eConfig struct {
	Broker                   string `def:"tcp://localhost:1883"`
	MqttQuiesceMs            int    `def:"500"`
	MqttAcknowledgeTimeoutMs int    `def:"3000"`

	RegistryAddress string

	RegistryUser     string `def:""`
	RegistryPassword string `def:""`

	DittoAddress string

	DittoUser     string `def:"ditto"`
	DittoPassword string `def:"ditto"`

	MqttAdapterAddress string
}

type thingConfig struct {
	DeviceID string `json:"deviceId"`
	TenantID string `json:"tenantId"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&tenantID, "tenant", "", "Hono tenant ID")

	flag.StringVar(&deviceID, "deviceID", "", "Test device ID, defaults to randomly generated. You should not use this param in a LAVA setting")
	flag.StringVar(&devicePass, "devicePass", "123456", "Test device password")

	flag.StringVar(&policyID, "policyID", "", "Test device's policy ID")

	flag.StringVar(&configFile, "configFile", "/etc/suite-connector/config.json", "Path to Suite Connector config file")
	flag.StringVar(&configFileBackup, "configFileBackup", "/etc/suite-connector/configBackup.json", "Path to Suite Connector config file backup")
	flag.StringVar(&certFile, "certFile", "/etc/suite-connector/iothub.crt", "Path to Suite Connector CA certificates file")
	flag.StringVar(&logFile, "logFile", "/var/log/suite-connector/suite-connector.log", "Path for Suite Connector log file")

	flag.StringVar(&envVariablesPrefix, "envprefix", "SCT", "Test environmental variables prefix")

	cleanup := flag.Bool("cleanup", false, "Clean up test resources")

	flag.Parse()

	assertFlag(tenantID, "tenantID")
	assertFlag(policyID, "policyID")

	err := initConfigFromEnv(&c2eCfg, envVariablesPrefix)
	if err != nil {
		fmt.Println(err)
		fmt.Println(getConfigHelp(c2eCfg, envVariablesPrefix))
		os.Exit(1)
	}

	configureDeviceID(*cleanup)

	if deviceID == "" {
		fmt.Println("can't find device id")

		os.Exit(1)
	}

	registryAPI := fmt.Sprintf("%s/v1", strings.TrimSuffix(c2eCfg.RegistryAddress, "/"))
	devicePath := fmt.Sprintf("%s/%s", tenantID, deviceID)

	resources := make([]*resource, 0, 4)
	resources = append(resources, &resource{base: registryAPI, path: devices + devicePath, method: http.MethodPost,
		body: deviceJSON, user: c2eCfg.RegistryUser, pass: c2eCfg.RegistryPassword, delete: true})

	authID = strings.ReplaceAll(deviceID, ":", "_")
	auth := fmt.Sprintf(authJSON, authID, devicePass)
	resources = append(resources, &resource{base: registryAPI, path: credentials + devicePath, method: http.MethodPut,
		body: auth, user: c2eCfg.RegistryUser, pass: c2eCfg.RegistryPassword, delete: false})

	dittoAPI := fmt.Sprintf("%s/api/2", strings.TrimSuffix(c2eCfg.DittoAddress, "/"))
	thing := fmt.Sprintf(thingJSON, policyID)
	resources = append(resources, &resource{base: dittoAPI, path: things + deviceID, method: http.MethodPut,
		body: thing, user: c2eCfg.DittoUser, pass: c2eCfg.DittoPassword, delete: true})

	code := 0
	if *cleanup {
		fmt.Printf("performing cleanup on device id: %s\n", deviceID)
		performCleanUp(resources)
		if err := copyFile(configFileBackup, configFile); err != nil {
			fmt.Printf(
				"unable to restore the backup copy of config file %s to %s: %v\n",
				configFileBackup, configFile, err)
			code = 1
		} else {
			code = restartSuiteConnector()
		}
	} else {
		code = performSetUp(resources)
		if code == 0 {
			code = restartSuiteConnector()
		}
		fmt.Println("setup complete")
	}
	os.Exit(code)
}

func configureDeviceID(cleanup bool) {
	if !cleanup {
		if deviceID == "" {
			deviceID = generateRandomDeviceID()
			fmt.Printf("generating a random device id, use mqtt to read it and use it later: \"%s\"\n", deviceID)
		} else {
			fmt.Printf("forcing deviceID \"%s\"\n", deviceID)
		}
	} else {
		opts := mqtt.NewClientOptions().
			AddBroker(c2eCfg.Broker).
			SetClientID(uuid.New().String()).
			SetKeepAlive(30 * time.Second).
			SetCleanSession(true).
			SetAutoReconnect(true)

		mqttClient := mqtt.NewClient(opts)

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			fmt.Printf("can't connect to MQTT broker: %v\n", token.Error())
		} else {
			defer mqttClient.Disconnect(uint(c2eCfg.MqttQuiesceMs))
			thingCfg, err := getThingConfig(mqttClient)
			if err != nil {
				fmt.Printf("can't get thing config from MQTT broker: %v\n", err)
			} else {
				fmt.Printf("found thing id from MQTT broker: %s", thingCfg.DeviceID)
				if deviceID != "" {
					fmt.Printf(", overriding %s", deviceID)
				}
				fmt.Println()
				deviceID = thingCfg.DeviceID
			}
		}
	}
}

func generateRandomDeviceID() string {
	return fmt.Sprintf("test:dev%d", rand.Intn(100_000))
}

func assertFlag(value string, name string) {
	if value == "" {
		fmt.Printf("'%s' must not be empty, but is not specified\n", name)
		flag.Usage()
		os.Exit(1)
	}
}

func performSetUp(resources []*resource) int {
	fmt.Println("performing setup...")
	for i, r := range resources {
		url := fmt.Sprintf("%s/%s", r.base, r.path)

		if b, err := doRequest(r.method, url, ([]byte)(r.body), r.user, r.pass); err != nil {
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

	if err := copyFile(configFile, configFileBackup); err != nil {
		fmt.Printf(
			"unable to save backup copy of config file %s to %s: %v\n",
			configFile, configFileBackup, err)

		performCleanUp(resources)
		return 1
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
	if deviceID == "" {
		return
	}

	fmt.Println("performing cleanup of initially created things...")
	for i := len(resources) - 1; i >= 0; i-- {
		r := resources[i]

		if !r.delete {
			continue
		}

		url := fmt.Sprintf("%s/%s", r.base, r.path)

		if _, err := doRequest(http.MethodDelete, url, nil, r.user, r.pass); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s '%s' deleted\n", indent, r.path)
		}
	}
	// Dependent devices use the same deviceID as prefix, followed by colon ':'
	deleteDependentDevices(deviceID + ":")
	fmt.Println("cleanup done")
}

func deleteDependentDevices(deviceIDPrefix string) {
	deleteHonoDevices(deviceIDPrefix)
	deleteDittoDevices(deviceIDPrefix)
}

func deleteHonoDevices(deviceIDPrefix string) {
	type honoThing struct {
		Id string `json:"id"`
	}

	type honoThings struct {
		Things []*honoThing `json:"result"`
	}

	// Delete Hono devices
	fmt.Println("performing additional cleanup in hono...")
	honoAPI := fmt.Sprintf(
		"%s/v1/%s%s/", strings.TrimSuffix(c2eCfg.RegistryAddress, "/"), devices, tenantID)
	thingsJson, err := doRequest("GET", honoAPI, nil, c2eCfg.RegistryUser, c2eCfg.RegistryPassword)
	if err != nil {
		fmt.Println(err)
	} else {
		honoDevices := &honoThings{}
		if err = json.Unmarshal(thingsJson, honoDevices); err != nil {
			fmt.Println(err)
		}
		for _, thing := range honoDevices.Things {
			if strings.HasPrefix(thing.Id, deviceIDPrefix) {
				if _, err = doRequest(http.MethodDelete, honoAPI+thing.Id, nil, c2eCfg.RegistryUser, c2eCfg.RegistryPassword); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%s '%s' deleted\n", indent, thing.Id)
				}
			}
		}
	}
}

func deleteDittoDevices(deviceIDPrefix string) {
	type dittoThing struct {
		ThingId string `json:"thingId"`
	}

	// Delete Ditto devices
	fmt.Println("performing additional cleanup in ditto...")
	dittoAPI := fmt.Sprintf("%s/api/2/%s", strings.TrimSuffix(c2eCfg.DittoAddress, "/"), things)
	thingsJson, err := doRequest("GET", dittoAPI, nil, c2eCfg.DittoUser, c2eCfg.DittoPassword)
	if err != nil {
		fmt.Println(err)
	} else {
		dittoDevices := []dittoThing{}
		if err = json.Unmarshal(thingsJson, &dittoDevices); err != nil {
			fmt.Println(err)
		}
		for _, thing := range dittoDevices {
			if strings.HasPrefix(thing.ThingId, deviceIDPrefix) {
				if _, err = doRequest(http.MethodDelete, dittoAPI+thing.ThingId, nil, c2eCfg.DittoUser, c2eCfg.DittoPassword); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%s '%s' deleted\n", indent, thing.ThingId)
				}
			}
		}
	}
}

func writeConfig(path string) error {
	type connectorConfig struct {
		CaCert   string `json:"caCert"`
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
	cfg.Address = c2eCfg.MqttAdapterAddress
	cfg.DeviceID = deviceID
	cfg.TenantID = tenantID
	cfg.AuthID = authID
	cfg.Password = devicePass

	b, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, fileFlags)
}

func doRequest(method string, url string, payload []byte, user string, pass string) ([]byte, error) {
	var r io.Reader

	if payload != nil {
		r = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	if user != "" {
		req.SetBasicAuth(user, pass)
	}

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

func getThingConfig(mqttClient mqtt.Client) (*thingConfig, error) {
	type result struct {
		cfg *thingConfig
		err error
	}

	ch := make(chan result)

	if token := mqttClient.Subscribe("edge/thing/response", 1, func(client mqtt.Client, message mqtt.Message) {
		var cfg thingConfig
		if err := json.Unmarshal(message.Payload(), &cfg); err != nil {
			ch <- result{nil, err}
		}
		ch <- result{&cfg, nil}
	}); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

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

func restartSuiteConnector() int {
	fmt.Println("restarting suite-connector...")
	code := 0
	cmd := exec.Command(systemctl, restart, suiteConnectorService)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("error restarting %s: %v", suiteConnectorService, err)
		code = -1
	}
	fmt.Println(string(stdout))
	if code == 0 {
		fmt.Println("... done")
	}
	return code
}

func copyFile(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, fileFlags)
}
