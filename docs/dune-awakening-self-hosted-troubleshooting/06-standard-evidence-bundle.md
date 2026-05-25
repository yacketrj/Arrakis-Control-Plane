# Standard Evidence Bundle

Use this document when support needs a consistent evidence package for any Dune: Awakening self-hosted issue.

This is not a replacement for the focused runbooks. It is the minimum baseline bundle that should exist before escalation.

## 1. Required Case Notes

```text
Reported issue:
Known working behavior:
Known failing behavior:
Hosting platform:
Runtime/orchestration layer:
Control panel, if any:
Exact UTC failure window:
Recent change:
Evidence confidence level:
```

## 2. Required Environment Values

```text
CLOUD_PROVIDER=
CLOUD_INSTANCE_ID=
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
SERVICE_NAME=
CONTAINER_NAME=
DIRECTOR_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_SERVICE=
DESTINATION_MAP=
CLIENT_IP=
PUBLIC_IP=
PRIVATE_IP=
```

Write `unknown` for values not yet discovered.

## 3. Linux Baseline Evidence

Run on: Linux host or Linux VM shell

```bash
mkdir -p evidence-bundle

hostnamectl | tee evidence-bundle/linux-hostnamectl.txt
cat /etc/os-release | tee evidence-bundle/linux-os-release.txt
uname -a | tee evidence-bundle/linux-uname.txt
systemd-detect-virt -v 2>&1 | tee evidence-bundle/linux-virt.txt || true
date -u | tee evidence-bundle/linux-date-utc.txt
timedatectl | tee evidence-bundle/linux-timedatectl.txt

ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox|AMP|Rabbit|docker|compose|postgres' | grep -v grep | sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig' | tee evidence-bundle/linux-processes.txt || true
sudo ss -tulpen | tee evidence-bundle/linux-listeners.txt
free -h | tee evidence-bundle/linux-memory.txt
df -h | tee evidence-bundle/linux-disk.txt
```

## 4. Windows Baseline Evidence

Run on: Windows host PowerShell

```powershell
New-Item -ItemType Directory -Path evidence-bundle -Force | Out-Null

Get-ComputerInfo | Out-File evidence-bundle\windows-computerinfo.txt
Get-Date -Format u | Out-File evidence-bundle\windows-date-utc.txt
w32tm /query /status | Out-File evidence-bundle\windows-time-status.txt

Get-Service | Where-Object { $_.Name -match 'dune|awakening|sandbox|amp|docker|rabbit|postgres' -or $_.DisplayName -match 'dune|awakening|sandbox|amp|docker|rabbit|postgres' } | Out-File evidence-bundle\windows-services.txt
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Awakening|Sandbox|AMP|Rabbit|Docker|postgres' } | Select-Object ProcessName, Id, Path | Out-File evidence-bundle\windows-processes.txt
Get-NetUDPEndpoint | Sort-Object LocalPort | Out-File evidence-bundle\windows-udp-listeners.txt
Get-NetTCPConnection | Sort-Object LocalPort | Out-File evidence-bundle\windows-tcp-listeners.txt
Get-Volume | Out-File evidence-bundle\windows-volumes.txt
```

## 5. Docker Baseline Evidence

Run on: Docker host shell, only if Docker is confirmed

```bash
mkdir -p evidence-bundle

docker version > evidence-bundle/docker-version.txt 2>&1
docker ps -a --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' > evidence-bundle/docker-ps.txt

docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format 'Image={{.Config.Image}} User={{.Config.User}} NetworkMode={{.HostConfig.NetworkMode}} Status={{.State.Status}} ExitCode={{.State.ExitCode}} OOMKilled={{.State.OOMKilled}}'
  docker inspect "$c" --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
 done > evidence-bundle/docker-inspect-summary.txt
```

## 6. Control Panel Evidence

Run in: control panel UI

```text
Export or screenshot:
- Instance status
- Console output
- Recent restart/task history
- Visible startup options
- Install path
- Log path
- Visible ports
```

## 7. Package Baseline Evidence

Run on: Linux host or Linux VM shell

```bash
tar -czf dune-evidence-bundle-$(date -u +%Y%m%d-%H%M%S).tar.gz evidence-bundle/
ls -lh dune-evidence-bundle-*.tar.gz
```

Run on: Windows host PowerShell

```powershell
Compress-Archive -Path evidence-bundle -DestinationPath "dune-evidence-bundle-$(Get-Date -Format yyyyMMdd-HHmmss).zip"
Get-ChildItem dune-evidence-bundle-*.zip
```

## 8. Redaction Reminder

Redact tokens, passwords, secrets, raw account IDs, and personal identifiers before broad sharing. Do not remove operational values required for troubleshooting unless the sharing audience requires it.
