package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if originAllowed(origin) {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Admin-Token, X-Admin-Reason")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logStartupSummary(addr string) {
	log.Printf("dune-admin listening on %s", addr)
	log.Printf("runtime: %s", normalizeRuntime(serverRuntime))
	log.Printf("ssh: %s@%s", sshUser, sshHost)

	mode := normalizedTunnelMode()
	tunnels := tunnelStatus()
	if len(tunnels) == 0 {
		log.Printf("ssh tunnel: none active (mode=%s)", mode)
		return
	}

	log.Printf("ssh tunnel mode: %s", mode)
	for _, tunnel := range tunnels {
		name := tunnel["name"]
		if name == "" {
			name = "unnamed"
		}
		log.Printf("ssh tunnel %s: %s -> %s", name, tunnel["local_addr"], tunnel["remote_addr"])
	}
}

func startServer(addr string) {
	addr = normalizeListenAddr(addr)
	if err := validateListenExposure(addr); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	registerRoutes(mux)

	logStartupSummary(addr)
	server := &http.Server{
		Addr:              addr,
		Handler:           corsMiddleware(auditMiddleware(mutationSafetyMiddleware(authMiddleware(mux)))),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonErr(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": RedactSensitiveText(err.Error())})
}

func decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func handlePublicStatus(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]any{
		"service": "dune-admin",
		"status":  "online",
	})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, buildStatusPayload())
}

func handleReconnect(w http.ResponseWriter, r *http.Request) {
	if globalDB != nil {
		globalDB.Close()
		globalDB = nil
	}
	closeManagedTunnels()
	if globalSSH != nil {
		globalSSH.Close()
		globalSSH = nil
	}
	msg, ok := cmdConnect().(msgConnect)
	if !ok || msg.err != nil {
		var errMsg string
		if msg.err != nil {
			errMsg = msg.err.Error()
		}
		jsonErr(w, fmt.Errorf("%s", errMsg), 500)
		return
	}
	if msg, ok := cmdFetchItemTemplates().(msgItemTemplates); ok {
		mergeItemTemplates(msg.templates)
	}
	handleStatus(w, r)
}
