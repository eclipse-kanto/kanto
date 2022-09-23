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
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func initConfigFromEnv(cfg interface{}, prefix string) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		configValue, ok := f.Tag.Lookup("def")

		envName := toSnakeCase(f.Name, prefix)
		if env, ok := os.LookupEnv(envName); ok {
			configValue = env
		}

		if !ok && configValue == "" {
			return fmt.Errorf("env variable %s not set", envName)
		}

		fv := v.Field(i)
		switch f.Type.Kind() {
		case reflect.String:
			fv.SetString(configValue)
		case reflect.Int:
			if i, err := strconv.Atoi(configValue); err == nil {
				fv.SetInt(int64(i))
			} else {
				return fmt.Errorf("invalid value %s for config property %s: %w", configValue, envName, err)
			}
		default:
			panic(fmt.Errorf("unsupported config value type: %v for field %s", f.Type, f.Name))
		}
	}

	return nil
}

func getConfigHelp(cfg interface{}, prefix string) string {
	result := strings.Builder{}

	result.WriteString("config environmental variables:")
	t := reflect.TypeOf(cfg)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		name := toSnakeCase(f.Name, prefix)

		result.WriteString("\n\t - ")
		result.WriteString(name)

		def, ok := f.Tag.Lookup("def")
		if ok {
			result.WriteString(fmt.Sprintf(" (default '%s')", def))
		}
	}

	return result.String()
}

func toSnakeCase(name string, prefix string) string {
	var result, word strings.Builder

	result.WriteString(prefix)

	for i, ch := range name {
		if i > 0 && unicode.IsUpper(ch) {
			if result.Len() > 0 {
				result.WriteByte('_')
			}

			result.WriteString(strings.ToUpper(word.String()))
			word.Reset()
		}

		word.WriteRune(ch)
	}

	if result.Len() > 0 {
		result.WriteByte('_')
	}
	result.WriteString(strings.ToUpper(word.String()))

	return result.String()
}
