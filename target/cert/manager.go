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

package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	pbc "ovs-gnxi/shared/gnoi/modeldata/generated/cert"
	"ovs-gnxi/shared/logging"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var log = logging.New("ovs-gnxi")

const (
	defaultKeysSize = 4096
	defaultCertID   = "c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667"
	activePath      = "active"
	certFileName    = "target.crt"
	keyFileName     = "target.key"
	caCertFileName  = "ca.crt"
)

type Package struct {
	Finalized             bool
	CertificateID         string
	certificateSystemPath string
	keySystemPath         string
	caSystemPath          string
	Certificate           *x509.Certificate
	PublicKey             *rsa.PublicKey
	PrivateKey            *rsa.PrivateKey
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

func (p *Package) ReadPEMToX509Cert(data []byte) error {
	certBlock, _ := pem.Decode(data)
	if certBlock == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}

	certs, err := x509.ParseCertificates(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err)
	}

	p.Certificate = certs[0]
	p.PublicKey = p.Certificate.PublicKey.(*rsa.PublicKey)
	p.CertInfo = []*pbc.CertificateInfo{
		{
			CertificateId: p.CertificateID,
			Certificate: &pbc.Certificate{
				Type:        pbc.CertificateType_CT_X509,
				Certificate: x509toPEM(p.Certificate),
			},
			ModificationTime: time.Now().UnixNano(),
		},
	}

	return nil
}

func (p *Package) ReadPEMToX509CACerts(rawCerts []*pbc.Certificate) error {
	var pemCACerts [][]byte

	for _, cert := range rawCerts {
		if cert.Type != pbc.CertificateType_CT_X509 {
			return fmt.Errorf("unexpected Certificate type: %d", cert.Type)
		}
		pemCACerts = append(pemCACerts, cert.Certificate)
	}

	rawCert := concatAppend(pemCACerts)

	caCertsBlock, _ := pem.Decode(rawCert)
	if caCertsBlock == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}

	caCerts, err := x509.ParseCertificates(caCertsBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err)
	}

	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(rawCert); !ok {
		return fmt.Errorf("failed to append CA certificate")
	}

	p.CACertificates = caCerts
	p.CertPool = certPool

	return nil
}

type Manager struct {
	active         *Package
	collection     map[string]*Package
	rootSystemPath string
	mu             sync.RWMutex
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
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.collection[certID]; !ok {
		return fmt.Errorf("unable to activate non existing cert package")
	}

	if m.collection[certID].Finalized != true {
		return fmt.Errorf("unable to activate non finalized cert package")
	}

	err := m.copyCertPackageToActivePath(certID)
	if err != nil {
		return err
	}

	m.active = m.collection[certID]

	log.Infof("Cert package %v is now active", m.active.CertificateID)

	return nil
}

func (m *Manager) GetActivePackageCertPath() string {
	return path.Join(path.Join(m.rootSystemPath, activePath), filepath.Base(m.active.certificateSystemPath))
}

func (m *Manager) GetActivePackageKeyPath() string {
	return path.Join(path.Join(m.rootSystemPath, activePath), filepath.Base(m.active.keySystemPath))
}

func (m *Manager) GetActivePackageCAPath() string {
	return path.Join(path.Join(m.rootSystemPath, activePath), filepath.Base(m.active.caSystemPath))
}

func (m *Manager) GetActivePackage() *Package {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.active
}

func (m *Manager) loadCertFromPath(path string) (*x509.Certificate, error) {
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

func (m *Manager) loadKeyFromPath(path string) (*rsa.PrivateKey, error) {
	keyFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read key: %v", err)
	}

	keyBlock, _ := pem.Decode(keyFile)
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to parse key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %v", err)
	}

	return key.(*rsa.PrivateKey), nil
}

func (m *Manager) loadCACertsFromPath(path string) ([]*x509.Certificate, *x509.CertPool, error) {
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

	cert, err := m.loadCertFromPath(certPath)
	if err != nil {
		return err
	}

	key, err := m.loadKeyFromPath(keyPath)
	if err != nil {
		return err
	}

	ca, caCertPool, err := m.loadCACertsFromPath(caPath)
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
		certificateSystemPath: certPath,
		keySystemPath:         keyPath,
		caSystemPath:          caPath,
		Certificate:           cert,
		PublicKey:             cert.PublicKey.(*rsa.PublicKey),
		PrivateKey:            key,
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

	certFile, err := os.Create(m.collection[certID].certificateSystemPath)
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

	keyFile, err := os.OpenFile(m.collection[certID].keySystemPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return fmt.Errorf("unable to create key file: %v", err)
	}

	keyBlock, err := x509.MarshalPKCS8PrivateKey(m.collection[certID].PrivateKey)

	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBlock})
	if err != nil {
		return fmt.Errorf("unable to write key to file: %v", err)
	}

	err = keyFile.Close()
	if err != nil {
		return fmt.Errorf("unable to close key file: %v", err)
	}

	caFile, err := os.Create(m.collection[certID].caSystemPath)
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

func (m *Manager) copyCertPackageToActivePath(certID string) error {
	basePath := path.Join(m.rootSystemPath, activePath)
	if _, err := os.Stat(activePath); os.IsNotExist(err) {
		os.Mkdir(activePath, 0644)
	}

	err := copySrcFileToDstPath(m.collection[certID].certificateSystemPath, path.Join(basePath, filepath.Base(m.collection[certID].certificateSystemPath)))
	if err != nil {
		return fmt.Errorf("unable to copy cert file: %v", err)
	}

	err = copySrcFileToDstPath(m.collection[certID].keySystemPath, path.Join(basePath, filepath.Base(m.collection[certID].keySystemPath)))
	if err != nil {
		return fmt.Errorf("unable to copy key file: %v", err)
	}

	err = copySrcFileToDstPath(m.collection[certID].caSystemPath, path.Join(basePath, filepath.Base(m.collection[certID].caSystemPath)))
	if err != nil {
		return fmt.Errorf("unable to copy ca file: %v", err)
	}

	return nil
}

func copySrcFileToDstPath(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func x509toPEM(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
}

func concatAppend(data [][]byte) []byte {
	var r []byte
	for _, s := range data {
		r = append(r, s...)
	}
	return r
}

func (m *Manager) InitializePackage() (*Package, error) {
	key, err := m.generatePrivateKey(defaultKeysSize)
	if err != nil {
		return nil, err
	}

	return &Package{
		Finalized:  false,
		PrivateKey: key,
	}, nil
}

func (m *Manager) FinalizePackage(p *Package) error {
	basePath := path.Join(m.rootSystemPath, p.CertificateID)
	p.certificateSystemPath = path.Join(basePath, filepath.Base(m.active.certificateSystemPath))
	p.keySystemPath = path.Join(basePath, filepath.Base(m.active.keySystemPath))
	p.caSystemPath = path.Join(basePath, filepath.Base(m.active.caSystemPath))

	m.collection[p.CertificateID] = p

	err := m.ExportCertPackageToPath(m.collection[p.CertificateID].CertificateID)
	if err != nil {
		return err
	}

	p.TLSCertKeyPair, err = m.loadTLSKeyPairFromPath(m.collection[p.CertificateID].CertificateID)
	if err != nil {
		return err
	}

	m.collection[p.CertificateID].Finalized = true

	return nil
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
