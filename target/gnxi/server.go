/* Copyright 2019 Google Inc.

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

// Contains code from source: https://github.com/google/gnxi/tree/8521faedac371b6e13ac76f752eb41079ca79bd7

package gnxi

import (
	"ovs-gnxi/shared"
	"ovs-gnxi/shared/gnmi/modeldata"
	"ovs-gnxi/shared/gnmi/modeldata/generated/ocstruct"
	"ovs-gnxi/shared/logging"
	"ovs-gnxi/target/cert"
	"ovs-gnxi/target/gnxi/service"
	"ovs-gnxi/target/gnxi/service/gnmi"
	"ovs-gnxi/target/ovs"
	"reflect"
)

const (
	certRootSystemPath = "certs"
	adminUsername      = "admin"
	adminPassword      = "testpassword"
)

var log = logging.New("ovs-gnxi")

type Server struct {
	Auth              *shared.Authenticator
	CertManager       *cert.Manager
	SystemBroker      *ovs.SystemBroker
	Service           *service.Service
	certificateChange chan struct{}
}

// NewServer creates an instance of Server.
func NewServer() (*Server, error) {
	log.Info("Initializing gNXI Server...")

	auth := shared.NewAuthenticator(adminUsername, adminPassword)

	certManager, err := cert.NewCertManager(certRootSystemPath)
	if err != nil {
		return nil, err
	}

	s := &Server{Auth: auth, CertManager: certManager}
	s.SystemBroker = ovs.NewSystemBroker(s.Service, s.CertManager)

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
	c, err := service.NewService(s.Auth, model, s.CertManager, []byte(config), s.SystemBroker.GNMIConfigSetupCallback, s.SystemBroker.GNMIConfigChangeCallback, s.SystemBroker.GNOIRebootCallback, s.SystemBroker.GNOIRotateCertificatesCallback)
	if err != nil {
		log.Fatalf("Error on creating gNMI service: %v", err)
	}

	return c
}
