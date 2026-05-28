package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var (
	globalSSH   *ssh.Client
	globalDB    *pgxpool.Pool
	globalPodIP string
	globalPodNS string
	globalPod   string
)

const requiredSSHKeyAlgorithm = ssh.KeyAlgoED25519

func defaultKnownHostsPath() string {
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return ""
	}
	return filepath.Join(home, ".ssh", "known_hosts")
}

func resolveKnownHostsPath() (string, error) {
	configured := strings.TrimSpace(sshKnownHostsPath)
	if configured != "" {
		return validateReadableFilePath("SSH_KNOWN_HOSTS", configured)
	}
	path := defaultKnownHostsPath()
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("SSH_KNOWN_HOSTS is required because the user home directory could not be resolved")
	}
	return validateReadableFilePath("SSH_KNOWN_HOSTS", path)
}

func sshKeyscanHint() string {
	host, port, err := net.SplitHostPort(sshHost)
	if err != nil || strings.TrimSpace(host) == "" {
		return "ssh-keyscan -t ed25519 -H <host> >> ~/.ssh/known_hosts"
	}
	if port == "22" {
		return fmt.Sprintf("ssh-keyscan -t ed25519 -H %s >> ~/.ssh/known_hosts", host)
	}
	return fmt.Sprintf("ssh-keyscan -t ed25519 -H -p %s %s >> ~/.ssh/known_hosts", port, host)
}

func validateEd25519Signer(signer ssh.Signer) error {
	if signer == nil || signer.PublicKey() == nil {
		return fmt.Errorf("SSH_KEY signer is empty")
	}
	if signer.PublicKey().Type() != requiredSSHKeyAlgorithm {
		return fmt.Errorf("SSH_KEY uses unsupported algorithm %s; only %s keys are allowed", signer.PublicKey().Type(), requiredSSHKeyAlgorithm)
	}
	return nil
}

func ed25519OnlyHostKeyCallback(callback ssh.HostKeyCallback) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if key == nil {
			return fmt.Errorf("SSH host key is empty")
		}
		if key.Type() != requiredSSHKeyAlgorithm {
			return fmt.Errorf("SSH host %s offered unsupported host key algorithm %s; only %s host keys are allowed", hostname, key.Type(), requiredSSHKeyAlgorithm)
		}
		return callback(hostname, remote, key)
	}
}

func newHostKeyCallback() (ssh.HostKeyCallback, string, error) {
	knownHostsPath, err := resolveKnownHostsPath()
	if err != nil {
		return nil, "", fmt.Errorf("%w. Add the remote Ed25519 host key first, for example: %s", err, sshKeyscanHint())
	}
	callback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, "", fmt.Errorf("load SSH known_hosts %s: %w", knownHostsPath, err)
	}
	return ed25519OnlyHostKeyCallback(callback), knownHostsPath, nil
}

func dialSSH(keyPath string) (*ssh.Client, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read key %s: %w", keyPath, err)
	}
	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("parse key: %w", err)
	}
	if err := validateEd25519Signer(signer); err != nil {
		return nil, err
	}
	hostKeyCallback, knownHostsPath, err := newHostKeyCallback()
	if err != nil {
		return nil, err
	}
	client, err := ssh.Dial("tcp", sshHost, &ssh.ClientConfig{
		User:              sshUser,
		Auth:              []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback:   hostKeyCallback,
		HostKeyAlgorithms: []string{requiredSSHKeyAlgorithm},
		Timeout:           15 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("SSH dial using Ed25519 known_hosts %s: %w", knownHostsPath, err)
	}
	return client, nil
}

func discoverDBPod(client *ssh.Client) (ns, pod, podIP string, err error) {
	endpoint, err := discoverDatabaseEndpoint(client)
	if err != nil {
		return "", "", "", err
	}
	applyDetectedRuntime(endpoint.Runtime)
	if endpoint.Port > 0 {
		dbPort = endpoint.Port
	}
	return endpoint.Namespace, endpoint.Name, endpoint.Host, nil
}

