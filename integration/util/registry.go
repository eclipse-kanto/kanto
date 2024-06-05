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
	"strings"
)

// Resource holds all needed properties to create resources for the device.
type Resource struct {
	URL string

	Method string
	Body   string

	User string
	Pass string

	Delete bool
}

// ConnectorConfiguration holds the minimum required configuration to suite connector to connect.
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
func CreateDeviceResources(newDeviceID, tenantID, policyID, password, registryAPI,
	registryAPIUsername, registryAPIPassword string, cfg *TestConfiguration) []*Resource {

	devicePath := tenantID + "/" + newDeviceID
	return []*Resource{
		&Resource{
			URL:    registryAPI + "/devices/" + devicePath,
			Method: http.MethodPost,
			Body:   `{"authorities":["auto-provisioning-enabled"]}`,
			User:   registryAPIUsername,
			Pass:   registryAPIPassword,
			Delete: true},
		&Resource{
			URL:    registryAPI + "/credentials/" + devicePath,
			Method: http.MethodPut,
			Body:   getCredentialsBody(strings.ReplaceAll(newDeviceID, ":", "_"), password),
			User:   registryAPIUsername,
			Pass:   registryAPIPassword},
		&Resource{
			URL:    GetThingURL(cfg.DigitalTwinAPIAddress, newDeviceID),
			Method: http.MethodPut,
			Body:   fmt.Sprintf(`{"policyId": "%s"}`, policyID),
			User:   cfg.DigitalTwinAPIUsername,
			Pass:   cfg.DigitalTwinAPIPassword,
			Delete: true},
	}
}

func getCredentialsBody(authID, pass string) string {
	type pwdPlain struct {
		PwdPlain string `json:"pwd-plain"`
	}

	type authStruct struct {
		TypeStr string     `json:"type"`
		AuthID  string     `json:"auth-id"`
		Secrets []pwdPlain `json:"secrets"`
	}
	auth := authStruct{"hashed-password", authID, []pwdPlain{pwdPlain{pass}}}

	data, _ := json.MarshalIndent([]authStruct{auth}, "", "\t")
	return string(data)
}

// RegisterDeviceResources registers all given resources. In case of error all resources registered by this function will be deleted.
func RegisterDeviceResources(cfg *TestConfiguration,
	resources []*Resource, deviceID, url, user, pass string) error {
	for i, r := range resources {
		if _, err := SendDeviceRegistryRequest(([]byte)(r.Body), r.Method, r.URL, r.User, r.Pass); err != nil {
			if i > 0 {
				DeleteResources(cfg, resources[:i], deviceID, url, user, pass)
			}
			return err
		}
	}
	return nil
}

// DeleteResources deletes all given resources and all related devices.
func DeleteResources(cfg *TestConfiguration, resources []*Resource, deviceID, url, user, pass string) error {
	var errors []error
	if err := deleteRelatedDevices(cfg, deviceID, url, user, pass); err != nil {
		errors = append(errors, err)
	}

	// Delete in reverse order of creation
	for i := len(resources) - 1; i >= 0; i-- {
		r := resources[i]

		if r.Delete {
			if _, err := SendDeviceRegistryRequest(nil, http.MethodDelete, r.URL, r.User, r.Pass); err != nil {
				errors = append(errors, err)
			}
		}

	}
	return CombineErrors(errors)
}

func deleteRelatedDevices(cfg *TestConfiguration, viaDeviceID, url, user, pass string) error {
	devicesVia, err := findDeviceRegistryDevicesVia(viaDeviceID, url, user, pass)
	if err != nil {
		return err
	}

	var errors []error
	// Digital Twin API things are created after Device Registry devices, so delete them first
	if err = deleteDigitalTwinThings(cfg, devicesVia); err != nil {
		errors = append(errors, err)
	}
	// Then delete Device Registry devices
	if err = deleteRegistryDevices(devicesVia, url, user, pass); err != nil {
		errors = append(errors, err)
	}
	return CombineErrors(errors)
}

func findDeviceRegistryDevicesVia(viaDeviceID, url, user, pass string) ([]string, error) {
	type registryDevice struct {
		ID  string   `json:"id"`
		Via []string `json:"via"`
	}
	type registryDevices struct {
		Devices []*registryDevice `json:"result"`
	}
	devicesJSON, err := SendDeviceRegistryRequest(nil, http.MethodGet, url, user, pass)
	if err != nil {
		return nil, err
	}
	devices := &registryDevices{}
	err = json.Unmarshal(devicesJSON, devices)
	if err != nil {
		return nil, err
	}
	var devicesVia []string
	for _, device := range devices.Devices {
		for _, via := range device.Via {
			if via == viaDeviceID {
				devicesVia = append(devicesVia, device.ID)
				break
			}
		}
	}

	return devicesVia, nil
}

func deleteDigitalTwinThings(cfg *TestConfiguration, things []string) error {
	var errors []error
	for _, thingID := range things {
		if _, err := SendDigitalTwinRequest(
			cfg, http.MethodDelete, GetThingURL(cfg.DigitalTwinAPIAddress, thingID), nil); err != nil {
			errors = append(errors, err)
		}
	}
	return CombineErrors(errors)
}

func deleteRegistryDevices(devices []string, tenantURL, user, pass string) error {
	var errors []error
	for _, device := range devices {
		if _, err := SendDeviceRegistryRequest(nil, http.MethodDelete, tenantURL+device, user, pass); err != nil {
			errors = append(errors, err)
		}
	}
	return CombineErrors(errors)
}
