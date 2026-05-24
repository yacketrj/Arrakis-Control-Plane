package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

var (
	managedTunnelsMu sync.Mutex
	managedTunnels   []*sshTunnel
)

type sshTunnel struct {
	name       string
	localAddr  string
	remoteAddr string
	listener   net.Listener
	client     *ssh.Client
	done       chan struct{}
	closeOnce  sync.Once
}

func normalizedTunnelMode() string {
	mode := strings.ToLower(strings.TrimSpace(sshTunnelMode))
	switch mode {
	case "", "auto", "on", "managed":
		return "auto"
	case "existing", "external", "local":
		return "existing"
	case "off", "direct", "disabled", "false", "0":
		return "off"
	default:
		return "auto"
	}
}

func newManagedTunnel(client *ssh.Client, name, localHost string, localPort int, remoteAddr string) (*sshTunnel, error) {
	if strings.TrimSpace(localHost) == "" {
		localHost = "127.0.0.1"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localHost, localPort))
	if err != nil {
		return nil, fmt.Errorf("listen for %s tunnel: %w", name, err)
	}
	tunnel := &sshTunnel{
		name:       name,
		localAddr:  listener.Addr().String(),
		remoteAddr: remoteAddr,
		listener:   listener,
		client:     client,
		done:       make(chan struct{}),
	}
	go tunnel.serve()
	managedTunnelsMu.Lock()
	managedTunnels = append(managedTunnels, tunnel)
	managedTunnelsMu.Unlock()
	return tunnel, nil
}

func dialManagedTunnel(client *ssh.Client, name, remoteAddr string) (net.Conn, error) {
	tunnel, err := newManagedTunnel(client, name, sshTunnelHost, 0, remoteAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("tcp", tunnel.localAddr)
	if err != nil {
		tunnel.close()
		return nil, err
	}
	return conn, nil
}

func (t *sshTunnel) serve() {
	defer close(t.done)
	for {
		localConn, err := t.listener.Accept()
		if err != nil {
			return
		}
		go t.proxy(localConn)
	}
}

func (t *sshTunnel) proxy(localConn net.Conn) {
	defer localConn.Close()
	remoteConn, err := t.client.Dial("tcp", t.remoteAddr)
	if err != nil {
		return
	}
	defer remoteConn.Close()
	done := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(remoteConn, localConn); done <- struct{}{} }()
	go func() { _, _ = io.Copy(localConn, remoteConn); done <- struct{}{} }()
	<-done
}

func (t *sshTunnel) close() {
	t.closeOnce.Do(func() {
		_ = t.listener.Close()
		<-t.done
	})
}

func closeManagedTunnels() {
	managedTunnelsMu.Lock()
	tunnels := managedTunnels
	managedTunnels = nil
	managedTunnelsMu.Unlock()
	for _, tunnel := range tunnels {
		tunnel.close()
	}
}

func tunnelStatus() []map[string]string {
	managedTunnelsMu.Lock()
	defer managedTunnelsMu.Unlock()
	status := make([]map[string]string, 0, len(managedTunnels))
	for _, tunnel := range managedTunnels {
		status = append(status, map[string]string{
			"name":        tunnel.name,
			"local_addr":  tunnel.localAddr,
			"remote_addr": tunnel.remoteAddr,
		})
	}
	return status
}
