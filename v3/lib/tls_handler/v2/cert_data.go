package tlshandler

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"
)

type CertData struct {
	Organization  string         `json:"organisation"`
	Country       string         `json:"country"`
	Province      string         `json:"province"`
	Locality      string         `json:"locality"`
	StreetAddress string         `json:"street_address"`
	PostalCode    string         `json:"postal_code"`
	NotAfter      NotAfterStruct `json:"add_date"`
	IPAddrresses  []net.IP       `json:"ip_addresses"`
}

type NotAfterStruct struct {
	Years  int `json:"years"`
	Months int `json:"months"`
	Days   int `json:"days"`
}

func GetSubject(cd CertData) pkix.Name {

	return pkix.Name{
		Organization:  []string{cd.Organization},
		Country:       []string{cd.Country},
		Province:      []string{cd.Province},
		Locality:      []string{cd.Locality},
		StreetAddress: []string{cd.StreetAddress},
		PostalCode:    []string{cd.PostalCode},
	}
}

func ServerTLSConf(serverCert tls.Certificate) (serverTLSConf *tls.Config) {
	serverTLSConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	return
}

func ClientTLSConf(caPems ...CAPem) (clientTLSConf *tls.Config) {

	certpool := x509.NewCertPool()
	for _, pem := range caPems {
		certpool.AppendCertsFromPEM(pem.PEM)
	}
	clientTLSConf = &tls.Config{
		RootCAs: certpool,
	}

	return
}

func ExampleCACert(ca_data CertData) (CA_X509 x509.Certificate) {

	CA_X509 = x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               GetSubject(ca_data),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(ca_data.NotAfter.Years, ca_data.NotAfter.Months, ca_data.NotAfter.Days),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	return
}

func ExampleTLSCert(ca_data CertData) (CA_X509 x509.Certificate) {

	CA_X509 = x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject:      GetSubject(ca_data),
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(ca_data.NotAfter.Years, ca_data.NotAfter.Months, ca_data.NotAfter.Days),

		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	return
}
