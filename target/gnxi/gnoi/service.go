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
	"crypto/tls"
	"crypto/x509"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"ovs-gnxi/shared"
	pbc "ovs-gnxi/shared/gnoi/modeldata/generated/cert"
	pbs "ovs-gnxi/shared/gnoi/modeldata/generated/system"
)

type RebootCallback func() error
type RotateCertificatesCallback func(certificates *shared.ServerCertificates) error

type Service struct {
	auth *shared.Authenticator
	//certServer          *cert.Server
	//certManager         *cert.Manager
	defaultCertificate  *tls.Certificate
	callbackReboot      RebootCallback
	callbackRotateCerts RotateCertificatesCallback
}

func NewService(auth *shared.Authenticator, defaultCertificate *tls.Certificate, callbackReboot RebootCallback, callbackRotateCerts RotateCertificatesCallback) (*Service, error) {
	//certManager := cert.NewManager(defaultCertificate.PrivateKey)
	//certServer := cert.NewServer(certManager)
	return &Service{
		auth: auth,
		//certServer:          certServer,
		//certManager:         certManager,
		defaultCertificate:  defaultCertificate,
		callbackReboot:      callbackReboot,
		callbackRotateCerts: callbackRotateCerts,
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
	//s.certServer.Register(g)
	pbs.RegisterSystemServer(g, s)
	//pbc.RegisterCertificateManagementServer(g, s)
	reflection.Register(g)
}

func (s *Service) Reboot(context.Context, *pbs.RebootRequest) (*pbs.RebootResponse, error) {
	if err := s.callbackReboot(); err != nil {
		return nil, err
	}

	return nil, status.Error(codes.Unimplemented, "Reboot is not implemented.")
}

func (s *Service) RebootStatus(context.Context, *pbs.RebootStatusRequest) (*pbs.RebootStatusResponse, error) {
	return nil, status.Error(codes.Unimplemented, "RebootStatus is not implemented.")
}

func (s *Service) CancelReboot(context.Context, *pbs.CancelRebootRequest) (*pbs.CancelRebootResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CancelReboot is not implemented.")
}

func (s *Service) Ping(*pbs.PingRequest, pbs.System_PingServer) error {
	return status.Error(codes.Unimplemented, "Ping is not implemented.")
}

func (s *Service) Traceroute(*pbs.TracerouteRequest, pbs.System_TracerouteServer) error {
	return status.Error(codes.Unimplemented, "Traceroute is not implemented.")
}

func (s *Service) Time(context.Context, *pbs.TimeRequest) (*pbs.TimeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Time is not implemented.")
}

func (s *Service) SetPackage(pbs.System_SetPackageServer) error {
	return status.Error(codes.Unimplemented, "SetPackage is not implemented.")
}

func (s *Service) SwitchControlProcessor(context.Context, *pbs.SwitchControlProcessorRequest) (*pbs.SwitchControlProcessorResponse, error) {
	return nil, status.Error(codes.Unimplemented, "SwitchControlProcessor is not implemented.")
}

func (s *Service) Rotate(pbc.CertificateManagement_RotateServer) error {
	return status.Error(codes.Unimplemented, "Rotate is not implemented.")
}

func (s *Service) GetCertificates(context.Context, *pbc.GetCertificatesRequest) (*pbc.GetCertificatesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetCertificates is not implemented.")
}

func (s *Service) Install(pbc.CertificateManagement_InstallServer) error {
	return status.Error(codes.Unimplemented, "Install is not implemented.")
}

func (s *Service) RevokeCertificates(context.Context, *pbc.RevokeCertificatesRequest) (*pbc.RevokeCertificatesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "RevokeCertificates is not implemented.")
}

func (s *Service) CanGenerateCSR(context.Context, *pbc.CanGenerateCSRRequest) (*pbc.CanGenerateCSRResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CanGenerateCSR is not implemented.")
}

/*
func (s *Service) RegisterNotifier(f cert.Notifier) {
	s.certManager.RegisterNotifier(f)
}*/
