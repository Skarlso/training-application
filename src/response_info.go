package main

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
)

type certInfo struct {
	Subject        string
	Issuer         string
}

type tlsInfo struct {
	ServerName  string
	CertInfos   []*certInfo
}

type responseInfo struct {
	Status     string
	Proto      string
	TlsInfo    *tlsInfo
	Header     map[string][]string
}

func newCertInfos(certs []*x509.Certificate) []*certInfo {
	certInfos := make([]*certInfo, len(certs))
	for i, cert := range certs {
		certInfos[i] = &certInfo{
			Subject:        strings.TrimSpace(cert.Subject.CommonName),
			Issuer:         strings.TrimSpace(cert.Issuer.CommonName),
		}
	} 
	return certInfos
}

func newTLSInfo(s *tls.ConnectionState) *tlsInfo {
	if s == nil {
		return nil
	} else {
		return &tlsInfo{
			ServerName:  s.ServerName,
			CertInfos:  newCertInfos(s.PeerCertificates),
		}
	}
}

func newResponseInfo(r *http.Response) *responseInfo {
	return &responseInfo{
		Status:     r.Status,
		Proto:      r.Proto,
		TlsInfo:    newTLSInfo(r.TLS),
		Header:     r.Header,
	}
}

func (ri *responseInfo) String() string {
	var sb strings.Builder
	sb.WriteString("Response Info:\n")
	sb.WriteString(fmt.Sprintf("\tStatus: %v\n", ri.Status))
	sb.WriteString(fmt.Sprintf("\tProto: %v\n", ri.Proto))
	if ri.TlsInfo == nil {
		sb.WriteString("\tTLS: false\n")
	} else {
		sb.WriteString("\tTLS:\n")
		sb.WriteString(fmt.Sprintf("\t\tServer Name: %s\n", ri.TlsInfo.ServerName))
		sb.WriteString("\t\tCertificates:\n")
		for _, ci := range ri.TlsInfo.CertInfos {
			sb.WriteString(fmt.Sprintf("\t\t\tCertificate Subject: %s - Certificate Issuer: %s\n", ci.Subject, ci.Issuer))
		}
	}
	sb.WriteString("\tHeader:\n")
	for key, value := range ri.Header {
		sb.WriteString(fmt.Sprintf("\t\t%v: %v\n", key, value))
	}
	return sb.String()
}
