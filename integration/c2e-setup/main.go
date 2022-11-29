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

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/eclipse-kanto/kanto/integration/util"
)

const (
	devices     = "devices/"
	credentials = "credentials/"

	indent = " "

	systemctl             = "systemctl"
	restart               = "restart"
	suiteConnectorService = "suite-connector.service"

	configDefaultMode = 0666
)

var (
	cfg    util.TestConfiguration
	c2eCfg c2eConfiguration

	caCert   string
	logFile  string
	deviceID string
	tenantID string
	password string

	policyID string

	configFile       string
	configFileBackup string

	authID string
)

type resource struct {
	url string

	method string
	body   string

	user string
	pass string

	delete bool
}

type c2eConfiguration struct {
	DeviceRegistryAPIAddress string `env:"DEVICE_REGISTRY_API_ADDRESS"`

	DeviceRegistryAPIUsername string `env:"DEVICE_REGISTRY_API_USERNAME" envDefault:"ditto"`
	DeviceRegistryAPIPassword string `env:"DEVICE_REGISTRY_API_PASSWORD" envDefault:"ditto"`

	MqttAdapterAddress string `env:"MQTT_ADAPTER_ADDRESS"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&caCert, "caCert", "/etc/suite-connector/iothub.crt", "Path to Suite Connector CA certificates file")
	flag.StringVar(&logFile, "logFile", "/var/log/suite-connector/suite-connector.log", "Path to Suite Connector log file")
	flag.StringVar(&deviceID, "deviceId", "", "Test device unique identifier, defaults to randomly generated")
	flag.StringVar(&tenantID, "tenantId", "", "Device registry tenant unique identifier")
	flag.StringVar(&password, "password", "123456", "Test device password")

	flag.StringVar(&policyID, "policyId", "", "Test device's policy unique identifier")

	flag.StringVar(&configFile, "configFile", "/etc/suite-connector/config.json",
		"Path to Suite Connector configuration file. "+
			"If set to the empty string, configuring Suite Connector and restarting it will be skipped")

	flag.StringVar(&configFileBackup, "configFileBackup", "/etc/suite-connector/configBackup.json",
		"Path to Suite Connector configuration file backup. "+
			"If set to the empty string, backing up the Suite Connector configuration file will be skipped")

	clean := flag.Bool("clean", false, "Clean up test resources")

	flag.Parse()

	envOpts := env.Options{RequiredIfNoDef: true}
	err := env.Parse(&cfg, envOpts)
	if err == nil {
		err = env.Parse(&c2eCfg, envOpts)
	}
	if err != nil {
		fmt.Printf("failed to process environment variables: %v\n", err)
		printConfigHelp()
		os.Exit(1)
	}

	if !*clean {
		if deviceID == "" {
			deviceID = generateRandomDeviceID()
			fmt.Printf("generating a random device id: \"%s\"\n", deviceID)
		} else {
			fmt.Printf("forcing device id: \"%s\"\n", deviceID)
		}
	} else if deviceID == "" || tenantID == "" {
		mqttClient, err := util.NewMQTTClient(&cfg)
		if err != nil {
			fmt.Printf("unable to open local MQTT connection to %s, error: %v\n", cfg.LocalBroker, err)
			os.Exit(1)
		}
		defer mqttClient.Disconnect(uint(cfg.MqttQuiesceMs))
		thingConfiguration, err := util.GetThingConfiguration(&cfg, mqttClient)
		if err != nil {
			fmt.Printf("unable to get thing configuration from the local MQTT %s, error: %v\n", cfg.LocalBroker, err)
			os.Exit(1)
		}
		deviceID = thingConfiguration.DeviceID
		tenantID = thingConfiguration.TenantID
	}

	if !*clean {
		assertFlag(tenantID, "tenant id")
		assertFlag(policyID, "policy id")
	}

	registryAPI := strings.TrimSuffix(c2eCfg.DeviceRegistryAPIAddress, "/") + "/v1"
	devicePath := tenantID + "/" + deviceID

	resources := make([]*resource, 0, 3)
	deviceResource := &resource{url: registryAPI + "/" + devices + devicePath, method: http.MethodPost,
		body: deviceJSON, user: c2eCfg.DeviceRegistryAPIUsername, pass: c2eCfg.DeviceRegistryAPIPassword, delete: true}
	resources = append(resources, deviceResource)

	authID = strings.ReplaceAll(deviceID, ":", "_")
	auth := fmt.Sprintf(authJSON, authID, password)
	resources = append(resources, &resource{url: registryAPI + "/" + credentials + devicePath, method: http.MethodPut,
		body: auth, user: c2eCfg.DeviceRegistryAPIUsername, pass: c2eCfg.DeviceRegistryAPIPassword, delete: false})

	thingURL := util.GetThingURL(cfg.DigitalTwinAPIAddress, deviceID)
	thing := fmt.Sprintf(thingJSON, policyID)
	resources = append(resources, &resource{url: thingURL, method: http.MethodPut,
		body: thing, user: cfg.DigitalTwinAPIUsername, pass: cfg.DigitalTwinAPIPassword, delete: true})

	var ok bool
	if !*clean {
		ok = performSetUp(deviceResource, resources)
		fmt.Println("setup complete")
	} else {
		ok = performCleanUp(resources)
		fmt.Println("cleanup complete")
	}
	if ok {
		os.Exit(0)
	}
	os.Exit(1)
}

func printHelp(cfg interface{}) {
	t := reflect.TypeOf(cfg)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		name, ok := f.Tag.Lookup("env")
		if ok {

			fmt.Printf("\n\t - %s", name)

			def, ok := f.Tag.Lookup("envDefault")
			if ok {
				fmt.Printf(" (default '%s')", def)
			}
		}
	}
}

func printConfigHelp() {
	fmt.Print("config environmental variables:")
	printHelp(cfg)
	printHelp(c2eCfg)
	fmt.Println()
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

func isDeviceIDPresentInRegistry(deviceID string, deviceResource *resource) bool {
	_, err := sendDeviceRegistryRequest(http.MethodGet, deviceResource.url)
	if err == nil {
		return true
	}
	return false
}

func performSetUp(deviceResource *resource, resources []*resource) bool {
	if isDeviceIDPresentInRegistry(deviceID, deviceResource) {
		fmt.Printf("device %s already exists in registry, aborting...\n", deviceID)
		return false
	}

	fmt.Println("performing setup...")

	if configFile != "" && configFileBackup != "" {
		fmt.Println("saving a backup of the suite-connector configuration file...")
		if err := copyFile(configFile, configFileBackup); err != nil {
			fmt.Printf(
				"unable to save backup copy of configuration file %s to %s: %v\n",
				configFile, configFileBackup, err)
			return false
		}
	}

	for i, r := range resources {
		if b, err := sendRequest(r.method, r.url, ([]byte)(r.body), r.user, r.pass); err != nil {
			fmt.Printf("unable to create device at %s, error: %v\n", r.url, err)

			if b != nil {
				fmt.Println(string(b))
				fmt.Println()
			}

			if i > 0 {
				deleteResources(resources[:i])
			}
			deleteBackupFile()
			return false
		}
		fmt.Printf("%s '%s' created\n", indent, r.url)
	}

	ok := true
	if configFile != "" {
		if err := writeConfigFile(configFile); err != nil {
			fmt.Printf("unable to write configuration file, error: %v\n", err)

			deleteResources(resources)
			deleteBackupFile()
			return false
		}

		fmt.Printf("%s configuration file '%s' written\n", indent, configFile)
		ok = restartSuiteConnector()
	}

	if ok {
		fmt.Println("setup successful")
	}

	return ok
}

func performCleanUp(resources []*resource) bool {
	ok := true
	if configFile != "" && configFileBackup != "" {
		fmt.Println("restoring suite-connector configuration file and restarting suite-connector")
		if err := copyFile(configFileBackup, configFile); err != nil {
			fmt.Printf(
				"unable to restore the backup copy of configuration file %s to %s: %v\n",
				configFileBackup, configFile, err)
			ok = false
		} else {
			ok = restartSuiteConnector()
		}
		if ok {
			// Delete suite-connector configuration backup file
			deleteBackupFile()
		}
	}
	// Delete devices and things
	fmt.Printf("performing cleanup on device id: %s\n", deviceID)
	deleteResources(resources)
	return ok
}

func deleteBackupFile() {
	if err := os.Remove(configFileBackup); err != nil {
		fmt.Printf("unable to delete configuration file backup %s, error: %v", configFileBackup, err)
	}
}

func deleteResources(resources []*resource) {
	deleteRelatedDevices(deviceID)
	fmt.Println("deleting initially created things...")
	// Delete in reverse order of creation
	for i := len(resources) - 1; i >= 0; i-- {
		r := resources[i]

		if !r.delete {
			continue
		}

		if _, err := sendRequest(http.MethodDelete, r.url, nil, r.user, r.pass); err != nil {
			fmt.Printf("%s unable to delete '%s', error: %v\n", indent, r.url, err)
		} else {
			fmt.Printf("%s '%s' deleted\n", indent, r.url)
		}
	}
}

func deleteRelatedDevices(viaDeviceID string) {
	devicesVia := findDeviceRegistryDevicesVia(viaDeviceID)
	// Digital Twin API things are created after Device Registry devices, so delete them first
	fmt.Println("deleting automatically created things in the digital twin API...")
	deleteDigitalTwinThings(devicesVia)
	// Then delete Device Registry devices
	fmt.Println("deleting automatically created devices in the device registry...")
	deleteDeviceRegistryDevices(devicesVia)
}

func findDeviceRegistryDevicesVia(viaDeviceID string) []string {
	var relations []string

	type registryDevice struct {
		ID  string   `json:"id"`
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

	devicesJSON, err := sendDeviceRegistryRequest(http.MethodGet, getTenantURL())
	if err != nil {
		fmt.Printf("unable to list devices from the device registry, error: %v\n", err)
	} else {
		devices := &registryDevices{}
		err = json.Unmarshal(devicesJSON, devices)
		if err != nil {
			fmt.Printf("unable to parse devices JSON returned from the device registry, error: %v\n", err)
			devices.Devices = nil
		}
		for _, device := range devices.Devices {
			if contains(device.Via, viaDeviceID) {
				relations = append(relations, device.ID)
			}
		}
	}

	return relations
}

func deleteDigitalTwinThings(things []string) {
	// Delete Digital Twin things
	for _, thingID := range things {
		url := util.GetThingURL(cfg.DigitalTwinAPIAddress, thingID)
		_, err := util.SendDigitalTwinRequest(&cfg, http.MethodDelete, url, nil)
		if err != nil {
			fmt.Printf("error deleting thing: %v\n", err)
		} else {
			fmt.Printf("%s '%s' deleted\n", indent, url)
		}
	}
}

func deleteDeviceRegistryDevices(devices []string) {
	// Delete Device Registry devices
	tenantURL := getTenantURL()
	for _, device := range devices {
		url := tenantURL + device
		if _, err := sendDeviceRegistryRequest(http.MethodDelete, url); err != nil {
			fmt.Printf("error deleting device: %v\n", err)
		} else {
			fmt.Printf("%s '%s' deleted\n", indent, url)
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
		CaCert:   caCert,
		LogFile:  logFile,
		Address:  c2eCfg.MqttAdapterAddress,
		DeviceID: deviceID,
		TenantID: tenantID,
		AuthID:   authID,
		Password: password,
	}

	jsonContents, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return fmt.Errorf("unable to marshal to json: %v", err)
	}

	// Preserve the file mode if the file already exists
	mode := getFileModeOrDefault(path, configDefaultMode)
	err = os.WriteFile(path, jsonContents, mode)
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

func sendDeviceRegistryRequest(method string, url string) ([]byte, error) {
	return sendRequest(method, url, nil, c2eCfg.DeviceRegistryAPIUsername, c2eCfg.DeviceRegistryAPIPassword)
}

func sendRequest(method string, url string, payload []byte, username string, password string) ([]byte, error) {
	var reqBody io.Reader

	if payload != nil {
		// Unlike util.SendDigitalTwinRequest, we use the payload directly
		reqBody = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s request failed: %s", method, url, resp.Status)
	}

	return body, err
}

func restartSuiteConnector() bool {
	fmt.Println("restarting suite-connector...")
	cmd := exec.Command(systemctl, restart, suiteConnectorService)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("error restarting %s: %v", suiteConnectorService, err)
		return false
	}
	fmt.Println(string(stdout))
	fmt.Println("... done")
	return true
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	// If the destination file exists, preserve its file mode.
	// If the destination file doesn't exist, use the file mode of the source file.
	srcMode := getFileModeOrDefault(src, configDefaultMode)
	dstMode := getFileModeOrDefault(dst, srcMode)
	return os.WriteFile(dst, data, dstMode)
}

func getTenantURL() string {
	return fmt.Sprintf(
		"%s/v1/%s%s/", strings.TrimSuffix(c2eCfg.DeviceRegistryAPIAddress, "/"), devices, tenantID)
}
