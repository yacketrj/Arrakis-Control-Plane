package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func runSetup() {
	r := bufio.NewReader(os.Stdin)

	prompt := func(label, def string) string {
		if def != "" {
			fmt.Printf("  %s [%s]: ", label, def)
		} else {
			fmt.Printf("  %s: ", label)
		}
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			return def
		}
		return line
	}

	ok := func(msg string) { fmt.Printf("  ✓ %s\n", msg) }
	fail := func(msg string) { fmt.Printf("  ✗ %s\n", msg) }

	fmt.Println()
	fmt.Println("=== Arrakis Control Panel setup ===")
	fmt.Println()

	fmt.Println("Checking for Ed25519 SSH key...")
	keyPath := resolveKeyPath()
	if _, err := os.Stat(keyPath); err != nil {
		fail("Ed25519 SSH key not found (checked ~/.ssh/dune_admin_ed25519 and ~/.ssh/id_ed25519)")
		fmt.Println()
		sshKeyPath = expandLocalPath(prompt("Path to Ed25519 SSH private key", ""))
		if sshKeyPath == "" {
			fmt.Fprintln(os.Stderr, "Ed25519 SSH key is required. Aborting.")
			os.Exit(1)
		}
		if _, err := os.Stat(sshKeyPath); err != nil {
			fmt.Fprintf(os.Stderr, "Key not found at %s. Aborting.\n", sshKeyPath)
			os.Exit(1)
		}
		keyPath = sshKeyPath
	} else {
		ok("Ed25519 SSH key: " + keyPath)
		sshKeyPath = keyPath
	}
	fmt.Println()

	fmt.Println("SSH connection:")
	sshHost = prompt("VM host:port", sshHost)
	sshUser = prompt("SSH user", sshUser)
	sshKnownHostsPath = prompt("SSH known_hosts path", envOr("SSH_KNOWN_HOSTS", defaultKnownHostsPath()))
	if _, err := resolveKnownHostsPath(); err != nil {
		fail("known_hosts validation failed: " + err.Error())
		fmt.Println()
		fmt.Println("  Add the remote Ed25519 host key before continuing, for example:")
		fmt.Println("    " + sshKeyscanHint())
		os.Exit(1)
	}
	ok("SSH known_hosts: " + expandLocalPath(sshKnownHostsPath))
	fmt.Println()

	fmt.Println("Database connection:")
	dbName = prompt("DB name", dbName)
	dbPortText := prompt("DB port", fmt.Sprintf("%d", dbPort))
	parsedDBPort, err := strconv.Atoi(strings.TrimSpace(dbPortText))
	if err != nil || parsedDBPort <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid DB port %q. Aborting.\n", dbPortText)
		os.Exit(1)
	}
	dbPort = parsedDBPort
	dbUser = prompt("DB user", dbUser)
	newDBPass := prompt("DB password", "")
	if newDBPass != "" {
		dbPass = newDBPass
	}
	if strings.TrimSpace(dbPass) == "" {
		fmt.Fprintln(os.Stderr, "DB password is required. Aborting.")
		os.Exit(1)
	}
	fmt.Println()

	fmt.Println("Admin security:")
	if strings.TrimSpace(adminToken) == "" {
		provided := prompt("ADMIN_TOKEN (press Enter to generate a strict 43-character token)", "")
		if strings.TrimSpace(provided) == "" {
			adminToken = generateAdminToken()
			ok("Generated strict ADMIN_TOKEN and will save it to .env")
		} else {
			adminToken = provided
			if err := validateStrictAdminToken(adminToken); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid ADMIN_TOKEN: %s\n", err)
				os.Exit(1)
			}
			ok("Using provided strict ADMIN_TOKEN")
		}
	} else {
		if err := validateStrictAdminToken(adminToken); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid existing ADMIN_TOKEN: %s\n", err)
			os.Exit(1)
		}
		ok("ADMIN_TOKEN already configured and valid")
	}
	ok("Admin reason enforcement default: false; set ADMIN_REQUIRE_REASON=true after UI reason prompts are fully wired")
	fmt.Println()

	if errs := requiredConfigErrors(); len(errs) > 0 {
		fail("Configuration validation failed:")
		for _, err := range errs {
			fmt.Println("    - " + err)
		}
		os.Exit(1)
	}

	fmt.Printf("Connecting via SSH to %s...\n", sshHost)
	client, err := dialSSH(keyPath)
	if err != nil {
		fail("SSH failed: " + err.Error())
		fmt.Println()
		fmt.Println("  Make sure:")
		fmt.Println("    - The VM is reachable at the given host:port")
		fmt.Println("    - The SSH key is Ed25519 and authorized on the VM for that user")
		fmt.Println("    - SSH_KNOWN_HOSTS contains the remote Ed25519 host key")
		fmt.Println("    - The SSH user has passwordless sudo for kubectl, or Docker is available")
		os.Exit(1)
	}
	ok("SSH connected with Ed25519 key and verified host identity")
	fmt.Println()

	if err := writeSetupEnv(true); err != nil {
		fail("Failed to write SSH and DB config to .env: " + err.Error())
	} else {
		ok("SSH, DB, and ADMIN_TOKEN config saved to .env")
	}
	fmt.Println()

	fmt.Println("Discovering database endpoint...")
	ns, pod, podIP, err := discoverDBPod(client)
	if err != nil {
		fail("Database endpoint discovery failed: " + err.Error())
		fmt.Println()
		fmt.Println("  SSH succeeded, but no supported Dune database service was found through Kubernetes or Docker.")
		fmt.Println("  Verify the game server stack is running on the target host and then re-run setup.")
		os.Exit(1)
	}
	globalSSH = client
	globalPodNS = ns
	globalPod = pod
	globalPodIP = podIP
	ok("Runtime detected: " + normalizeRuntime(serverRuntime))
	ok("Database endpoint: " + pod)
	fmt.Println()

	fmt.Println("Connecting to database...")
	pool, err := connectDB(context.Background(), dbUser, dbPass)
	if err != nil {
		fail("DB connect failed: " + err.Error())
		fmt.Println()
		fmt.Println("  SSH and endpoint discovery succeeded, but PostgreSQL rejected or closed the connection.")
		fmt.Println("  Verify the DB user, DB password, database name, and detected database port.")
		os.Exit(1)
	}
	globalDB = pool
	ok("Database connected as: " + dbUser)
	fmt.Println()

	fmt.Println("Server config:")
	listenAddr = prompt("HTTP listen address", listenAddr)
	if err := validateListenExposure(listenAddr); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid LISTEN_ADDR: %s\n", err)
		os.Exit(1)
	}
	fmt.Println()

	if errs := requiredConfigErrors(); len(errs) > 0 {
		fail("Configuration validation failed:")
		for _, err := range errs {
			fmt.Println("    - " + err)
		}
		os.Exit(1)
	}

	if err := writeSetupEnv(true); err != nil {
		fail("Failed to write .env: " + err.Error())
		os.Exit(1)
	}
	ok(".env written with saved credentials, runtime, SSH trust, and ADMIN_TOKEN")
	fmt.Println()

	fmt.Println("Setup complete.")
	fmt.Println()
	fmt.Println("  Build and run:  make build && ./arrakis-control-panel")
	fmt.Println("  Run (no build): go run .")
	fmt.Println("  Frontend token: paste ADMIN_TOKEN from .env into the frontend settings gear")
	fmt.Println()
}

