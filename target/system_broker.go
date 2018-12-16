/* Copyright 2018 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/openconfig/ygot/ygot"
	"ovs-gnxi/target/gnxi"
	"ovs-gnxi/target/gnxi/gnmi"
	"ovs-gnxi/target/ovs"
)

type SystemBroker struct {
	OVSClient  *ovs.Client
	GNXIServer *gnxi.Server
}

func NewSystemBroker() *SystemBroker {
	b := SystemBroker{}
	return &b
}

func (s *SystemBroker) OVSConfigChangeCallback(ovsConfig *ovs.Config) error {
	log.Debug("Received new change by OVS device")
	gnmiConfig, err := gnmi.GenerateConfig(ovsConfig)
	if err != nil {
		log.Errorf("Unable to generate gNMI config from OVS config source: %v", err)
		return err
	}
	s.GNXIServer.ServiceGNMI.OverwriteConfig(gnmiConfig)

	return nil
}

func (s *SystemBroker) GNMIConfigChangeCallback(config ygot.ValidatedGoStruct) error {
	return nil
}
