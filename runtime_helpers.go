package main

import (
	"net"
	"net/http"
	"strings"
)

func limitBody(w http.ResponseWriter, r *http.Request, maxBytes int64) {
	if r != nil && r.Body != nil {
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	}
}

func isLoopbackAddr(addr string) bool {
	host, _, err := net.SplitHostPort(strings.TrimSpace(addr))
	if err != nil {
		return false
	}
	if host == "" || strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func isAllowedLogTarget(namespace string, target string) bool {
	namespace = strings.TrimSpace(namespace)
	target = strings.TrimSpace(target)
	if target == "" || !isValidK8sName(target) {
		return false
	}
	if namespace == globalPodNS || namespace == "funcom-operators" {
		return true
	}
	return false
}
