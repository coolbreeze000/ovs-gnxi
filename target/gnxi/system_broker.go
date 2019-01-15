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

package gnxi

import (
	"github.com/openconfig/ygot/ygot"
	"os"
	"ovs-gnxi/target/gnxi/gnmi"
	"ovs-gnxi/target/ovs"
	"sync"
)

const (
	ovsAddress  = "ovs.gnxi.lan"
	ovsProtocol = "tcp"
	ovsPort     = "6640"
)

type SystemBroker struct {
	GNXIServer *Server
	OVSClient  *ovs.Client
}

func NewSystemBroker(gNXIServer *Server) *SystemBroker {
	var err error
	s := &SystemBroker{GNXIServer: gNXIServer}

	log.Info("Initializing OVS Client...")

	s.OVSClient, err = ovs.NewClient(ovsAddress, ovsProtocol, ovsPort, s.GNXIServer.Certs.KeySystemPath, s.GNXIServer.Certs.CertSystemPath, s.GNXIServer.Certs.CASystemPath, ovs.NewConfig(nil))
	if err != nil {
		log.Errorf("Unable to initialize OVS Client: %v", err)
		os.Exit(1)
	}

	return s
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

func (s *SystemBroker) GNOICertificateChangeCallback(certs *ServerCertificates) error {
	return nil
}

func (s *SystemBroker) RunOVSClient(wg *sync.WaitGroup) {
	defer s.OVSClient.Connection.Disconnect()
	defer wg.Done()
	s.OVSClient.StartMonitorAll()
	log.Error("OVS Client exit")
}
