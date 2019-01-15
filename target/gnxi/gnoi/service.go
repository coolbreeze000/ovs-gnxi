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

// Package gnoi contains required services for running a gnoi server.
package gnoi

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"github.com/google/gnxi/gnoi/cert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"ovs-gnxi/target/gnxi"
)

const (
	RSABitSize = 4096
)

type CertificateCallback func(gnxi.ServerCertificates) error

type Service struct {
	server             *gnxi.Server
	certServer         *cert.Server
	certManager        *cert.Manager
	defaultCertificate *tls.Certificate
	callback           CertificateCallback
}

func NewService(server *gnxi.Server, privateKey crypto.PrivateKey, defaultCertificate *tls.Certificate) (*Service, error) {
	certManager := cert.NewManager(defaultCertificate.PrivateKey)
	certServer := cert.NewServer(certManager)
	return &Service{
		server:             server,
		certServer:         certServer,
		certManager:        certManager,
		defaultCertificate: defaultCertificate,
	}, nil
}

func (s *Service) PrepareServer(certificates []tls.Certificate, certPool *x509.CertPool) (*grpc.Server, error) {
	opts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: certificates,
		ClientCAs:    certPool,
	}))}

	return grpc.NewServer(opts...), nil
}

func (s *Service) RegisterService(g *grpc.Server) {
	s.certServer.Register(g)
}

// RegisterNotifier registers a function that will be called everytime the number
// of Certificates or CA Certificates changes.
func (s *Service) RegisterNotifier(f cert.Notifier) {
	s.certManager.RegisterNotifier(f)
}
