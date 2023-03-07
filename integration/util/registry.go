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
	deviceJSON = `{"authorities":["auto-provisioning-enabled"]}`

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
func CreateDeviceResources(newDeviceId, tenantID, policyID, password, registryAPI,
	registryAPIUsername, registryAPIPassword string,
	cfg *TestConfiguration) (*Resource, []*Resource) {

	resources := make([]*Resource, 0, 3)
	devicePath := tenantID + "/" + newDeviceId
	deviceResource := &Resource{url: registryAPI + "/devices/" + devicePath, method: http.MethodPost,
		body: deviceJSON, user: registryAPIUsername, pass: registryAPIPassword, delete: true}
	resources = append(resources, deviceResource)

	authID := strings.ReplaceAll(newDeviceId, ":", "_")
	resources = append(resources, &Resource{url: registryAPI + "/credentials/" + devicePath, method: http.MethodPut,
		body: getCredentialsBody(authID, password), user: registryAPIUsername, pass: registryAPIPassword, delete: false})

	thingURL := GetThingURL(cfg.DigitalTwinAPIAddress, newDeviceId)
	thing := fmt.Sprintf(thingJSON, policyID)
	resources = append(resources, &Resource{url: thingURL, method: http.MethodPut,
		body: thing, user: cfg.DigitalTwinAPIUsername, pass: cfg.DigitalTwinAPIPassword, delete: true})
	return deviceResource, resources
}

func getCredentialsBody(authID, pass string) string {
	type pwdPlain struct {
		PwdPlain string `json:"pwd-plain"`
	}

	type authStruct struct {
		TypeStr string     `json:"type"`
		AuthId  string     `json:"auth-id"`
		Secrets []pwdPlain `json:"secrets"`
	}
	pwds := []pwdPlain{pwdPlain{pass}}

	auth := authStruct{"hashed-password", authID, pwds}
	authJson := []authStruct{auth}
	data, _ := json.MarshalIndent(authJson, "", "\t")
	return string(data)
}

// WriteConfigFile writes interface data to the path file, creating it if necessary.
func WriteConfigFile(path string, cfg interface{}) error {
	jsonContents, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return fmt.Errorf("unable to marshal to json: %v", err)
	}

	// Preserve the file mode if the file already exists
	mode, err := getFileModeOrDefault(path, configDefaultMode)
	if err != nil {
		return fmt.Errorf("unable to get file mode %s: %v", path, err)
	}
	if err = os.WriteFile(path, jsonContents, mode); err != nil {
		return fmt.Errorf("unable to save file %s: %v", path, err)
	}
	return nil
}

func getFileModeOrDefault(path string, defaultMode os.FileMode) (os.FileMode, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return defaultMode, err
	}
	return fileInfo.Mode(), nil
}

// DeleteResources deletes all given resources and all related devices.
func DeleteResources(cfg *TestConfiguration, resources []*Resource, deviceId, url, user, pass string) []error {
	deleteErrors := deleteRelatedDevices(cfg, deviceId, url, user, pass)
	// Delete in reverse order of creation
	for i := len(resources) - 1; i >= 0; i-- {
		r := resources[i]

		if !r.delete {
			continue
		}

		if _, err := SendDeviceRegistryRequest(nil, http.MethodDelete, r.url, r.user, r.pass); err != nil {
			deleteErrors = append(deleteErrors, err)
		}
	}
	return deleteErrors
}

func deleteRelatedDevices(cfg *TestConfiguration, viaDeviceID, url, user, pass string) []error {
	var deleteErrors []error
	devicesVia, err := findDeviceRegistryDevicesVia(viaDeviceID, url, user, pass)
	if err != nil {
		deleteErrors = append(deleteErrors, err)
	}
	// Digital Twin API things are created after Device Registry devices, so delete them first
	if delErrors := deleteDigitalTwinThings(cfg, devicesVia); delErrors != nil {
		deleteErrors = append(deleteErrors, delErrors...)
	}
	// Then delete Device Registry devices
	if delErrors := deleteRegistryDevices(devicesVia, url, user, pass); delErrors != nil {
		deleteErrors = append(deleteErrors, delErrors...)
	}
	return deleteErrors
}

func findDeviceRegistryDevicesVia(viaDeviceID, url, user, pass string) ([]string, error) {
	var devicesVia []string

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
		return devicesVia, err
	} else {
		devices := &registryDevices{}
		err = json.Unmarshal(devicesJSON, devices)
		if err != nil {
			return devicesVia, err
		}
		for _, device := range devices.Devices {
			if contains(device.Via, viaDeviceID) {
				devicesVia = append(devicesVia, device.ID)
			}
		}
	}

	return devicesVia, nil
}

func deleteDigitalTwinThings(cfg *TestConfiguration, things []string) []error {
	var deleteErrors []error
	for _, thingID := range things {
		url := GetThingURL(cfg.DigitalTwinAPIAddress, thingID)
		if _, err := SendDigitalTwinRequest(cfg, http.MethodDelete, url, nil); err != nil {
			deleteErrors = append(deleteErrors, err)
		}
	}
	return deleteErrors
}

func deleteRegistryDevices(devices []string, tenantURL, user, pass string) []error {
	var deleteErrors []error
	for _, device := range devices {
		url := tenantURL + device
		if _, err := SendDeviceRegistryRequest(nil, http.MethodDelete, url, user, pass); err != nil {
			deleteErrors = append(deleteErrors, err)
		}
	}
	return deleteErrors
}

// ExecuteCommandToService executes command to the service with given name
func ExecuteCommandToService(command, service string) ([]byte, error) {
	cmd := exec.Command(systemctl, command, service)
	return cmd.Output()
}

// CopyFile copies source file to the destination.
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	// If the destination file exists, preserve its file mode.
	// If the destination file doesn't exist, use the file mode of the source file.
	srcMode, err := getFileModeOrDefault(src, configDefaultMode)
	if err != nil {
		return err
	}
	dstMode, err := getFileModeOrDefault(dst, srcMode)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, dstMode)
}
