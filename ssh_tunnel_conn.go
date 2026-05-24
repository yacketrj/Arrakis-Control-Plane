package main

import (
	"net"

	"golang.org/x/crypto/ssh"
)

type managedTunnelConn struct {
	net.Conn
	tunnel *sshTunnel
}

func (c *managedTunnelConn) Close() error {
	err := c.Conn.Close()
	c.tunnel.close()
	return err
}

func dialManagedTunnelConn(client *ssh.Client, name, remoteAddr string) (net.Conn, error) {
	tunnel, err := newManagedTunnel(client, name, sshTunnelHost, 0, remoteAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("tcp", tunnel.localAddr)
	if err != nil {
		tunnel.close()
		return nil, err
	}
	return &managedTunnelConn{Conn: conn, tunnel: tunnel}, nil
}
