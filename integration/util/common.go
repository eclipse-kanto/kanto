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
	"encoding/json"
	"fmt"
	"os"
)

const (
	configDefaultMode = 0666
)

// Convert marshals an object(e.g. map) to a JSON payload and unmarshals it to the given structure
func Convert(from interface{}, to interface{}) error {
	jsonValue, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonValue, to)
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
