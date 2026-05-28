package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const remoteExposureEnv = "DUNE_ADMIN_REMOTE_EXPOSURE"
const remoteExposureRequiredValue = "reverse-proxy-tls"

func remoteExposureAcknowledged() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(remoteExposureEnv)), remoteExposureRequiredValue)
}

func validateListenExposure(addr string) error {
	normalized := normalizeListenAddr(addr)
	if isLoopbackAddr(normalized) {
		return nil
	}
	if remoteExposureAcknowledged() {
		log.Printf("security: LISTEN_ADDR %q binds beyond loopback because %s=%s is explicitly configured; ensure HTTPS/WSS, a trusted reverse proxy or VPN, and firewall allow-listing", normalized, remoteExposureEnv, remoteExposureRequiredValue)
		return nil
	}
	return fmt.Errorf("LISTEN_ADDR %q is not loopback; refusing to start. Keep LISTEN_ADDR on 127.0.0.1, or set %s=%s only when the backend is protected behind HTTPS/WSS and a trusted reverse proxy or VPN", normalized, remoteExposureEnv, remoteExposureRequiredValue)
}
