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

const (
	deviceJSON = `{"authorities":["auto-provisioning-enabled"]}`

	authJSON = `[{
		"type": "hashed-password",
		"auth-id": "%s",
		"secrets": [{
			"pwd-plain": "%s"
		}]
	}]`

	policyJSON = `{
		"entries":{
		"DEFAULT":{
			"subjects":{
				"{{ request:subjectId }}":{
					"type":"Ditto user authenticated via nginx"
				}
			},
			"resources":{
				"thing:/":{
					"grant":[
					"READ",
					"WRITE"
					],
					"revoke":[
					
					]
				},
				"policy:/":{
					"grant":[
					"READ",
					"WRITE"
					],
					"revoke":[
					
					]
				},
				"message:/":{
					"grant":[
					"READ",
					"WRITE"
					],
					"revoke":[
					
					]
				}
			}
		},
		"HONO":{
			"subjects":{
				"pre-authenticated:%s":{
					"type":"Connection to Eclipse Hono"
				}
			},
			"resources":{
				"thing:/":{
					"grant":[
					"READ",
					"WRITE"
					],
					"revoke":[
					
					]
				},
				"message:/":{
					"grant":[
					"READ",
					"WRITE"
					],
					"revoke":[
					
					]
				}
			}
		}
		}
	}`

	thingJSON = `{"policyId": "%s"}`
)
