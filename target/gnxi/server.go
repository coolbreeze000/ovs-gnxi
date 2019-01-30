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
	"crypto/tls"
	"fmt"
	"net"
	"ovs-gnxi/shared"
	"ovs-gnxi/shared/logging"
	"ovs-gnxi/target/gnxi/gnmi"
	"ovs-gnxi/target/gnxi/gnmi/modeldata/generated/ocstruct"
	"ovs-gnxi/target/ovs"
	"sync"

	"ovs-gnxi/target/gnxi/gnmi/modeldata"
	"ovs-gnxi/target/gnxi/gnoi"
	"reflect"
)

const (
	caPATH        = "certs/ca.crt"
	certPATH      = "certs/target.crt"
	keyPATH       = "certs/target.key"
	adminUsername = "admin"
	adminPassword = "testpassword"
	gnxiProtocol  = "tcp"
	portGNMI      = "10161"
	portGNOI      = "10162"
)

var log = logging.New("ovs-gnxi")

// Server struct maintains the data structure for device config and implements the interface of gnxi server. It supports Capabilities, Get, and Set APIs.
type Server struct {
	Auth              *shared.Authenticator
	certs             *shared.ServerCertificates
	SystemBroker      *ovs.SystemBroker
	serviceGNMI       *gnmi.Service
	serviceGNOI       *gnoi.Service
	certificateChange chan struct{}
}

// NewServer creates an instance of Server.
func NewServer() (*Server, error) {
	log.Info("Initializing gNXI Server...")

	auth := shared.NewAuthenticator(adminUsername, adminPassword)

	certs, err := shared.NewServerCertificates(caPATH, certPATH, keyPATH)
	if err != nil {
		return nil, err
	}

	log.Info("TEST3")

	s := &Server{Auth: auth, certs: certs}
	s.SystemBroker = ovs.NewSystemBroker(s.serviceGNMI, s.serviceGNOI, s.certs)

	log.Info("TEST4")

	return s, nil
}

func (s *Server) InitializeServices() {
	s.serviceGNMI = s.createGNMIService()
	s.serviceGNOI = s.CreateGNOIService()
	s.SystemBroker.ServiceGNMI = s.serviceGNMI
	s.SystemBroker.ServiceGNOI = s.serviceGNOI
}

func (s *Server) createGNMIService() *gnmi.Service {
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
	c, err := gnmi.NewService(s.Auth, model, []byte(config), s.SystemBroker.GNMIConfigSetupCallback, s.SystemBroker.GNMIConfigChangeCallback)
	if err != nil {
		log.Fatalf("Error on creating gNMI service: %v", err)
	}

	return c
}

func (s *Server) CreateGNOIService() *gnoi.Service {
	c, err := gnoi.NewService(s.Auth, &s.certs.Certificate, s.SystemBroker.GNOICertificateChangeCallback)
	if err != nil {
		log.Fatalf("Error on creating gNOI service: %v", err)
	}

	return c
}

func (s *Server) RunGNMIService(wg *sync.WaitGroup) {
	defer wg.Done()

	g, err := s.serviceGNMI.PrepareServer([]tls.Certificate{s.certs.Certificate}, s.certs.CertPool)
	if err != nil {
		log.Fatalf("Failed to prepare gNMI server: %v", err)
	}

	s.serviceGNMI.RegisterService(g)

	log.Infof("Starting to listen")
	listen, err := net.Listen(gnxiProtocol, fmt.Sprintf(":%s", portGNMI))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Info("Starting to serve gNMI")
	if err := g.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Error("GNMI Server exit")
}

func (s *Server) RunGNOIService(wg *sync.WaitGroup) {
	defer wg.Done()

	g, err := s.serviceGNOI.PrepareServer([]tls.Certificate{s.certs.Certificate}, s.certs.CertPool)
	if err != nil {
		log.Fatalf("Failed to prepare gNMI server: %v", err)
	}

	s.serviceGNOI.RegisterService(g)

	log.Infof("Starting to listen")
	listen, err := net.Listen(gnxiProtocol, fmt.Sprintf(":%s", portGNOI))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Info("Starting to serve gNOI")
	if err := g.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Error("GNMI Server exit")
}

// TODO(dherkel@google.com): Implement Certificate Rotation Channel Receiver