func writeSetupEnv(includeDatabasePassword bool) error {
	quote := func(value string) string {
		if strings.ContainsAny(value, " \t#\"'") {
			return "\"" + strings.ReplaceAll(value, "\"", "\\\"") + "\""
		}
		return value
	}
	lines := []string{
		"# Generated by: arrakis-control-panel -setup",
		"",
		"SERVER_RUNTIME=" + quote(normalizeRuntime(serverRuntime)),
		"",
		"SSH_HOST=" + quote(sshHost),
		"SSH_USER=" + quote(sshUser),
		"SSH_KEY=" + quote(keyPathForEnv()),
		"SSH_KNOWN_HOSTS=" + quote(knownHostsPathForEnv()),
		"SSH_TUNNEL_MODE=" + quote(normalizedTunnelMode()),
		"SSH_TUNNEL_LOCAL_HOST=" + quote(sshTunnelHost),
		fmt.Sprintf("DB_TUNNEL_LOCAL_PORT=%d", dbTunnelLocalPort),
		"",
		fmt.Sprintf("DB_PORT=%d", dbPort),
		"DB_USER=" + quote(dbUser),
	}
	if includeDatabasePassword {
		lines = append(lines, "DB_PASS="+quote(dbPass))
	} else {
		lines = append(lines, "DB_PASS=")
	}
	lines = append(lines,
		"DB_NAME="+quote(dbName),
		"DB_SCHEMA="+quote(dbSchema),
		"",
		fmt.Sprintf("SCRIP_CURRENCY=%d", scripCurrencyID),
		"ADMIN_TOKEN="+quote(effectiveAdminToken()),
		"ADMIN_REQUIRE_REASON="+quote(envOr("ADMIN_REQUIRE_REASON", "false")),
		"ALLOWED_ORIGINS="+quote(allowedOrigins),
		"LISTEN_ADDR="+quote(listenAddr),
	)
	return os.WriteFile(".env", []byte(strings.Join(lines, "\n")+"\n"), 0600)
}

func keyPathForEnv() string {
	if sshKeyPath != "" {
		return expandLocalPath(sshKeyPath)
	}
	return resolveKeyPath()
}

func knownHostsPathForEnv() string {
	if sshKnownHostsPath != "" {
		return expandLocalPath(sshKnownHostsPath)
	}
	return defaultKnownHostsPath()
}
