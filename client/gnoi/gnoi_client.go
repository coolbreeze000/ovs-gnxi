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

package gnoi

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/gnxi/utils/credentials"
	"github.com/google/gnxi/utils/entity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pbc "ovs-gnxi/shared/gnoi/modeldata/generated/cert"
	pbs "ovs-gnxi/shared/gnoi/modeldata/generated/system"
	"ovs-gnxi/shared/logging"
)

var (
	log = logging.New("ovs-gnxi-client")
)

type Client struct {
	targetAddress string
	targetName    string
	encodingName  string
	caEntity      *entity.Entity
}

// NewClient returns an instance of Client struct.
func NewClient(targetAddress, targetName, encodingName string, caEntity *entity.Entity) *Client {
	return &Client{
		targetAddress: targetAddress,
		targetName:    targetName,
		encodingName:  encodingName,
		caEntity:      caEntity,
	}
}

func (c *Client) Reboot(ctx context.Context, rebootMessage string) (*pbs.RebootResponse, error) {
	opts := credentials.ClientCredentials(c.targetName)
	conn, err := grpc.Dial(c.targetAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pbs.NewSystemClient(conn)

	request := &pbs.RebootRequest{
		Message: rebootMessage,
	}

	log.Debug("== Request:")
	log.Debug(proto.MarshalTextString(request))

	response, err := cli.Reboot(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error in rebooting system: %v", err)
	}

	return response, nil
}

// RotateCertificates rotates a certificate.
func (c *Client) RotateCertificates(ctx context.Context, certID string) error {
	opts := credentials.ClientCredentials(c.targetName)
	conn, err := grpc.Dial(c.targetAddress, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := pbc.NewCertificateManagementClient(conn)

	requestGen := &pbc.RotateCertificateRequest{
		RotateRequest: &pbc.RotateCertificateRequest_GenerateCsr{
			GenerateCsr: &pbc.GenerateCSRRequest{
				CsrParams: &pbc.CSRParams{
					Type:               pbc.CertificateType_CT_X509,
					MinKeySize:         4096,
					KeyType:            pbc.KeyType_KT_RSA,
					CommonName:         "target.gnxi.lan",
					Country:            "AT",
					State:              "Vienna",
					City:               "Vienna",
					Organization:       "Test",
					OrganizationalUnit: "Test",
				},
			},
		},
	}

	log.Debug("== Request:")
	log.Debug(proto.MarshalTextString(requestGen))

	certClient, err := cli.Rotate(ctx)
	if err != nil {
		return fmt.Errorf("failed stream: %v", err)
	}

	if err = certClient.Send(requestGen); err != nil {
		return fmt.Errorf("failed to send GenerateCSRRequest: %v", err)
	}

	resp, err := certClient.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive RotateCertificateResponse: %v", err)
	}

	genCSR := resp.GetGeneratedCsr()
	if genCSR == nil || genCSR.Csr == nil {
		return fmt.Errorf("expected GenerateCSRRequest, got something else")
	}

	derCSR, _ := pem.Decode(genCSR.Csr.Csr)
	if derCSR == nil {
		return fmt.Errorf("failed to decode CSR PEM block")
	}

	csr, err := x509.ParseCertificateRequest(derCSR.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CSR DER")
	}

	signedCert, err := c.sign(csr)
	if err != nil {
		return fmt.Errorf("failed to sign the CSR: %v", err)
	}

	certPEM := x509toPEM(signedCert)

	var caCertificates []*pbc.Certificate
	for _, caCert := range []*x509.Certificate{c.caEntity.Certificate.Leaf} {
		caCertificates = append(caCertificates, &pbc.Certificate{
			Type:        pbc.CertificateType_CT_X509,
			Certificate: x509toPEM(caCert),
		})
	}

	requestLoad := &pbc.RotateCertificateRequest{
		RotateRequest: &pbc.RotateCertificateRequest_LoadCertificate{
			LoadCertificate: &pbc.LoadCertificateRequest{
				Certificate: &pbc.Certificate{
					Type:        pbc.CertificateType_CT_X509,
					Certificate: certPEM,
				},
				KeyPair:        nil,
				CertificateId:  certID,
				CaCertificates: caCertificates,
			},
		},
	}

	log.Debug("== Request:")
	log.Debug(proto.MarshalTextString(requestLoad))

	if err = certClient.Send(requestLoad); err != nil {
		return fmt.Errorf("failed to send LoadCertificateRequest: %v", err)
	}

	if resp, err = certClient.Recv(); err != nil {
		return fmt.Errorf("failed to receive RotateCertificateResponse: %v", err)
	}
	loadCertificateResponse := resp.GetLoadCertificate()
	if loadCertificateResponse == nil {
		return fmt.Errorf("expected LoadCertificateResponse, got something else")
	}

	certs, err := c.GetCertificates(ctx)
	if err != nil {
		return err
	}

	if _, ok := certs[certID]; !ok {
		return fmt.Errorf("failed to validate rotated certificate: %v", err)
	}

	if err := certClient.Send(&pbc.RotateCertificateRequest{
		RotateRequest: &pbc.RotateCertificateRequest_FinalizeRotation{FinalizeRotation: &pbc.FinalizeRequest{}},
	}); err != nil {
		return fmt.Errorf("failed to send LoadCertificateRequest: %v", err)
	}

	return nil
}

// GetCertificates gets a map of certificates in the target, certID to certificate
func (c *Client) GetCertificates(ctx context.Context) (map[string]*x509.Certificate, error) {
	opts := credentials.ClientCredentials(c.targetName)
	conn, err := grpc.Dial(c.targetAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pbc.NewCertificateManagementClient(conn)

	resp, err := cli.GetCertificates(ctx, &pbc.GetCertificatesRequest{})
	if err != nil {
		return nil, err
	}

	certs := map[string]*x509.Certificate{}
	for _, certInfo := range resp.CertificateInfo {
		if certInfo.Certificate == nil {
			continue
		}
		x509Cert, err := PEMtox509(certInfo.Certificate.Certificate)
		if err != nil {
			return nil, fmt.Errorf("failed to decode certificate: %v", err)
		}
		certs[certInfo.CertificateId] = x509Cert
	}
	return certs, nil
}

// EncodeCert encodes a x509.Certificate into a PEM block.
func x509toPEM(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
}

// PEMtox509 decodes a PEM block into a x509.Certificate.
func PEMtox509(bytes []byte) (*x509.Certificate, error) {
	certDERBlock, _ := pem.Decode(bytes)
	if certDERBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	certificate, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DER bytes")
	}
	return certificate, nil
}

// sign is called to create a Certificate from a CSR.
func (c *Client) sign(csr *x509.CertificateRequest) (*x509.Certificate, error) {
	e, err := entity.FromSigningRequest(csr)
	if err != nil {
		return nil, fmt.Errorf("failed generating a cert from a CSR: %v", err)
	}
	if err := e.SignWith(c.caEntity); err != nil {
		return nil, fmt.Errorf("failed to sign the certificate: %v", err)
	}
	return e.Certificate.Leaf, nil
}
