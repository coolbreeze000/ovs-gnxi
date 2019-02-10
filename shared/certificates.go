package shared

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	pbc "ovs-gnxi/shared/gnoi/modeldata/generated/cert"
	"ovs-gnxi/shared/logging"
)

var log = logging.New("ovs-gnxi")

const defaultUUID string = "c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667"

type ServerCertificates struct {
	CertificateID  string
	CASystemPath   string
	CertSystemPath string
	KeySystemPath  string
	Certificate    tls.Certificate
	CACertificates []*x509.Certificate
	CertPool       *x509.CertPool
	CertInfo       []*pbc.CertificateInfo
}

// LoadCertificates loads certificates from file.
func NewServerCertificates(ca, cert, key string) (*ServerCertificates, error) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("could not load target key pair: %v", err)
	}

	certPool := x509.NewCertPool()

	caFile, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate: %v", err)
	}

	block, _ := pem.Decode(caFile)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	log.Info(string(caFile))

	caCertificate, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(caFile); !ok {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	return &ServerCertificates{
			CertificateID:  defaultUUID,
			CASystemPath:   ca,
			CertSystemPath: cert,
			KeySystemPath:  key,
			Certificate:    certificate,
			CACertificates: caCertificate,
			CertPool:       certPool,
			CertInfo:       []*pbc.CertificateInfo{},
		},
		nil
}
