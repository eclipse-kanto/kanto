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
	"errors"
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

	MQTT "github.com/eclipse/paho.mqtt.golang"
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

	configDefaultMode = 0666
)

var (
	c2eCfg c2eConfig

	tenantId string
	policyId string

	deviceId   string
	devicePass string

	authId string

	configFile       string
	configFileBackup string
	certFile         string
	logFile          string
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

	DeviceRegistryApiAddress string

	DeviceRegistryApiUser     string `def:""`
	DeviceRegistryApiPassword string `def:""`

	DigitalTwinApiAddress string

	DigitalTwinApiUser     string `def:"ditto"`
	DigitalTwinApiPassword string `def:"ditto"`

	MqttAdapterAddress string
}

type thingConfig struct {
	DeviceId string `json:"deviceId"`
	TenantId string `json:"tenantId"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&tenantId, "tenant", "", "Hono tenant id")

	flag.StringVar(&deviceId, "device", "", "Test device id, defaults to randomly generated")
	flag.StringVar(&devicePass, "devicePass", "123456", "Test device password")

	flag.StringVar(&policyId, "policy", "", "Test device's policy id")

	flag.StringVar(&configFile, "configFile", "/etc/suite-connector/config.json", "Path to Suite Connector configuration file")
	flag.StringVar(&configFileBackup, "configFileBackup", "/etc/suite-connector/configBackup.json", "Path to Suite Connector configuration file backup")
	flag.StringVar(&certFile, "certFile", "/etc/suite-connector/iothub.crt", "Path to Suite Connector CA certificates file")
	flag.StringVar(&logFile, "logFile", "/var/log/suite-connector/suite-connector.log", "Path to Suite Connector log file")

	clean := flag.Bool("clean", false, "Clean up test resources")

	flag.Parse()

	err := initConfigFromEnv(&c2eCfg)
	if err != nil {
		fmt.Println(err)
		fmt.Println(getConfigHelp(c2eCfg))
		os.Exit(1)
	}

	if !*clean {
		if deviceId == "" {
			deviceId = generateRandomDeviceID()
			fmt.Printf("generating a random device id: \"%s\"\n", deviceId)
		} else {
			fmt.Printf("forcing device id: \"%s\"\n", deviceId)
		}
	} else if deviceId == "" || tenantId == "" {
		readIdsFromMQTT()
	}

	if !*clean {
		assertFlag(tenantId, "tenant id")
		assertFlag(policyId, "policy id")
	}

	registryAPI := fmt.Sprintf("%s/v1", strings.TrimSuffix(c2eCfg.DeviceRegistryApiAddress, "/"))
	devicePath := fmt.Sprintf("%s/%s", tenantId, deviceId)

	resources := make([]*resource, 0, 4)
	deviceResource := &resource{base: registryAPI, path: devices + devicePath, method: http.MethodPost,
		body: deviceJSON, user: c2eCfg.DeviceRegistryApiUser, pass: c2eCfg.DeviceRegistryApiPassword, delete: true}
	resources = append(resources, deviceResource)

	authId = strings.ReplaceAll(deviceId, ":", "_")
	auth := fmt.Sprintf(authJSON, authId, devicePass)
	resources = append(resources, &resource{base: registryAPI, path: credentials + devicePath, method: http.MethodPut,
		body: auth, user: c2eCfg.DeviceRegistryApiUser, pass: c2eCfg.DeviceRegistryApiPassword, delete: false})

	dittoAPI := fmt.Sprintf("%s/api/2", strings.TrimSuffix(c2eCfg.DigitalTwinApiAddress, "/"))
	thing := fmt.Sprintf(thingJSON, policyId)
	resources = append(resources, &resource{base: dittoAPI, path: things + deviceId, method: http.MethodPut,
		body: thing, user: c2eCfg.DigitalTwinApiUser, pass: c2eCfg.DigitalTwinApiPassword, delete: true})

	var code int
	if !*clean {
		code = performSetUp(deviceResource, resources)
		fmt.Println("setup complete")
	} else {
		code = performCleanUp(resources)
		fmt.Println("cleanup complete")
	}
	os.Exit(code)
}

func readIdsFromMQTT() {
	deviceId = ""
	policyId = ""

	opts := MQTT.NewClientOptions().
		AddBroker(c2eCfg.Broker).
		SetClientID(uuid.New().String()).
		SetKeepAlive(30 * time.Second).
		SetCleanSession(true).
		SetAutoReconnect(true)

	mqttClient := MQTT.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("can't connect to MQTT broker: %v\n", token.Error())
	} else {
		defer mqttClient.Disconnect(uint(c2eCfg.MqttQuiesceMs))
		thingCfg, err := getThingConfig(mqttClient)
		if err != nil {
			fmt.Printf("can't get thing configuration from MQTT broker: %v\n", err)
		} else {
			deviceId = thingCfg.DeviceId
			tenantId = thingCfg.TenantId
		}
	}

	if deviceId == "" {
		fmt.Println("can't find device id")
		os.Exit(1)
	}
	fmt.Printf("found thing id from MQTT broker: %s\n", deviceId)
	if tenantId == "" {
		fmt.Println("can't find tenant id")
		os.Exit(1)
	}
	fmt.Printf("found tenant id from MQTT broker: %s\n", tenantId)
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

func checkDeviceInRegistry(deviceId string, deviceResource *resource) error {
	url := fmt.Sprintf("%s/%s", deviceResource.base, deviceResource.path)
	devicesBytes, err := doRequest("GET", url, nil, c2eCfg.DeviceRegistryApiUser, c2eCfg.DeviceRegistryApiPassword)
	if err == nil {
		devicesJson := string(devicesBytes)
		if strings.Contains(devicesJson, "status") &&
			strings.Contains(devicesJson, "created") &&
			strings.Contains(devicesJson, "authorities") {
			return nil
		}
	}
	if err == nil {
		return fmt.Errorf("device %s hasn't been created", deviceId)
	}
	return err
}

func performSetUp(deviceResource *resource, resources []*resource) int {
	if err := checkDeviceInRegistry(deviceId, deviceResource); err == nil {
		fmt.Printf("device %s already exists in registry, aborting...\n", deviceId)
		os.Exit(1)
	}

	fmt.Println("performing setup...")

	for i, r := range resources {
		url := fmt.Sprintf("%s/%s", r.base, r.path)

		if b, err := doRequest(r.method, url, ([]byte)(r.body), r.user, r.pass); err != nil {
			fmt.Println(err)
			if b != nil {
				fmt.Println(string(b))
				fmt.Println()
			}

			if i > 0 {
				deleteResources(resources[:i])
			}

			return 1
		}
		fmt.Printf("%s '%s' created\n", indent, r.path)
	}

	fmt.Println("checking if the device was successfully created in the registry")
	if err := checkDeviceInRegistry(deviceId, deviceResource); err != nil {
		fmt.Printf("%v\n", err)
		deleteResources(resources)
		return 1
	}

	if err := copyFile(configFile, configFileBackup); err != nil {
		fmt.Printf(
			"unable to save backup copy of configuration file %s to %s: %v\n",
			configFile, configFileBackup, err)

		deleteResources(resources)
		return 1
	}

	if err := writeConfigFile(configFile); err != nil {
		fmt.Println(err)

		deleteResources(resources)
		return 1
	}

	fmt.Printf("%s configuration file '%s' written\n", indent, configFile)

	code := restartSuiteConnector()

	if code == 0 {
		fmt.Println("setup successful")
	}

	return code
}

func performCleanUp(resources []*resource) int {
	var code int
	fmt.Println("restoring suite-connector configuration file and restarting suite-connector")
	if err := copyFile(configFileBackup, configFile); err != nil {
		fmt.Printf(
			"unable to restore the backup copy of configuration file %s to %s: %v\n",
			configFileBackup, configFile, err)
		code = 1
	} else {
		code = restartSuiteConnector()
	}
	if code == 0 {
		// Delete suite-connector configuration backup file
		if _, err := os.Stat(configFileBackup); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("unable to delete configuration file backup %s, file does not exist: %v", configFileBackup, err)
		} else if err := os.Remove(configFileBackup); err != nil {
			fmt.Printf("unable to delete configuration file backup %s, error: %v", configFileBackup, err)
		}
	}
	// Delete devices and things
	fmt.Printf("performing cleanup on device id: %s\n", deviceId)
	deleteResources(resources)
	return code
}

func deleteResources(resources []*resource) {
	deleteRelatedDevices(deviceId)
	fmt.Println("deleting initially created things...")
	// Delete in reverse order of creation
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
}

func deleteRelatedDevices(viaDeviceId string) {
	devicesVia := findHonoDevicesVia(viaDeviceId)
	// Ditto things are created after Hono devices, so delete them first
	deleteDittoThings(devicesVia)
	// Then delete Hono devices
	deleteHonoDevices(devicesVia)
}

func findHonoDevicesVia(viaDeviceId string) []string {
	var relations []string

	type registryDevice struct {
		Id  string   `json:"id"`
		Via []string `json:"via"`
	}

	type registryDevices struct {
		Devices []*registryDevice `json:"result"`
	}

	contains := func(where []string, what string) bool {
		for _, item := range where {
			if item == what {
				return true
			}
		}
		return false
	}

	honoAPI := fmt.Sprintf(
		"%s/v1/%s%s/", strings.TrimSuffix(c2eCfg.DeviceRegistryApiAddress, "/"), devices, tenantId)
	devicesJson, err := doRequest("GET", honoAPI, nil, c2eCfg.DeviceRegistryApiUser, c2eCfg.DeviceRegistryApiPassword)
	if err != nil {
		fmt.Println(err)
	} else {
		devices := &registryDevices{}
		err = json.Unmarshal(devicesJson, devices)
		if err != nil {
			fmt.Println(err)
			devices.Devices = nil
		}
		for _, device := range devices.Devices {
			if contains(device.Via, viaDeviceId) {
				relations = append(relations, device.Id)
			}
		}
	}

	return relations
}

func deleteDittoThings(relations []string) {
	// Delete related Ditto things
	fmt.Println("deleting automatically created things in ditto...")
	dittoAPI := fmt.Sprintf("%s/api/2/", strings.TrimSuffix(c2eCfg.DigitalTwinApiAddress, "/"))
	for _, thingId := range relations {
		_, err := doRequest(http.MethodDelete, dittoAPI+thingId, nil, c2eCfg.DigitalTwinApiUser, c2eCfg.DigitalTwinApiPassword)
		if err != nil {
			fmt.Printf("error deleting thing: %v\n", err)
		} else {
			fmt.Printf("%s '%s' deleted\n", indent, thingId)
		}
	}
}

func deleteHonoDevices(relations []string) {
	// Delete related Hono devices
	fmt.Println("deleting automatically created devices in hono...")
	honoAPI := fmt.Sprintf(
		"%s/v1/%s%s/", strings.TrimSuffix(c2eCfg.DeviceRegistryApiAddress, "/"), devices, tenantId)
	for _, device := range relations {
		if _, err := doRequest(http.MethodDelete, honoAPI+device, nil, c2eCfg.DeviceRegistryApiUser, c2eCfg.DeviceRegistryApiPassword); err != nil {
			fmt.Printf("error deleting device: %v\n", err)
		} else {
			fmt.Printf("%s '%s' deleted\n", indent, device)
		}
	}
}

func writeConfigFile(path string) error {
	type connectorConfig struct {
		CaCert   string `json:"caCert"`
		LogFile  string `json:"logFile"`
		Address  string `json:"address"`
		TenantID string `json:"tenantId"`
		DeviceID string `json:"deviceId"`
		AuthID   string `json:"authId"`
		Password string `json:"password"`
	}

	cfg := &connectorConfig{
		CaCert:   certFile,
		LogFile:  logFile,
		Address:  c2eCfg.MqttAdapterAddress,
		DeviceID: deviceId,
		TenantID: tenantId,
		AuthID:   authId,
		Password: devicePass,
	}

	jsonContents, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return fmt.Errorf("unable to marshal to json: %v", err)
	}

	// Preserve the file mode if the file already exists
	mode := getFileModeOrDefault(path, configDefaultMode)
	err = ioutil.WriteFile(path, jsonContents, mode)
	if err != nil {
		return fmt.Errorf("unable to save suite-connector configuration json file %s: %v", path, err)
	}
	return nil
}

func getFileModeOrDefault(path string, defaultMode os.FileMode) os.FileMode {
	mode := defaultMode
	fileInfo, err := os.Stat(path)
	if err == nil {
		mode = fileInfo.Mode()
	}
	return mode
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

func getThingConfig(mqttClient MQTT.Client) (*thingConfig, error) {
	type result struct {
		cfg *thingConfig
		err error
	}

	ch := make(chan result)

	if token := mqttClient.Subscribe("edge/thing/response", 1, func(client MQTT.Client, message MQTT.Message) {
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
		return nil, fmt.Errorf("thing configuration not received in %v", timeout)
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
	// Preserve the file mode if the file already exists
	mode := getFileModeOrDefault(dst, configDefaultMode)
	return ioutil.WriteFile(dst, data, mode)
}
