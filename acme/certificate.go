package acme

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func (client *Client) GetCertificate(certUrl string) (cert *x509.Certificate, err error) {
	_, body, err := client.get(certUrl)
	if err != nil {
		return
	}
	p, _ := pem.Decode(body)
	if p == nil {
		err = fmt.Errorf("no PEM data found")
		return
	}
	return x509.ParseCertificate(p.Bytes)
}

func (client *Client) RevokeCert(cert *x509.Certificate, reason int) (err error) {
	revokeReq := RevokeCertRequest{
		Certificate: base64.RawURLEncoding.EncodeToString(cert.Raw),
		Reason:      reason,
	}
	_, _, err = client.post(client.Directory.RevokeCert, revokeReq)
	return
}
