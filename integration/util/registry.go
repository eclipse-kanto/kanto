// Copyright (c) 2023 Contributors to the Eclipse Foundation
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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	Indent = " "

	deletedTemplate = "%s '%s' deleted\n"
	done            = "... done"

	deviceJSON = `{"authorities":["auto-provisioning-enabled"]}`

	authJSON = `[{
		"type": "hashed-password",
		"auth-id": "%s",
		"secrets": [{
			"pwd-plain": "%s"
		}]
	}]`

	thingJSON = `{"policyId": "%s"}`

	configDefaultMode = 0666

	systemctl = "systemctl"
)

type Resource struct {
	url string

	method string
	body   string

	user string
	pass string

	delete bool
}

type BootstrapConfiguration struct {
	LogFile             string   `json:"logFile"`
	PostBootstrapFile   string   `json:"postBootstrapFile"`
	PostBootstrapScript []string `json:"postBootstrapScript"`
	CaCert              string   `json:"caCert"`
	Address             string   `json:"address"`
	TenantID            string   `json:"tenantId"`
	DeviceID            string   `json:"deviceId"`
	AuthID              string   `json:"authId"`
	Password            string   `json:"password"`
}

type ConnectorConfiguration struct {
	CaCert   string `json:"caCert"`
	LogFile  string `json:"logFile"`
	Address  string `json:"address"`
	TenantID string `json:"tenantId"`
	DeviceID string `json:"deviceId"`
	AuthID   string `json:"authId"`
	Password string `json:"password"`
}

// CreateDeviceResources creates device resources.
func CreateDeviceResources(newDeviceId string, resources []*Resource,
	tenantID, policyID, password, registryAPI, registryAPIUsername, registryAPIPassword string,
	cfg *TestConfiguration) (*Resource, []*Resource) {

	devicePath := tenantID + "/" + newDeviceId
	deviceResource := &Resource{url: registryAPI + "/devices/" + devicePath, method: http.MethodPost,
		body: deviceJSON, user: registryAPIUsername, pass: registryAPIPassword, delete: true}
	resources = append(resources, deviceResource)

	authID := strings.ReplaceAll(newDeviceId, ":", "_")
	auth := fmt.Sprintf(authJSON, authID, password)
	resources = append(resources, &Resource{url: registryAPI + "/credentials/" + devicePath, method: http.MethodPut,
		body: auth, user: registryAPIUsername, pass: registryAPIPassword, delete: false})

	thingURL := GetThingURL(cfg.DigitalTwinAPIAddress, newDeviceId)
	thing := fmt.Sprintf(thingJSON, policyID)
	resources = append(resources, &Resource{url: thingURL, method: http.MethodPut,
		body: thing, user: cfg.DigitalTwinAPIUsername, pass: cfg.DigitalTwinAPIPassword, delete: true})
	return deviceResource, resources
}

// WriteConfigFile writes interface data to the path file, creating it if necessary.
func WriteConfigFile(path string, cfg interface{}) error {
	jsonContents, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return fmt.Errorf("unable to marshal to json: %v", err)
	}

	// Preserve the file mode if the file already exists
	mode := getFileModeOrDefault(path, configDefaultMode)
	err = os.WriteFile(path, jsonContents, mode)
	if err != nil {
		return fmt.Errorf("unable to save file %s: %v", path, err)
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

// DeleteResources deletes all given resources and all related devices.
func DeleteResources(cfg *TestConfiguration, resources []*Resource, deviceId, url, user, pass string) bool {
	ok := deleteRelatedDevices(cfg, deviceId, url, user, pass)
	fmt.Println("deleting initially created things...")
	// Delete in reverse order of creation
	for i := len(resources) - 1; i >= 0; i-- {
		r := resources[i]

		if !r.delete {
			continue
		}

		if _, err := SendDeviceRegistryRequest(nil, http.MethodDelete, r.url, r.user, r.pass); err != nil {
			ok = false
			fmt.Printf("%s unable to delete '%s', error: %v\n", Indent, r.url, err)
		} else {
			fmt.Printf(deletedTemplate, Indent, r.url)
		}
	}
	return ok
}

func deleteRelatedDevices(cfg *TestConfiguration, viaDeviceID, url, user, pass string) bool {
	devicesVia, ok := findDeviceRegistryDevicesVia(viaDeviceID, url, user, pass)
	// Digital Twin API things are created after Device Registry devices, so delete them first
	fmt.Println("deleting automatically created things...")
	if !deleteDigitalTwinThings(cfg, devicesVia) {
		ok = false
	}
	// Then delete Device Registry devices
	fmt.Println("deleting automatically created devices...")
	if !deleteRegistryDevices(devicesVia, url, user, pass) {
		ok = false
	}
	return ok
}

func findDeviceRegistryDevicesVia(viaDeviceID, url, user, pass string) ([]string, bool) {
	var devicesVia []string
	ok := true

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

	devicesJSON, err := SendDeviceRegistryRequest(nil, http.MethodGet, url, user, pass)
	if err != nil {
		ok = false
		fmt.Printf("unable to list devices from the device registry, error: %v\n", err)
	} else {
		devices := &registryDevices{}
		err = json.Unmarshal(devicesJSON, devices)
		if err != nil {
			ok = false
			fmt.Printf("unable to parse devices JSON returned from the device registry, error: %v\n", err)
			devices.Devices = nil
		}
		for _, device := range devices.Devices {
			if contains(device.Via, viaDeviceID) {
				devicesVia = append(devicesVia, device.ID)
			}
		}
	}

	return devicesVia, ok
}

func deleteDigitalTwinThings(cfg *TestConfiguration, things []string) bool {
	ok := true
	for _, thingID := range things {
		url := GetThingURL(cfg.DigitalTwinAPIAddress, thingID)
		_, err := SendDigitalTwinRequest(cfg, http.MethodDelete, url, nil)
		if err != nil {
			ok = false
			fmt.Printf("error deleting thing: %v\n", err)
		} else {
			fmt.Printf(deletedTemplate, Indent, url)
		}
	}
	return ok
}

func deleteRegistryDevices(devices []string, tenantURL, user, pass string) bool {
	ok := true
	for _, device := range devices {
		url := tenantURL + device
		if _, err := SendDeviceRegistryRequest(nil, http.MethodDelete, url, user, pass); err != nil {
			ok = false
			fmt.Printf("error deleting device: %v\n", err)
		} else {
			fmt.Printf(deletedTemplate, Indent, url)
		}
	}
	return ok
}

// RestartService restarts the service with given name
func RestartService(service string) bool {
	fmt.Printf("restarting %s...", service)
	cmd := exec.Command(systemctl, "restart", service)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("error restarting %s: %v", service, err)
		return false
	}
	fmt.Println(string(stdout))
	fmt.Println(done)
	return true
}

// StopService stops the service with given name
func StopService(service string) bool {
	fmt.Printf("stopping %s...", service)
	cmd := exec.Command(systemctl, "stop", service)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("error stopping %s: %v", service, err)
		return false
	}
	fmt.Println(string(stdout))
	fmt.Println(done)
	return true
}

// CopyFile copies source file to the destination.
func CopyFile(src, dst string) error {
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

// DeleteFile removes the named file or directory.
func DeleteFile(path string) bool {
	if err := os.Remove(path); err != nil {
		fmt.Printf("unable to delete file %s, error: %v", path, err)
		return false
	}
	return true
}
