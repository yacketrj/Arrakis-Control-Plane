package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var version = "dev" // set by goreleaser ldflags

// ── config ────────────────────────────────────────────────────────────────────

var (
	captureMode       bool
	setupMode         bool
	diagnoseMode      bool
	serverRuntime     string
	sshHost           string
	sshUser           string
	sshKeyPath        string
	sshKnownHostsPath string
	sshTunnelMode     string
	sshTunnelHost     string
	dbTunnelLocalPort int
	itemDataPath      string
	scripCurrencyID   int
	dbPort            int
	dbUser            string
	dbPass            string
	dbName            string
	dbSchema          string
	listenAddr        string
	startupConnectErr string
)

func init() {
	loadDotEnv()
	flag.StringVar(&serverRuntime, "runtime", envOr("SERVER_RUNTIME", "auto"), "Server runtime mode: auto, kubernetes, or docker")
	flag.StringVar(&sshHost, "host", envOr("SSH_HOST", "192.168.0.72:22"), "SSH host:port")
	flag.StringVar(&sshUser, "user", envOr("SSH_USER", "dune"), "SSH user")
	flag.StringVar(&sshKeyPath, "key", envOr("SSH_KEY", ""), "SSH Ed25519 private key path (auto-detected if empty)")
	flag.StringVar(&sshKnownHostsPath, "knownhosts", envOr("SSH_KNOWN_HOSTS", ""), "SSH known_hosts file path containing the remote Ed25519 host key (defaults to ~/.ssh/known_hosts)")
	flag.StringVar(&sshTunnelMode, "sshtunnel", envOr("SSH_TUNNEL_MODE", "auto"), "SSH tunnel mode for game-management traffic: auto, existing, or off")
	flag.StringVar(&sshTunnelHost, "tunnelhost", envOr("SSH_TUNNEL_LOCAL_HOST", "127.0.0.1"), "Local bind host for SSH tunnels")
	flag.IntVar(&dbTunnelLocalPort, "dbtunnelport", envIntOr("DB_TUNNEL_LOCAL_PORT", 0), "Local DB tunnel port; 0 chooses an available port in auto mode")
	flag.StringVar(&itemDataPath, "itemdata", envOr("ITEM_DATA", ""), "Item data JSON path")
	flag.IntVar(&scripCurrencyID, "scripcurrency", envIntOr("SCRIP_CURRENCY", 1), "Scrip currency id")
	flag.IntVar(&dbPort, "dbport", envIntOr("DB_PORT", 15432), "PostgreSQL port inside the selected runtime")
	flag.StringVar(&dbUser, "dbuser", envOr("DB_USER", "dune"), "PostgreSQL user")
	flag.StringVar(&dbPass, "dbpass", envOr("DB_PASS", ""), "PostgreSQL password")
	flag.StringVar(&dbName, "dbname", envOr("DB_NAME", "dune"), "PostgreSQL database name")
	flag.StringVar(&dbSchema, "schema", envOr("DB_SCHEMA", "dune"), "PostgreSQL schema")
	flag.StringVar(&listenAddr, "addr", envOr("LISTEN_ADDR", ":8080"), "HTTP listen address")
	flag.BoolVar(&captureMode, "capture", false, "Capture RabbitMQ messages (grant + notifications) and print to stdout")
	flag.BoolVar(&setupMode, "setup", false, "Interactive setup wizard — writes .env")
	flag.BoolVar(&diagnoseMode, "diagnose", false, "Run staged SSH/runtime/DB connectivity diagnostics and exit")
}

func resolveKeyPath() string {
	if sshKeyPath != "" {
		return expandLocalPath(sshKeyPath)
	}
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, ".ssh", "dune_admin_ed25519"),
		filepath.Join(home, ".ssh", "id_ed25519"),
		filepath.Join(os.Getenv("USERPROFILE"), ".ssh", "dune_admin_ed25519"),
		filepath.Join(os.Getenv("USERPROFILE"), ".ssh", "id_ed25519"),
	}
	for _, p := range candidates {
		p = expandLocalPath(p)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return expandLocalPath(candidates[0])
}

