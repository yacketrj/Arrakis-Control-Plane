package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var bgCmdAllowlist = map[string]bool{
	"start": true, "stop": true, "restart": true,
	"update": true, "backup": true, "restore": true,
}

type battlegroupHealthSection struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Output      string `json:"output"`
	Error       string `json:"error,omitempty"`
}

type battlegroupHealthResponse struct {
	Namespace string                     `json:"namespace"`
	Runtime   string                     `json:"runtime"`
	CheckedAt string                     `json:"checked_at"`
	Sections  []battlegroupHealthSection `json:"sections"`
}

type battlegroupHealthSpec struct {
	Name        string
	Description string
	Command     string
}

func battlegroupHealthSpecs(namespace string) []battlegroupHealthSpec {
	return runtimeHealthSpecs(namespace)
}

func handleBGStatus(w http.ResponseWriter, r *http.Request) {
	if err := validateRuntimeCommandNamespace(); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	cmd := runtimeBGStatusCommand()
	out, err := sshExec(cmd)
	if err != nil {
		jsonErr(w, fmt.Errorf("runtime %s status failed: %w — output: %s", normalizeRuntime(serverRuntime), err, RedactSensitiveText(out)), 500)
		return
	}
	jsonOK(w, map[string]string{"output": RedactSensitiveText(out), "runtime": normalizeRuntime(serverRuntime)})
}

func handleBGHealth(w http.ResponseWriter, r *http.Request) {
	if err := validateRuntimeCommandNamespace(); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}

	response := battlegroupHealthResponse{
		Namespace: globalPodNS,
		Runtime:   normalizeRuntime(serverRuntime),
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
	for _, spec := range battlegroupHealthSpecs(globalPodNS) {
		section := battlegroupHealthSection{
			Name:        spec.Name,
			Description: spec.Description,
			Command:     spec.Command,
		}
		out, err := sshExec(spec.Command)
		section.Output = RedactSensitiveText(out)
		if err != nil {
			section.Error = RedactSensitiveText(err.Error())
		}
		response.Sections = append(response.Sections, section)
	}
	jsonOK(w, response)
}

func handleBGExec(w http.ResponseWriter, r *http.Request) {
	if runtimeUsesDockerCommands() {
		jsonErr(w, fmt.Errorf("battlegroup script commands are not supported for Docker runtimes yet"), http.StatusNotImplemented)
		return
	}

	var req struct {
		Cmd string `json:"cmd"`
	}
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, 400)
		return
	}
	cmd, err := normalizeBattlegroupCommand(req.Cmd)
	if err != nil {
		jsonErr(w, err, 400)
		return
	}
	out, err := sshExec(fmt.Sprintf("sudo ~/.dune/download/scripts/battlegroup.sh %s 2>&1", cmd))
	if err != nil {
		jsonErr(w, fmt.Errorf("exec: %w — output: %s", err, RedactSensitiveText(out)), 500)
		return
	}
	jsonOK(w, map[string]string{"output": RedactSensitiveText(out)})
}

func handleBGPods(w http.ResponseWriter, r *http.Request) {
	if err := validateRuntimeCommandNamespace(); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	cmd := runtimeBGPodsCommand()
	out, err := sshExec(cmd)
	if err != nil {
		jsonErr(w, fmt.Errorf("runtime %s container/pod list failed: %w — output: %s", normalizeRuntime(serverRuntime), err, RedactSensitiveText(out)), 500)
		return
	}
	lines := splitAndRedactLines(out)
	jsonOK(w, map[string]any{"pods": lines, "namespace": globalPodNS, "runtime": normalizeRuntime(serverRuntime)})
}

func normalizeBattlegroupCommand(raw string) (string, error) {
	cmd := strings.ToLower(strings.TrimSpace(raw))
	if cmd == "" {
		return "", fmt.Errorf("cmd required")
	}
	if containsUnsafeControl(cmd) {
		return "", fmt.Errorf("cmd contains unsupported control characters")
	}
	if !bgCmdAllowlist[cmd] {
		return "", fmt.Errorf("unknown command %q", sanitizedAuditString(cmd, 64))
	}
	return cmd, nil
}

func validateRuntimeCommandNamespace() error {
	if runtimeUsesDockerCommands() {
		return nil
	}
	if !isValidK8sName(globalPodNS) {
		return fmt.Errorf("invalid namespace %q", sanitizedAuditString(globalPodNS, 128))
	}
	return nil
}

func splitAndRedactLines(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, RedactSensitiveText(line))
		}
	}
	return out
}
