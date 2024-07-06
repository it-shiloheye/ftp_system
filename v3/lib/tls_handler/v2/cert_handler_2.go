package tlshandler

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	// "log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

func main() {
	// get our ca and server certificate
	serverTLSConf, clientTLSConf, err := certsetup()
	if err != nil {
		panic(err)
	}

	// set up the httptest.Server using our certificate signed by our CA
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "success!")
	}))
	server.TLS = serverTLSConf
	server.StartTLS()
	defer server.Close()

	// communicate with the server using an http.Client configured to trust our CA
	transport := &http.Transport{
		TLSClientConfig: clientTLSConf,
	}
	http := http.Client{
		Transport: transport,
	}
	resp, err := http.Get(server.URL)
	if err != nil {
		panic(err)
	}

	// verify the response
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	body := strings.TrimSpace(string(respBodyBytes[:]))
	if body == "success!" {
		fmt.Println(body)
	} else {
		panic("not successful!")
	}
}

func certsetup() (serverTLSConf *tls.Config, clientTLSConf *tls.Config, err error) {
	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return nil, nil, err
	}

	serverTLSConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPEM.Bytes())
	clientTLSConf = &tls.Config{
		RootCAs: certpool,
	}

	return
}

type CAPem struct {
	X509       x509.Certificate `json:"x509"`
	PrivKey    *rsa.PrivateKey  `json:"priv_key"`
	Cert       []byte           `json:"ca"`
	PEM        []byte           `json:"pem"`
	PrivKeyPEM []byte           `json:"priv_key_pem"`
}

func GenerateCAPem(X509_ x509.Certificate) (caPem CAPem, err error) {
	loc := log_item.Loc("func GenerateCAPem() (caPem CAPem, err error)")
	caPem = CAPem{
		X509: X509_,
	}

	// create our private and public key
	caPem.PrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("caPem.PrivKey, err = rsa.GenerateKey(rand.Reader, 4096)").
			SetMessage(err.Error()).
			AppendParentError(err)
		return
	}

	// create the CA
	caPem.Cert, err = x509.CreateCertificate(rand.Reader, &caPem.X509, &caPem.X509, &caPem.PrivKey.PublicKey, caPem.PrivKey)
	if err != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("caPem.Cert, err = x509.CreateCertificate(rand.Reader, &caPem.X509, &caPem.X509, &caPem.PrivKey.PublicKey, caPem.PrivKey)").
			SetMessage(err.Error()).
			AppendParentError(err)
		return
	}
	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caPem.Cert,
	})
	caPem.PEM, _ = io.ReadAll(caPEM)

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPem.PrivKey),
	})
	caPem.PrivKeyPEM, _ = io.ReadAll(caPrivKeyPEM)

	return
}

type TLSCert struct {
	CAPem

	TlsCert tls.Certificate `json:"tls_cert"`
}

func GenerateTLSCert(ca_pem *CAPem, X509_ x509.Certificate) (tls_cert *TLSCert, err error) {
	loc := log_item.Loc("func GenerateTLSCert(ca_pem CAPem) (tls_cert TLSCert, err error)")
	tls_cert = &TLSCert{}
	// set up our server certificate
	tls_cert.X509 = X509_

	tls_cert.PrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("tls_cert.PrivKey, err = rsa.GenerateKey(rand.Reader, 4096)").
			SetMessage(err.Error()).
			AppendParentError(err)
		return
	}

	// log.Println(tls_cert)
	tls_cert.Cert, err = x509.CreateCertificate(rand.Reader, &tls_cert.X509, &ca_pem.X509, &tls_cert.PrivKey.PublicKey, ca_pem.PrivKey)
	if err != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("tls_cert.Cert, err = x509.CreateCertificate(rand.Reader, &tls_cert.X509, &ca_pem.X509, &tls_cert.PrivKey.PublicKey, ca_pem.PrivKey)").
			SetMessage(err.Error()).
			AppendParentError(err)
		return
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: tls_cert.Cert,
	})
	tls_cert.PEM, _ = io.ReadAll(certPEM)

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(tls_cert.PrivKey),
	})
	tls_cert.PrivKeyPEM, _ = io.ReadAll(certPrivKeyPEM)

	tls_cert.TlsCert, err = tls.X509KeyPair(tls_cert.PEM, tls_cert.PrivKeyPEM)
	if err != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("tls_cert.TlsCert, err = tls.X509KeyPair(tls_cert.PEM, tls_cert.PrivKeyPEM)").
			SetMessage(err.Error()).
			AppendParentError(err)
		return
	}

	return
}
