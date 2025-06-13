package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"
)

type requestInfo struct {
	Method     string
	Url        string
	Proto      string
	Host       string
	RemoteAddr string
	RequestUri string
	TLS        bool
	Header     map[string][]string
}

func newRequestInfo(r *http.Request) *requestInfo {
	return &requestInfo{
		Method:     r.Method,
		Url:        r.URL.String(),
		Proto:      r.Proto,
		Host:       r.Host,
		RemoteAddr: r.RemoteAddr,
		RequestUri: r.RequestURI,
		TLS:        r.TLS != nil,
		Header:     r.Header,
	}
}

func (ri *requestInfo) String() string {
	var sb strings.Builder
	sb.WriteString("Request Info:\n")
	sb.WriteString(fmt.Sprintf("\tMethod:     %v\n", ri.Method))
	sb.WriteString(fmt.Sprintf("\tUrl:        %v\n", ri.Url))
	sb.WriteString(fmt.Sprintf("\tProto:      %v\n", ri.Proto))
	sb.WriteString(fmt.Sprintf("\tHost:       %v\n", ri.Host))
	sb.WriteString(fmt.Sprintf("\tRemoteAddr: %v\n", ri.RemoteAddr))
	sb.WriteString(fmt.Sprintf("\tRequestUri: %v\n", ri.RequestUri))
	sb.WriteString(fmt.Sprintf("\tTLS:        %v\n", ri.TLS))
	sb.WriteString("\tHeader:\n")
	for key, value := range ri.Header {
		sb.WriteString(fmt.Sprintf("\t\t%v: %v\n", key,value))
	}
	return sb.String()
}