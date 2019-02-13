package shared

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	pbc "ovs-gnxi/shared/gnoi/modeldata/generated/cert"
	"ovs-gnxi/shared/logging"
	"time"
)

var log = logging.New("ovs-gnxi")

const (
	KeySize = 4096
)

type TargetCertificates struct {
	CertificateID   string
	CASystemPath    string
	CertSystemPath  string
	KeySystemPath   string
	TLSCertificates []tls.Certificate
	CACertificates  []*x509.Certificate
	CertPool        *x509.CertPool
	CertInfo        []*pbc.CertificateInfo
}

func LoadCertificatesFromFile(path string) ([]*x509.Certificate, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate: %v", err)
	}

	block, _ := pem.Decode(file)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	certificates, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return certificates, nil
}

// LoadCertificates loads certificates from file.
func NewServerCertificates(id, ca, cert, key string) (*TargetCertificates, error) {
	tlsCertificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("could not load target key pair: %v", err)
	}

	certPool := x509.NewCertPool()

	certificates, err := LoadCertificatesFromFile(cert)
	if err != nil {
		return nil, err
	}

	caCertificates, err := LoadCertificatesFromFile(ca)
	if err != nil {
		return nil, err
	}

	caFile, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(caFile); !ok {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	return &TargetCertificates{
			CertificateID:   id,
			CASystemPath:    ca,
			CertSystemPath:  cert,
			KeySystemPath:   key,
			TLSCertificates: []tls.Certificate{tlsCertificate},
			CACertificates:  caCertificates,
			CertPool:        certPool,
			CertInfo: []*pbc.CertificateInfo{
				{
					CertificateId: id,
					Certificate: &pbc.Certificate{
						Type:        pbc.CertificateType_CT_X509,
						Certificate: x509toPEM(certificates[0]),
					},
					ModificationTime: time.Now().UnixNano(),
				},
			},
		},
		nil
}

func x509toPEM(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
}

func CreateCSR(rand io.Reader, template *x509.CertificateRequest, priv interface{}) ([]byte, error) {
	der, err := x509.CreateCertificateRequest(rand, template, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSR: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: der,
	}), nil
}

func GeneratePrivateKey(size int) (*rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key")
	}
	if bits := priv.N.BitLen(); bits != size {
		return nil, fmt.Errorf("key too short (%d vs %d)", bits, size)
	}

	return priv, nil
}
