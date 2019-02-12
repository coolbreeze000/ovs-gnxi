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

// Package gnxi implements a gnxi server.
package gnxi

import (
	"ovs-gnxi/shared"
	"ovs-gnxi/shared/gnmi/modeldata"
	"ovs-gnxi/shared/gnmi/modeldata/generated/ocstruct"
	"ovs-gnxi/shared/logging"
	"ovs-gnxi/target/gnxi/service"
	"ovs-gnxi/target/gnxi/service/gnmi"
	"ovs-gnxi/target/ovs"
	"reflect"
)

const (
	caPATH        = "certs/ca.crt"
	certPATH      = "certs/target.crt"
	keyPATH       = "certs/target.key"
	adminUsername = "admin"
	adminPassword = "testpassword"
	defaultCertID = "c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667"
)

var log = logging.New("ovs-gnxi")

// Server struct maintains the data structure for device config and implements the interface of gnxi server. It supports Capabilities, Get, and Set APIs.
type Server struct {
	Auth              *shared.Authenticator
	Certs             *shared.TargetCertificates
	SystemBroker      *ovs.SystemBroker
	Service           *service.Service
	certificateChange chan struct{}
}

// NewServer creates an instance of Server.
func NewServer() (*Server, error) {
	log.Info("Initializing gNXI Server...")

	auth := shared.NewAuthenticator(adminUsername, adminPassword)

	certs, err := shared.NewServerCertificates(defaultCertID, caPATH, certPATH, keyPATH)
	if err != nil {
		return nil, err
	}

	s := &Server{Auth: auth, Certs: certs}
	s.SystemBroker = ovs.NewSystemBroker(s.Service, s.Certs)

	return s, nil
}

func (s *Server) InitializeService() {
	s.Service = s.createService()
	s.SystemBroker.GNXIService = s.Service
}

func (s *Server) createService() *service.Service {
	<-s.SystemBroker.OVSClient.Config.Initialized

	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*ocstruct.Device)(nil)),
		ocstruct.SchemaTree["Device"],
		ocstruct.Unmarshal,
		ocstruct.ΛEnum)

	log.Info("Start generating initial gNMI config from OVS system source...")

	config, err := s.SystemBroker.GenerateConfig(s.SystemBroker.OVSClient.Config)
	if err != nil {
		log.Fatalf("Unable to generate gNMI Config: %v", err)
	}

	log.Debugf("Using following initial config data: %s", config)

	s.SystemBroker.OVSClient.Config.OverwriteCallback(s.SystemBroker.OVSConfigChangeCallback)
	c, err := service.NewService(s.Auth, model, s.Certs, []byte(config), s.SystemBroker.GNMIConfigSetupCallback, s.SystemBroker.GNMIConfigChangeCallback, s.SystemBroker.GNOIRebootCallback, s.SystemBroker.GNOIRotateCertificatesCallback)
	if err != nil {
		log.Fatalf("Error on creating gNMI service: %v", err)
	}

	return c
}
