# Runbook: Login and Authentication Failure

Use this runbook when players cannot log in, authentication fails, password/token validation fails, or the player reaches login but never enters the game world.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Capture the User-Visible Login Symptom

Ask the player or environment owner to report the exact behavior.

```text
Can the player see the server?
Can the player select the server?
Does the client request a password?
Does no password work?
Does a known password fail?
Does the client hang after submit?
Does the client return to menu?
Does the client show an error message?
Exact UTC time of login attempt:
```

## 2. Capture Control-Plane or Login Logs

Run on: Docker host shell, only if the login/control-plane service is a container

```bash
docker logs --since 30m "$DIRECTOR_SERVICE" 2>&1 | tee login-control-plane.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
grep -RniE 'LoginRequest|Login|Password|Token|Authentication|Authorization|FLS|Travel request|ServerAuthenticator|failed|error|exception' "$LOG_PATH" 2>/dev/null | head -300 | tee login-log-search.txt
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'LoginRequest','Login','Password','Token','Authentication','Authorization','FLS','Travel request','ServerAuthenticator','failed','error','exception' -ErrorAction SilentlyContinue | Select-Object -First 300 | Tee-Object -FilePath login-log-search.txt
```

Run in: control panel UI

```text
Open the control-plane, gateway, director, or login service console/log view.
Export or copy log output covering the login attempt.
```

## 3. Capture Destination or Starting Map Logs

Run on: Docker host shell, only if the starting map is a container

```bash
docker logs --since 30m "$DESTINATION_SERVICE" 2>&1 | tee login-destination.log
```

Run on: Linux host or Linux VM shell

```bash
grep -RniE 'PreLogin|Login|Welcome|VerifyFlsIdentity|VerifyFlsAuthorization|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|FinishSpawn|Disconnect|failed|error|exception' "$LOG_PATH" 2>/dev/null | head -300 | tee login-destination-search.txt
```

Run on: Windows host PowerShell

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'PreLogin','Login','Welcome','VerifyFlsIdentity','VerifyFlsAuthorization','DatabaseLogin','CharacterDownload','Join','GameModeLogin','StartingNewPlayer','FlsLogin','FinishSpawn','Disconnect','failed','error','exception' -ErrorAction SilentlyContinue | Select-Object -First 300 | Tee-Object -FilePath login-destination-search.txt
```

## 4. Check Listener and Traffic State

Run on: Linux host or Linux VM shell

```bash
sudo ss -tulpen | tee login-listeners.txt
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath login-udp-listeners.txt
Get-NetTCPConnection | Sort-Object LocalPort | Tee-Object -FilePath login-tcp-listeners.txt
```

If a specific client IP is known, capture one login attempt.

Run on: Linux host or Linux VM shell

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP}" | tee login-client-packets.log
```

Run on: Windows host PowerShell

```powershell
pktmon start --capture --pkt-size 0
```

After the login attempt:

```powershell
pktmon stop
pktmon format PktMon.etl -o login-pktmon-capture.txt
```

## 5. Interpret Results

```text
No login/control-plane logs appear:
  The request may not be reaching the login/control-plane service. Check server visibility, network path, and correct log source.

Login request appears but password/token validation fails:
  Check configured password/token source, secrets, environment variables, and whether the client is submitting the expected value.

Login request appears and succeeds, but starting map has no matching PreLogin or Join stages:
  Check travel handoff, destination availability, and listener state.

Starting map reaches PreLogin or VerifyFls stages but fails later:
  Check auth/session state, database login, character download, and disconnect reason.

Packets arrive but no server response leaves:
  Check service health, firewall, auth/session state, and process handling.

No packets arrive:
  Check cloud firewall, host firewall, NAT, Hyper-V/Proxmox bridge, Docker port publishing, and advertised public IP.
```

## 6. Evidence to Escalate

```text
User-visible login symptom
Exact UTC login attempt time
Control-plane/director/login logs
Starting map/destination logs
Listener state
Packet capture, if available
Runtime launch arguments or service config
Cloud/firewall/NAT evidence, if applicable
```
