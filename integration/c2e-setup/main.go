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
	"flag"
	"fmt"
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
	indent = " "

	systemctl = "systemctl"
	restart   = "restart"
	stop      = "stop"

	suiteConnectorService     = "suite-connector.service"
	suiteBootstrappingService = "suite-bootstrapping.service"
	ldtService                = "local-digital-twins.service"

	deleteResourcesTemplate = "%s unable to delete resources, error: %v\n"
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

	configConnectorFile       string
	configConnectorFileBackup string

	authID string

	bootstrapCaCert           string
	logBootstrapFile          string
	postBootstrapFile         string
	postBootstrapScript       string
	configBootstrapFile       string
	configBootstrapFileBackup string

	ldtCaCert           string
	logLdtFile          string
	thingsDb            string
	configLdtFile       string
	configLdtFileBackup string

	bootstrap bool
	ldt       bool
)

type c2eConfiguration struct {
	DeviceRegistryAPIAddress string `env:"DEVICE_REGISTRY_API_ADDRESS"`

	DeviceRegistryAPIUsername string `env:"DEVICE_REGISTRY_API_USERNAME" envDefault:"ditto"`
	DeviceRegistryAPIPassword string `env:"DEVICE_REGISTRY_API_PASSWORD" envDefault:"ditto"`

	MQTTAdapterAddress string `env:"MQTT_ADAPTER_ADDRESS"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&caCert, "caCert", "/etc/suite-connector/iothub.crt", "Path to Suite Connector CA certificates file")
	flag.StringVar(&logFile, "logFile", "/var/log/suite-connector/suite-connector.log", "Path to Suite Connector log file")
	flag.StringVar(&deviceID, "deviceId", "", "Test device unique identifier, defaults to randomly generated")
	flag.StringVar(&tenantID, "tenantId", "", "Device registry tenant unique identifier")
	flag.StringVar(&password, "password", "123456", "Test device password")

	flag.StringVar(&policyID, "policyId", "", "Test device's policy unique identifier")

	flag.StringVar(&configConnectorFile, "configFile", "/etc/suite-connector/config.json",
		"Path to Suite Connector configuration file. "+
			"If set to the empty string, configuring Suite Connector and restarting it will be skipped")

	flag.StringVar(&configConnectorFileBackup, "configFileBackup", "/etc/suite-connector/configBackup.json",
		"Path to Suite Connector configuration file backup. "+
			"If set to the empty string, backing up the Suite Connector configuration file will be skipped")

	flag.StringVar(&bootstrapCaCert, "bootstrapCaCert", "/etc/suite-bootstrapping/iothub.crt", "Path to Suite Bootstrapping CA certificates file")
	flag.StringVar(&logBootstrapFile, "logBootstrapFile", "/var/log/suite-bootstrapping/suite-bootstrapping.log",
		"Path to Suite Bootstrapping log file")
	flag.StringVar(&postBootstrapFile, "postBootstrapFile", "/etc/suite-connector/config.json",
		"Path to the file used for a bootstrapping response data")
	flag.StringVar(&postBootstrapScript, "postBootstrapScript", "/var/tmp/suite-bootstrapping/post_script.sh",
		"Path to the script that is executed after a bootstrapping response")

	flag.StringVar(&configBootstrapFile, "configBootstrapFile", "/etc/suite-bootstrapping/config.json",
		"Path to Suite Bootstrapping configuration file. "+
			"If set to the empty string, configuring Suite Bootstrapping and restarting it will be skipped")

	flag.StringVar(&configBootstrapFileBackup, "configBootstrapFileBackup", "/etc/suite-bootstrapping/configBackup.json",
		"Path to Suite Bootstrapping configuration file backup. "+
			"If set to the empty string, backing up the Suite Bootstrapping configuration file will be skipped")

	flag.StringVar(&ldtCaCert, "ldtCaCert", "/etc/local-digital-twins/iothub.crt", "Path to Local Digital Twins CA certificates file")
	flag.StringVar(&logLdtFile, "logLdtFile", "/var/log/local-digital-twins/local-digital-twins.log",
		"Path to Local Digital Twins log file")
	flag.StringVar(&thingsDb, "thingsDb", "/var/lib/local-digital-twins/thing.db",
		"Path to the file where digital twins will be stored")

	flag.StringVar(&configLdtFile, "configLdtFile", "/etc/local-digital-twins/config.json",
		"Path to Local Digital Twins configuration file. "+
			"If set to the empty string, configuring Local Digital Twins and restarting it will be skipped")

	flag.StringVar(&configLdtFileBackup, "configLdtFileBackup", "/etc/local-digital-twins/configBackup.json",
		"Path to Local Digital Twins configuration file backup. "+
			"If set to the empty string, backing up the Local Digital Twins configuration file will be skipped")

	clean := flag.Bool("clean", false, "Clean up test resources")
	flag.BoolVar(&bootstrap, "bootstrap", false, "Create bootstrapping resources")
	flag.BoolVar(&ldt, "ldt", false, "Create local-digital-twins resources")

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
		assertFlag(tenantID, "tenant id")
		assertFlag(policyID, "policy id")
	} else if deviceID == "" || tenantID == "" {
		mqttClient, err := util.NewMQTTClient(&cfg)
		if err != nil {
			fmt.Printf("unable to open local MQTT connection to %s, error: %v\n", cfg.LocalBroker, err)
			os.Exit(1)
		}
		defer mqttClient.Disconnect(uint(cfg.MQTTQuiesceMS))
		thingConfiguration, err := util.GetThingConfiguration(&cfg, mqttClient)
		if err != nil {
			fmt.Printf("unable to get thing configuration from the local MQTT %s, error: %v\n", cfg.LocalBroker, err)
			os.Exit(1)
		}
		deviceID = thingConfiguration.DeviceID
		tenantID = thingConfiguration.TenantID
	}

	authID = strings.ReplaceAll(deviceID, ":", "_")
	registryAPI := strings.TrimSuffix(c2eCfg.DeviceRegistryAPIAddress, "/") + "/v1"

	resources := util.CreateDeviceResources(deviceID, tenantID, policyID, password, registryAPI,
		c2eCfg.DeviceRegistryAPIUsername, c2eCfg.DeviceRegistryAPIPassword, &cfg)

	var ok bool
	if !*clean {
		ok = performSetUp(resources)
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

func isDeviceIDPresentInRegistry(deviceResource *util.Resource) bool {
	_, err := util.SendDeviceRegistryRequest(nil, http.MethodGet,
		deviceResource.URL, c2eCfg.DeviceRegistryAPIUsername, c2eCfg.DeviceRegistryAPIPassword)
	return err == nil
}

func getServiceNameConfigAndBackupFile() (string, string, string) {
	if bootstrap {
		return suiteBootstrappingService, configBootstrapFile, configBootstrapFileBackup
	} else if ldt {
		return ldtService, configLdtFile, configLdtFileBackup
	}
	return suiteConnectorService, configConnectorFile, configConnectorFileBackup
}

func performSetUp(resources []*util.Resource) bool {
	serviceName, configFile, configFileBackup := getServiceNameConfigAndBackupFile()

	if len(resources) > 0 && isDeviceIDPresentInRegistry(resources[0]) {
		fmt.Printf("device %s already exists in registry, aborting...\n", deviceID)
		return false
	}

	fmt.Println("performing setup...")

	if configFile != "" && configFileBackup != "" {
		fmt.Printf("saving a backup of the %s configuration file...\n", serviceName)
		if err := util.CopyFile(configFile, configFileBackup); err != nil {
			fmt.Printf(
				"unable to save backup copy of configuration file %s to %s: %v\n",
				configFile, configFileBackup, err)
			return false
		}
	}

	for i, r := range resources {
		if b, err := util.SendDeviceRegistryRequest(([]byte)(r.Body), r.Method, r.URL, r.User, r.Pass); err != nil {
			fmt.Printf("unable to create device at %s, error: %v\n", r.URL, err)

			if b != nil {
				fmt.Println(string(b))
				fmt.Println()
			}

			if i > 0 {
				if err = util.DeleteResources(&cfg, resources[:i], deviceID, getTenantURL(),
					c2eCfg.DeviceRegistryAPIUsername, c2eCfg.DeviceRegistryAPIPassword); err != nil {
					fmt.Printf(deleteResourcesTemplate, indent, err)
				}
			}
			deleteFile(configFileBackup)
			return false
		}
		fmt.Printf("%s '%s' created\n", indent, r.URL)
	}

	ok := true
	if configFile != "" {
		var err error
		if serviceName == suiteBootstrappingService {
			err = writeConfigBootstrapFile(configFile)
		} else if serviceName == ldtService {
			err = writeConfigLdtFile(configFile)
		} else {
			err = writeConfigFile(configFile)
		}
		if err != nil {
			fmt.Printf("unable to write configuration file, error: %v\n", err)
			if err = util.DeleteResources(&cfg, resources, deviceID, getTenantURL(),
				c2eCfg.DeviceRegistryAPIUsername, c2eCfg.DeviceRegistryAPIPassword); err != nil {
				fmt.Printf(deleteResourcesTemplate, indent, err)
			}
			deleteFile(configFileBackup)
			return false
		}

		fmt.Printf("%s configuration file '%s' written\n", indent, configFile)

		if serviceName != suiteConnectorService {
			stdout, err := exec.Command(systemctl, stop, suiteConnectorService).Output()
			if stdout != nil {
				fmt.Println(string(stdout))
			}
			if err != nil {
				fmt.Printf("error stopping %s: %v", suiteConnectorService, err)
				return false
			}
		}
		ok = restartService(serviceName)
	}

	if ok {
		fmt.Println("setup successful")
	}

	return ok
}

func performCleanUp(resources []*util.Resource) bool {
	ok := true
	serviceName, configFile, configFileBackup := getServiceNameConfigAndBackupFile()

	if configFile != "" && configFileBackup != "" {
		fmt.Printf("restoring %s configuration file and restarting %s\n", serviceName, serviceName)
		if err := util.CopyFile(configFileBackup, configFile); err != nil {
			ok = false
			fmt.Printf(
				"unable to restore the backup copy of configuration file %s to %s: %v\n",
				configFileBackup, configFile, err)
		} else {
			if serviceName != suiteConnectorService {
				stdout, err := exec.Command(systemctl, stop, serviceName).Output()
				if stdout != nil {
					fmt.Println(string(stdout))
				}
				if err != nil {
					fmt.Printf("error stopping %s: %v", serviceName, err)
					ok = false
				}
			}
			ok = restartService(suiteConnectorService)
		}
		if ok {
			// Delete service configuration backup file
			ok = deleteFile(configFileBackup)
		}
	}
	// Delete devices and things
	fmt.Printf("performing cleanup on device id: %s\n", deviceID)
	if err := util.DeleteResources(&cfg, resources, deviceID, getTenantURL(),
		c2eCfg.DeviceRegistryAPIUsername, c2eCfg.DeviceRegistryAPIPassword); err != nil {
		fmt.Printf(deleteResourcesTemplate, indent, err)
		ok = false
	}
	return ok
}

func deleteFile(path string) bool {
	if err := os.Remove(path); err != nil {
		fmt.Printf("unable to delete file %s, error: %v", path, err)
		return false
	}
	return true
}

func getTenantURL() string {
	return fmt.Sprintf(
		"%s/v1/devices/%s/", strings.TrimSuffix(c2eCfg.DeviceRegistryAPIAddress, "/"), tenantID)
}

func writeConfigFile(path string) error {
	cfg := &util.ConnectorConfiguration{
		CaCert:   caCert,
		LogFile:  logFile,
		Address:  c2eCfg.MQTTAdapterAddress,
		DeviceID: deviceID,
		TenantID: tenantID,
		AuthID:   authID,
		Password: password,
	}
	return util.WriteConfigFile(path, cfg)
}

func writeConfigBootstrapFile(path string) error {
	postScriptArr := []string{postBootstrapScript}
	cfg := &util.BootstrapConfiguration{
		LogFile:             logBootstrapFile,
		PostBootstrapFile:   postBootstrapFile,
		PostBootstrapScript: postScriptArr,
		Address:             c2eCfg.MQTTAdapterAddress,
		TenantID:            tenantID,
		DeviceID:            deviceID,
		AuthID:              authID,
		Password:            password,
	}
	return util.WriteConfigFile(path, cfg)
}

func writeConfigLdtFile(path string) error {
	type ldtConnectorConfig struct {
		CaCert   string `json:"caCert"`
		LogFile  string `json:"logFile"`
		Address  string `json:"address"`
		TenantID string `json:"tenantId"`
		DeviceID string `json:"deviceId"`
		AuthID   string `json:"authId"`
		Password string `json:"password"`
		ThingsDb string `json:"thingsDb"`
	}

	cfg := &ldtConnectorConfig{
		CaCert:   ldtCaCert,
		LogFile:  logLdtFile,
		Address:  c2eCfg.MQTTAdapterAddress,
		DeviceID: deviceID,
		TenantID: tenantID,
		AuthID:   authID,
		Password: password,
		ThingsDb: thingsDb,
	}
	return util.WriteConfigFile(path, cfg)
}

func restartService(service string) bool {
	fmt.Printf("restarting %s...", service)
	stdout, err := exec.Command(systemctl, restart, service).Output()
	if stdout != nil {
		fmt.Println(string(stdout))
	}
	if err != nil {
		fmt.Printf("error restarting %s: %v", service, err)
		return false
	}
	fmt.Println("... done")
	return true
}
