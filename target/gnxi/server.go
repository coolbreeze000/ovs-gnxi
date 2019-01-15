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
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"ovs-gnxi/shared/gnxi"
	"ovs-gnxi/shared/logging"
	"ovs-gnxi/target/gnxi/gnmi"
	"ovs-gnxi/target/gnxi/gnmi/modeldata/generated/ocstruct"
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
	Auth              *gnxi.Authenticator
	Certs             *ServerCertificates
	SystemBroker      *SystemBroker
	ServiceGNMI       *gnmi.Service
	ServiceGNOI       *gnoi.Service
	certificateChange chan struct{}
}

// NewServer creates an instance of Server.
func NewServer() (*Server, error) {
	log.Info("Initializing gNXI Server...")

	auth := &gnxi.Authenticator{}
	adminUser := gnxi.NewUser(adminUsername, adminPassword)
	auth.AddUser(adminUser)

	certs, err := NewServerCertificates(caPATH, certPATH, keyPATH)
	if err != nil {
		return nil, err
	}

	s := &Server{Auth: auth, Certs: certs}
	s.SystemBroker = NewSystemBroker(s)

	return s, nil
}

func (s *Server) InitializeServices() {
	s.ServiceGNMI = s.createGNMIService()
	s.ServiceGNOI = s.CreateGNOIService()
}

type ServerCertificates struct {
	CASystemPath   string
	CertSystemPath string
	KeySystemPath  string
	Certificate    tls.Certificate
	CACertificates []*x509.Certificate
	CertPool       *x509.CertPool
}

// LoadCertificates loads certificates from file.
func NewServerCertificates(ca, cert, key string) (*ServerCertificates, error) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not load target key pair: %s", err))
	}

	certPool := x509.NewCertPool()

	caFile, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not read CA certificate: %s", err))
	}

	caCertificate, err := x509.ParseCertificates([]byte(caFile))
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	if ok := certPool.AppendCertsFromPEM(caFile); !ok {
		return nil, errors.New("failed to append CA certificate")
	}

	return &ServerCertificates{
			CASystemPath:   ca,
			CertSystemPath: cert,
			KeySystemPath:  key,
			Certificate:    certificate,
			CACertificates: caCertificate,
			CertPool:       certPool},
		nil
}

func (s *Server) createGNMIService() *gnmi.Service {
	<-s.SystemBroker.OVSClient.Config.Initialized

	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*ocstruct.Device)(nil)),
		ocstruct.SchemaTree["Device"],
		ocstruct.Unmarshal,
		ocstruct.ΛEnum)

	config, err := gnmi.GenerateConfig(s.SystemBroker.OVSClient.Config)
	if err != nil {
		log.Fatalf("Unable to generate gNMI Config: %v", err)
	}

	log.Info(fmt.Sprintf("%s", config))

	s.SystemBroker.OVSClient.Config.OverwriteCallback(s.SystemBroker.OVSConfigChangeCallback)
	c, err := gnmi.NewService(s, model, []byte(config), s.SystemBroker.GNMIConfigChangeCallback)
	if err != nil {
		log.Fatalf("Error on creating gNMI service: %v", err)
	}

	return c
}

func (s *Server) CreateGNOIService() *gnoi.Service {
	c, err := gnoi.NewService(s, s.Certs.Certificate.PrivateKey, &s.Certs.Certificate, s.SystemBroker.GNOICertificateChangeCallback())
	if err != nil {
		log.Fatalf("Error on creating gNOI service: %v", err)
	}

	return c
}

func (s *Server) RunGNMIService(wg *sync.WaitGroup) {
	defer wg.Done()

	g, err := s.ServiceGNMI.PrepareServer([]tls.Certificate{s.Certs.Certificate}, s.Certs.CertPool)
	if err != nil {
		log.Fatalf("Failed to prepare gNMI server: %v", err)
	}

	s.ServiceGNMI.RegisterService(g)

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

	g, err := s.ServiceGNOI.PrepareServer([]tls.Certificate{s.Certs.Certificate}, s.Certs.CertPool)
	if err != nil {
		log.Fatalf("Failed to prepare gNMI server: %v", err)
	}

	s.ServiceGNOI.RegisterService(g)

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
