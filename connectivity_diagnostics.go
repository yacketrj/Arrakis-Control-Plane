package main

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const diagnosticStageTimeout = 8 * time.Second

type connectivityDiagnosticStage struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
	Error  string `json:"error,omitempty"`
}

type connectivityDiagnosticsPayload struct {
	OK         bool                          `json:"ok"`
	Mode       string                        `json:"tunnel_mode"`
	Runtime    string                        `json:"runtime"`
	Stages     []connectivityDiagnosticStage `json:"stages"`
	NextAction string                        `json:"next_action,omitempty"`
}

type dbDiscoveryResult struct {
	ns   string
	pod  string
	host string
	err  error
}

func diagnosticOK(name, detail string) connectivityDiagnosticStage {
	return connectivityDiagnosticStage{Name: name, OK: true, Detail: detail}
}

func diagnosticFail(name string, err error, nextAction string, stages []connectivityDiagnosticStage) connectivityDiagnosticsPayload {
	stages = append(stages, connectivityDiagnosticStage{Name: name, OK: false, Error: redactDiagnosticText(err.Error())})
	return connectivityDiagnosticsPayload{OK: false, Mode: normalizedTunnelMode(), Runtime: normalizeRuntime(serverRuntime), Stages: stages, NextAction: nextAction}
}

func redactDiagnosticText(value string) string {
	redacted := value
	secretValues := []string{dbPass, adminToken, effectiveAdminToken()}
	for _, secret := range secretValues {
		secret = strings.TrimSpace(secret)
		if secret != "" {
			redacted = strings.ReplaceAll(redacted, secret, "<redacted>")
		}
	}
	return redacted
}

func runRemoteDiagnostic(client *ssh.Client, command string) (string, error) {
	sess, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	type result struct {
		out string
		err error
	}
	ch := make(chan result, 1)
	go func() {
		out, err := sess.CombinedOutput(command)
		ch <- result{out: strings.TrimSpace(string(out)), err: err}
	}()

	select {
	case res := <-ch:
		return res.out, res.err
	case <-time.After(diagnosticStageTimeout):
		_ = sess.Close()
		return "", fmt.Errorf("remote command timed out after %s", diagnosticStageTimeout)
	}
}

func diagnosticRuntimeDetail(client *ssh.Client) string {
	out, err := runRemoteDiagnostic(client, "bash -lc 'printf kubectl=; command -v kubectl >/dev/null && printf yes || printf no; printf \" docker=\"; command -v docker >/dev/null && printf yes || printf no; printf \" battlegroup=\"; command -v battlegroup >/dev/null && printf yes || printf no'")
	if err != nil || strings.TrimSpace(out) == "" {
		return "remote runtime tools could not be queried"
	}
	return out
}

func discoverDBPodWithTimeout(client *ssh.Client) (string, string, string, error) {
	ch := make(chan dbDiscoveryResult, 1)
	go func() {
		ns, pod, host, err := discoverDBPod(client)
		ch <- dbDiscoveryResult{ns: ns, pod: pod, host: host, err: err}
	}()

	select {
	case res := <-ch:
		return res.ns, res.pod, res.host, res.err
	case <-time.After(diagnosticStageTimeout):
		return "", "", "", fmt.Errorf("database discovery timed out after %s", diagnosticStageTimeout)
	}
}

func sshDialRemoteWithTimeout(client *ssh.Client, remoteAddr string) (net.Conn, error) {
	type result struct {
		conn net.Conn
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		conn, err := client.Dial("tcp", remoteAddr)
		ch <- result{conn: conn, err: err}
	}()

	select {
	case res := <-ch:
		return res.conn, res.err
	case <-time.After(diagnosticStageTimeout):
		return nil, fmt.Errorf("remote TCP dial timed out after %s", diagnosticStageTimeout)
	}
}

func runConnectivityDiagnostics() connectivityDiagnosticsPayload {
	var stages []connectivityDiagnosticStage

	if errs := requiredConfigErrors(); len(errs) > 0 {
		stages = append(stages, connectivityDiagnosticStage{Name: "config", OK: false, Error: redactDiagnosticText(strings.Join(errs, "; "))})
		return connectivityDiagnosticsPayload{OK: false, Mode: normalizedTunnelMode(), Runtime: normalizeRuntime(serverRuntime), Stages: stages, NextAction: "Run setup and correct invalid configuration values."}
	}
	stages = append(stages, diagnosticOK("config", "required configuration values passed validation"))

	keyPath, err := resolveValidatedKeyPath()
	if err != nil {
		return diagnosticFail("ssh_key", err, "Verify the configured SSH private key path and file permissions.", stages)
	}
	stages = append(stages, diagnosticOK("ssh_key", fmt.Sprintf("readable key file: %s", filepath.Base(keyPath))))

	client, err := dialSSH(keyPath)
	if err != nil {
		return diagnosticFail("ssh_dial", err, "Verify SSH host, port, username, private key, firewall, and bastion reachability.", stages)
	}
	defer client.Close()
	stages = append(stages, diagnosticOK("ssh_dial", "SSH connection established"))

	stages = append(stages, diagnosticOK("remote_runtime_tools", diagnosticRuntimeDetail(client)))

	ns, pod, host, err := discoverDBPodWithTimeout(client)
	if err != nil {
		return diagnosticFail("db_discovery", err, "Verify runtime mode, DB pod/container naming, battlegroup state, and discovery commands on the remote host.", stages)
	}
	stages = append(stages, diagnosticOK("db_discovery", fmt.Sprintf("namespace=%s pod=%s endpoint=%s", ns, pod, host)))

	remoteAddr := fmt.Sprintf("%s:%d", host, dbPort)
	remoteConn, err := sshDialRemoteWithTimeout(client, remoteAddr)
	if err != nil {
		return diagnosticFail("remote_db_tcp", err, "DB endpoint was discovered but is not reachable through SSH from the remote host context.", stages)
	}
	_ = remoteConn.Close()
	stages = append(stages, diagnosticOK("remote_db_tcp", "remote DB endpoint accepted TCP connection through SSH client"))

	bindHost := strings.TrimSpace(sshTunnelHost)
	if bindHost == "" {
		bindHost = "127.0.0.1"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bindHost, 0))
	if err != nil {
		return diagnosticFail("local_tunnel_bind", err, "Local tunnel bind failed. Verify local bind host and whether another process or policy blocks local listening sockets.", stages)
	}
	localAddr := listener.Addr().String()
	_ = listener.Close()
	stages = append(stages, diagnosticOK("local_tunnel_bind", fmt.Sprintf("local bind available at %s", localAddr)))

	return connectivityDiagnosticsPayload{OK: true, Mode: normalizedTunnelMode(), Runtime: normalizeRuntime(serverRuntime), Stages: stages}
}

func handleConnectivityDiagnostics(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, runConnectivityDiagnostics())
}
