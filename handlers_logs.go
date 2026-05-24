package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return originAllowed(r.Header.Get("Origin")) },
}

type logPod struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func handleLogPods(w http.ResponseWriter, r *http.Request) {
	cmd := runtimeLogPodsCommand()
	out, err := sshExec(cmd)
	if err != nil {
		jsonErr(w, fmt.Errorf("runtime %s log target discovery failed: %w — output: %s", normalizeRuntime(serverRuntime), err, out), 500)
		return
	}

	var pods []logPod
	if runtimeUsesDockerCommands() {
		for _, line := range splitLines(out) {
			id := strings.TrimSpace(line)
			if id != "" && isValidRuntimeLogTarget("docker", id) {
				pods = append(pods, logPod{Namespace: normalizeRuntime(serverRuntime), Name: id})
			}
		}
	} else {
		for _, line := range splitLines(out) {
			name := strings.TrimSpace(line)
			if name != "" && isValidK8sName(name) && !strings.Contains(name, "db-dbdepl") {
				pods = append(pods, logPod{Namespace: globalPodNS, Name: name})
			}
		}
		out2, _ := sshExec("sudo kubectl get pods -n funcom-operators --no-headers -o custom-columns=NAME:.metadata.name 2>&1")
		for _, line := range splitLines(out2) {
			name := strings.TrimSpace(line)
			if name != "" && isValidK8sName(name) {
				pods = append(pods, logPod{Namespace: "funcom-operators", Name: name})
			}
		}
	}
	if pods == nil {
		pods = []logPod{}
	}
	jsonOK(w, pods)
}

func handleLogStream(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("ns")
	pod := r.URL.Query().Get("pod")
	if ns == "" || pod == "" {
		http.Error(w, "ns and pod required", 400)
		return
	}
	if !isValidRuntimeLogTarget(ns, pod) {
		http.Error(w, "invalid log stream target", 400)
		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	cmd := runtimeLogStreamCommand(ns, pod)
	ch, cancel, err := sshStream(cmd)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("error: "+err.Error()))
		return
	}
	defer cancel()

	for line := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			return
		}
	}
}

func splitLines(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n")
}

func handleGetCheatLog(w http.ResponseWriter, r *http.Request) {
	msg, ok := cmdFetchCheatLogFixed()().(msgCheatLog)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	rows := msg.rows
	if rows == nil {
		rows = []cheatEntry{}
	}
	jsonOK(w, rows)
}
