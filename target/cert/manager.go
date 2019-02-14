package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	pbc "ovs-gnxi/shared/gnoi/modeldata/generated/cert"
	"ovs-gnxi/shared/logging"
	"path"
	"time"
)

var log = logging.New("ovs-gnxi")

const (
	defaultKeysSize = 4096
	defaultCertID   = "c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667"
	certFileName    = "target.crt"
	keyFileName     = "target.key"
	caCertFileName  = "ca.crt"
)

type Package struct {
	Finalized             bool
	CertificateID         string
	CertificateSystemPath string
	KeySystemPath         string
	CASystemPath          string
	Certificate           *x509.Certificate
	PublicKey             rsa.PublicKey
	PrivateKey            rsa.PrivateKey
	TLSCertKeyPair        []tls.Certificate
	CACertificates        []*x509.Certificate
	CertPool              *x509.CertPool
	CertInfo              []*pbc.CertificateInfo
}

func (p *Package) CreateCSR(country, organization, organizationalUnit, commonName string) ([]byte, error) {
	subject := pkix.Name{
		Country:            []string{country},
		Organization:       []string{organization},
		OrganizationalUnit: []string{organizationalUnit},
		CommonName:         commonName,
	}

	template := &x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	der, err := x509.CreateCertificateRequest(rand.Reader, template, p.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSR: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: der,
	}), nil
}

type Manager struct {
	active         *Package
	collection     map[string]*Package
	rootSystemPath string
}

func NewCertManager(rootSystemPath string) (*Manager, error) {
	m := &Manager{
		rootSystemPath: rootSystemPath,
		collection:     make(map[string]*Package),
	}

	err := m.ImportCertPackageFromPath(defaultCertID)
	if err != nil {
		return nil, err
	}

	err = m.ActivatePackage(defaultCertID)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) ActivatePackage(certID string) error {
	if _, ok := m.collection[certID]; !ok {
		return fmt.Errorf("unable to activate non existing cert package")
	}

	if m.collection[certID].Finalized != true {
		return fmt.Errorf("unable to activate non finalized cert package")
	}

	m.active = m.collection[certID]

	return nil
}

func (m *Manager) GetActivePackage() *Package {
	return m.active
}

func (m *Manager) loadCertFromPath(path, certID string) (*x509.Certificate, error) {
	certFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate: %v", err)
	}

	certBlock, _ := pem.Decode(certFile)
	if certBlock == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	certs, err := x509.ParseCertificates(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return certs[0], nil
}

func (m *Manager) loadKeyFromPath(path, certID string) (*rsa.PrivateKey, error) {
	keyFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read key: %v", err)
	}

	key, err := x509.ParsePKCS1PrivateKey(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %v", err)
	}

	return key, nil
}

func (m *Manager) loadCACertsFromPath(path, certID string) ([]*x509.Certificate, *x509.CertPool, error) {
	caCertsFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read certificate: %v", err)
	}

	caCertsBlock, _ := pem.Decode(caCertsFile)
	if caCertsBlock == nil {
		return nil, nil, fmt.Errorf("failed to parse certificate PEM")
	}

	caCerts, err := x509.ParseCertificates(caCertsBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(caCertsFile); !ok {
		return nil, nil, fmt.Errorf("failed to append CA certificate")
	}

	return caCerts, certPool, nil
}

func (m *Manager) loadTLSKeyPairFromPath(certID string) ([]tls.Certificate, error) {
	tlsPair, err := tls.LoadX509KeyPair(path.Join(m.rootSystemPath, certID, certFileName), path.Join(m.rootSystemPath, certID, keyFileName))
	if err != nil {
		return nil, fmt.Errorf("could not load key pair: %v", err)
	}

	return []tls.Certificate{tlsPair}, nil
}

func (m *Manager) ImportCertPackageFromPath(certID string) error {
	certPath := path.Join(m.rootSystemPath, certID, certFileName)
	keyPath := path.Join(m.rootSystemPath, certID, keyFileName)
	caPath := path.Join(m.rootSystemPath, certID, caCertFileName)

	cert, err := m.loadCertFromPath(certPath, certID)
	if err != nil {
		return err
	}

	key, err := m.loadKeyFromPath(keyPath, certID)
	if err != nil {
		return err
	}

	ca, caCertPool, err := m.loadCACertsFromPath(caPath, certID)
	if err != nil {
		return err
	}

	tlsPair, err := m.loadTLSKeyPairFromPath(certID)
	if err != nil {
		return err
	}

	m.collection[certID] = &Package{
		Finalized:             true,
		CertificateID:         certID,
		CertificateSystemPath: certPath,
		KeySystemPath:         keyPath,
		CASystemPath:          caPath,
		Certificate:           cert,
		PublicKey:             cert.PublicKey.(rsa.PublicKey),
		PrivateKey:            *key,
		TLSCertKeyPair:        tlsPair,
		CACertificates:        ca,
		CertPool:              caCertPool,
		CertInfo: []*pbc.CertificateInfo{
			{
				CertificateId: certID,
				Certificate: &pbc.Certificate{
					Type:        pbc.CertificateType_CT_X509,
					Certificate: x509toPEM(cert),
				},
				ModificationTime: time.Now().UnixNano(),
			},
		},
	}

	return nil
}

func (m *Manager) ExportCertPackageToPath(certID string) error {
	basePath := path.Join(m.rootSystemPath, certID)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.Mkdir(basePath, 0644)
	}

	certFile, err := os.Create(m.collection[certID].CertificateSystemPath)
	if err != nil {
		return fmt.Errorf("unable to create cert file: %v", err)
	}

	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: m.collection[certID].Certificate.Raw})
	if err != nil {
		return fmt.Errorf("unable to write cert to file: %v", err)
	}

	err = certFile.Close()
	if err != nil {
		return fmt.Errorf("unable to close cert file: %v", err)
	}

	keyFile, err := os.OpenFile(m.collection[certID].KeySystemPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return fmt.Errorf("unable to create key file: %v", err)
	}

	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(&m.collection[certID].PrivateKey)})
	if err != nil {
		return fmt.Errorf("unable to write key to file: %v", err)
	}

	err = keyFile.Close()
	if err != nil {
		return fmt.Errorf("unable to close key file: %v", err)
	}

	caFile, err := os.Create(m.collection[certID].CASystemPath)
	if err != nil {
		return fmt.Errorf("unable to create ca file: %v", err)
	}

	for _, c := range m.collection[certID].CACertificates {
		err = pem.Encode(caFile, &pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})
		if err != nil {
			return fmt.Errorf("unable to write ca to file: %v", err)
		}
	}

	err = caFile.Close()
	if err != nil {
		return fmt.Errorf("unable to close ca file: %v", err)
	}

	return nil
}

func x509toPEM(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
}

func (m *Manager) InitializePackage() (*Package, error) {
	key, err := m.generatePrivateKey(defaultKeysSize)
	if err != nil {
		return nil, err
	}

	return &Package{
		Finalized:  false,
		PrivateKey: *key,
	}, nil
}

func (m *Manager) generatePrivateKey(size int) (*rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key")
	}
	if bits := priv.N.BitLen(); bits != size {
		return nil, fmt.Errorf("key too short (%d vs %d)", bits, size)
	}

	return priv, nil
}