func resolveItemDataPath() string {
	if itemDataPath != "" {
		return itemDataPath
	}
	candidates := []string{"./item-data.json", "../item-data.json"}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

var itemData itemDataFile

func loadItemData() error {
	path := resolveItemDataPath()
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read item data %s: %w", path, err)
	}
	var parsed itemDataFile
	if err := json.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("parse item data %s: %w", path, err)
	}
	normalizedItems := make(map[string]itemRule, len(parsed.Items))
	for k, v := range parsed.Items {
		normalizedItems[strings.ToLower(k)] = v
	}
	parsed.Items = normalizedItems
	normalizedNames := make(map[string]string, len(parsed.Names))
	for k, v := range parsed.Names {
		normalizedNames[strings.ToLower(k)] = v
	}
	parsed.Names = normalizedNames
	itemData = parsed
	return nil
}

// ── main ──────────────────────────────────────────────────────────────────────

func needsSetup() bool {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return true
	}
	return len(requiredConfigErrors()) > 0
}

func printConfigErrors(errs []string) {
	if len(errs) == 0 {
		return
	}
	fmt.Fprintln(os.Stderr, "Configuration is missing or invalid:")
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, "  - "+err)
	}
}

func printDiagnostics(payload connectivityDiagnosticsPayload) {
	fmt.Printf("Connectivity diagnostics: ok=%t runtime=%s tunnel_mode=%s\n", payload.OK, payload.Runtime, payload.Mode)
	for _, stage := range payload.Stages {
		state := "FAIL"
		if stage.OK {
			state = "OK"
		}
		fmt.Printf("[%s] %s", state, stage.Name)
		if stage.Detail != "" {
			fmt.Printf(" - %s", stage.Detail)
		}
		if stage.Error != "" {
			fmt.Printf(" - %s", stage.Error)
		}
		fmt.Println()
	}
	if payload.NextAction != "" {
		fmt.Println("Next action:", payload.NextAction)
	}
}

func main() {
	flag.Parse()
	serverRuntime = normalizeRuntime(serverRuntime)

	// Explicit -setup flag: reconfigure and exit (don't start server).
	if setupMode {
		runSetup()
		return
	}

	if diagnoseMode {
		payload := runConnectivityDiagnostics()
		printDiagnostics(payload)
		if !payload.OK {
			os.Exit(1)
		}
		return
	}

	if captureMode {
		if msg, ok := cmdConnect().(msgConnect); ok && msg.err != nil {
			fmt.Fprintln(os.Stderr, "SSH connect:", msg.err)
			os.Exit(1)
		}
		runCapture()
		return
	}

	if err := loadItemData(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Auto-run setup when .env is missing or required values are missing/invalid.
	alreadyConnected := false
	if needsSetup() {
		if _, err := os.Stat(".env"); err == nil {
			printConfigErrors(requiredConfigErrors())
		}
		runSetup()
		alreadyConnected = true
		fmt.Println()
		fmt.Printf("Starting server on %s...\n", listenAddr)
	}

	defer func() {
		if globalDB != nil {
			globalDB.Close()
		}
		closeManagedTunnels()
		if globalSSH != nil {
			globalSSH.Close()
		}
	}()

	if !alreadyConnected {
		if errs := requiredConfigErrors(); len(errs) > 0 {
			printConfigErrors(errs)
			fmt.Fprintf(os.Stderr, "Run .\\%s -setup to repair configuration.\n", appWindowsExecutable)
			os.Exit(1)
		}
		// Connect synchronously. If this fails, the API starts in degraded mode so /api/v1/status and /api/v1/reconnect remain available.
		if msg, ok := cmdConnect().(msgConnect); ok && msg.err != nil {
			startupConnectErr = msg.err.Error()
			fmt.Fprintln(os.Stderr, "connect:", startupConnectErr)
			fmt.Fprintln(os.Stderr, "Startup degraded: DB-backed features are disabled until /api/v1/reconnect succeeds.")
		} else {
			if msg, ok := cmdFetchItemTemplates().(msgItemTemplates); ok {
				mergeItemTemplates(msg.templates)
			}
		}
	} else {
		if msg, ok := cmdFetchItemTemplates().(msgItemTemplates); ok {
			mergeItemTemplates(msg.templates)
		}
	}

	startServer(listenAddr)
}
