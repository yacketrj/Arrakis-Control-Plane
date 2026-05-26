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
	out, err := sess.CombinedOutput(command)
	return strings.TrimSpace(string(out)), err
}

func diagnosticRuntimeDetail(client *ssh.Client) string {
	out, err := runRemoteDiagnostic(client, "bash -lc 'printf kubectl=; command -v kubectl >/dev/null && printf yes || printf no; printf \" docker=\"; command -v docker >/dev/null && printf yes || printf no; printf \" battlegroup=\"; command -v battlegroup >/dev/null && printf yes || printf no'")
	if err != nil || strings.TrimSpace(out) == "" {
		return "remote runtime tools could not be queried"
	}
	return out
}

func handleConnectivityDiagnostics(w http.ResponseWriter, r *http.Request) {
	var stages []connectivityDiagnosticStage

	if errs := requiredConfigErrors(); len(errs) > 0 {
		stages = append(stages, connectivityDiagnosticStage{Name: "config", OK: false, Error: redactDiagnosticText(strings.Join(errs, "; "))})
		jsonOK(w, connectivityDiagnosticsPayload{OK: false, Mode: normalizedTunnelMode(), Runtime: normalizeRuntime(serverRuntime), Stages: stages, NextAction: "Run setup and correct invalid configuration values."})
		return
	}
	stages = append(stages, diagnosticOK("config", "required configuration values passed validation"))

	keyPath, err := resolveValidatedKeyPath()
	if err != nil {
		jsonOK(w, diagnosticFail("ssh_key", err, "Verify the configured SSH private key path and file permissions.", stages))
		return
	}
	stages = append(stages, diagnosticOK("ssh_key", fmt.Sprintf("readable key file: %s", filepath.Base(keyPath))))

	client, err := dialSSH(keyPath)
	if err != nil {
		jsonOK(w, diagnosticFail("ssh_dial", err, "Verify SSH host, port, username, private key, firewall, and bastion reachability.", stages))
		return
	}
	defer client.Close()
	stages = append(stages, diagnosticOK("ssh_dial", "SSH connection established"))

	stages = append(stages, diagnosticOK("remote_runtime_tools", diagnosticRuntimeDetail(client)))

	ns, pod, host, err := discoverDBPod(client)
	if err != nil {
		jsonOK(w, diagnosticFail("db_discovery", err, "Verify runtime mode, DB pod/container naming, battlegroup state, and discovery commands on the remote host.", stages))
		return
	}
	stages = append(stages, diagnosticOK("db_discovery", fmt.Sprintf("namespace=%s pod=%s endpoint=%s", ns, pod, host)))

	remoteAddr := fmt.Sprintf("%s:%d", host, dbPort)
	remoteConn, err := client.Dial("tcp", remoteAddr)
	if err != nil {
		jsonOK(w, diagnosticFail("remote_db_tcp", err, "DB endpoint was discovered but is not reachable through SSH from the remote host context.", stages))
		return
	}
	_ = remoteConn.Close()
	stages = append(stages, diagnosticOK("remote_db_tcp", "remote DB endpoint accepted TCP connection through SSH client"))

	bindHost := strings.TrimSpace(sshTunnelHost)
	if bindHost == "" {
		bindHost = "127.0.0.1"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bindHost, 0))
	if err != nil {
		jsonOK(w, diagnosticFail("local_tunnel_bind", err, "Local tunnel bind failed. Verify local bind host and whether another process or policy blocks local listening sockets.", stages))
		return
	}
	localAddr := listener.Addr().String()
	_ = listener.Close()
	stages = append(stages, diagnosticOK("local_tunnel_bind", fmt.Sprintf("local bind available at %s", localAddr)))

	_ = time.Now()
	jsonOK(w, connectivityDiagnosticsPayload{OK: true, Mode: normalizedTunnelMode(), Runtime: normalizeRuntime(serverRuntime), Stages: stages})
}