func cmdConnect() Msg {
	closeManagedTunnels()
	keyPath := resolveKeyPath()
	client, err := dialSSH(keyPath)
	if err != nil {
		return msgConnect{err: err}
	}

	ns, pod, podIP, err := discoverDBPod(client)
	if err != nil {
		client.Close()
		return msgConnect{err: err}
	}
	globalPodNS = ns
	globalPod = pod
	globalPodIP = podIP
	globalSSH = client

	pool, err := connectDB(context.Background(), dbUser, dbPass)
	if err != nil {
		closeManagedTunnels()
		client.Close()
		globalSSH = nil
		return msgConnect{err: fmt.Errorf("DB connect: %w", err)}
	}
	globalDB = pool
	return msgConnect{}
}

func connectDB(ctx context.Context, user, pass string) (*pgxpool.Pool, error) {
	mode := normalizedTunnelMode()
	host := sshTunnelHost
	port := dbPort
	remoteAddr := fmt.Sprintf("%s:%d", globalPodIP, dbPort)

	switch mode {
	case "auto":
		tunnel, err := newManagedTunnel(globalSSH, "postgres", sshTunnelHost, dbTunnelLocalPort, remoteAddr)
		if err != nil {
			return nil, err
		}
		parsedHost, parsedPort, err := net.SplitHostPort(tunnel.localAddr)
		if err != nil {
			return nil, fmt.Errorf("parse tunnel address: %w", err)
		}
		parsedPortNumber, err := net.LookupPort("tcp", parsedPort)
		if err != nil {
			return nil, fmt.Errorf("parse tunnel port: %w", err)
		}
		host = parsedHost
		port = parsedPortNumber
	case "existing":
		if dbTunnelLocalPort > 0 {
			port = dbTunnelLocalPort
		}
	case "off":
		host = globalPodIP
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName)
	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, fmt.Sprintf(`SET search_path TO %s, public`, pgx.Identifier{dbSchema}.Sanitize()))
		return err
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	dbUser = user
	dbPass = pass
	return pool, nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func battlegroupFromPod(pod string) string {
	const suffix = "-db-dbdepl-sts-"
	if idx := strings.LastIndex(pod, suffix); idx != -1 {
		return pod[:idx]
	}
	return ""
}

func listBattlegroups(client *ssh.Client) []string {
	sess, err := client.NewSession()
	if err != nil {
		return nil
	}
	defer sess.Close()
	out, err := sess.CombinedOutput("bash -lc 'battlegroup list' 2>/dev/null")
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return nil
	}
	var names []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			if name := strings.TrimSpace(line[2:]); name != "" {
				names = append(names, name)
			}
		}
	}
	return names
}

func extractPasswordFromYAML(client *ssh.Client, filePath string) (user, pass string) {
	sess, err := client.NewSession()
	if err != nil {
		return "", ""
	}
	defer sess.Close()
	out, err := sess.CombinedOutput(fmt.Sprintf("cat %s 2>/dev/null", shellQuote(filePath)))
	if err != nil || len(out) == 0 {
		sess2, err2 := client.NewSession()
		if err2 != nil {
			return "", ""
		}
		defer sess2.Close()
		out, err = sess2.CombinedOutput(fmt.Sprintf(`bash -c 'cat %s'`, filePath))
		if err != nil || len(out) == 0 {
			return "", ""
		}
	}
	return parseDeploymentCredentials(out)
}

func sshExec(cmd string) (string, error) {
	if globalSSH == nil {
		return "", fmt.Errorf("not connected")
	}
	sess, err := globalSSH.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()
	out, err := sess.CombinedOutput(cmd)
	return strings.TrimSpace(string(out)), err
}

func sshStream(cmd string) (<-chan string, func(), error) {
	if globalSSH == nil {
		return nil, func() {}, fmt.Errorf("not connected")
	}
	sess, err := globalSSH.NewSession()
	if err != nil {
		return nil, func() {}, err
	}
	pipe, err := sess.StdoutPipe()
	if err != nil {
		sess.Close()
		return nil, func() {}, err
	}
	if err := sess.Start(cmd); err != nil {
		sess.Close()
		return nil, func() {}, err
	}
	ch := make(chan string, 256)
	go func() {
		defer close(ch)
		sc := bufio.NewScanner(pipe)
		for sc.Scan() {
			ch <- sc.Text()
		}
		sess.Wait()
	}()
	cancel := func() { sess.Close() }
	return ch, cancel, nil
}
