# Runbook: Configuration and Launch Argument Review

Use this runbook when the suspected issue may involve server configuration, startup arguments, environment variables, map settings, bind addresses, public addresses, ports, passwords, tokens, or runtime-generated config files.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Define What Configuration Is Being Reviewed

Record:

```text
Reason for review:
Affected service/process:
Known working behavior:
Known failing behavior:
Configuration file path, if known:
Startup command source:
Recent config change, if known:
```

Do not change configuration until the current state is captured.

---

## 2. Capture Running Process Arguments

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox' | grep -v grep | sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig' | tee launch-arguments-linux.txt
```

Run on: Windows host PowerShell

```powershell
Get-CimInstance Win32_Process | Where-Object { $_.CommandLine -match 'Dune|Awakening|Sandbox' } | Select-Object ProcessId, CommandLine | Tee-Object -FilePath launch-arguments-windows.txt
```

If secrets appear in Windows output, redact before broad sharing.

---

## 3. Capture Runtime or Service Configuration

Run on: Linux systemd host, only if systemd is used

```bash
systemctl cat "$SERVICE_NAME" | tee service-unit.txt
systemctl show "$SERVICE_NAME" -p ExecStart -p User -p Group -p WorkingDirectory -p Environment --no-pager | tee service-show.txt
```

Run on: Windows service host, only if Windows service is used

```powershell
Get-CimInstance Win32_Service | Where-Object { $_.Name -eq $env:SERVICE_NAME } | Select-Object Name, DisplayName, State, StartName, PathName | Tee-Object -FilePath windows-service-config.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker inspect "$CONTAINER_NAME" --format 'Image={{.Config.Image}} User={{.Config.User}} NetworkMode={{.HostConfig.NetworkMode}}'
docker inspect "$CONTAINER_NAME" --format '{{range .Config.Env}}{{println .}}{{end}}' | sed -E 's/(TOKEN|SECRET|PASSWORD|KEY)=.*/\1=<redacted>/Ig' | tee docker-env.txt
docker inspect "$CONTAINER_NAME" --format '{{json .Args}}' | tee docker-args.txt
docker inspect "$CONTAINER_NAME" --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}' | tee docker-mounts.txt
```

Run in: control panel UI

```text
Record startup command, startup parameters, environment/config fields, visible ports, map settings, and generated config paths shown in the panel.
Do not copy tokens or secrets into case notes.
```

---

## 4. Locate Configuration Files

Run on: Linux host or Linux VM shell

```bash
find "$INSTANCE_PATH" -maxdepth 8 -type f \( -name '*.ini' -o -name '*.env' -o -name '*.json' -o -name '*.yaml' -o -name '*.yml' -o -name '*.cfg' \) 2>/dev/null | sort | tee config-files-linux.txt
```

Run on: Windows host PowerShell

```powershell
Get-ChildItem -Path $env:INSTANCE_PATH -Recurse -ErrorAction SilentlyContinue -Include *.ini,*.env,*.json,*.yaml,*.yml,*.cfg | Select-Object FullName | Tee-Object -FilePath config-files-windows.txt
```

---

## 5. Search for Common Configuration Values

Run on: Linux host or Linux VM shell

```bash
grep -RniE 'Port|IGW|Bind|External|Public|Region|ServerName|DisplayName|Password|Token|Rabbit|Database|Partition|Map|Instancing|FLS|Gateway' "$INSTANCE_PATH" 2>/dev/null | sed -E 's/(Token|Secret|Password|Key)=?[^, ]+/\1=<redacted>/Ig' | head -500 | tee config-value-search.txt
```

Run on: Windows host PowerShell

```powershell
Select-String -Path "$env:INSTANCE_PATH\**\*" -Pattern 'Port','IGW','Bind','External','Public','Region','ServerName','DisplayName','Password','Token','Rabbit','Database','Partition','Map','Instancing','FLS','Gateway' -ErrorAction SilentlyContinue | Select-Object -First 500 | Tee-Object -FilePath config-value-search.txt
```

---

## 6. Compare Working and Failing Services

If one map/service works and another fails, compare:

```text
Startup command
Game/client port
Server-to-server or IGW port
Bind address
Advertised public address
Map name
Partition/index/destination value
Runtime image or binary version
Mounted Saved/config path
Environment variables
Database host/name/user
RabbitMQ/messaging host and port
```

Run on: Docker host shell, only if Docker is confirmed

```bash
for c in $(docker ps --format '{{.Names}}'); do
  echo "===== $c ====="
  docker inspect "$c" --format 'Image={{.Config.Image}} NetworkMode={{.HostConfig.NetworkMode}} Args={{json .Args}}'
  docker inspect "$c" --format '{{range .Config.Env}}{{println .}}{{end}}' | sed -E 's/(TOKEN|SECRET|PASSWORD|KEY)=.*/\1=<redacted>/Ig'
done | tee docker-config-comparison.txt
```

---

## 7. Interpret Results

```text
Working and failing services use different image/build versions:
  Use update-patch-and-version-validation.md.

Failing service lacks a required port or bind argument:
  Use port-and-network-listener-validation.md.

Failing service points to a different Saved/config/database path:
  Verify active instance path and backup before changing anything.

Secrets or tokens are missing or empty:
  Escalate to the environment owner or authorized admin. Do not invent replacement secrets.

Generated config differs from source config:
  Identify which tool generates the runtime file before editing.
```

## 8. Evidence to Escalate

```text
Reason for config review
Running launch arguments with secrets redacted
Service/runtime configuration
Config file list
Working vs failing service comparison
Recent config changes
Exact UTC failure window
```
