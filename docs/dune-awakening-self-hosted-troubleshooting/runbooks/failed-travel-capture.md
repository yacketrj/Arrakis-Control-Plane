# Runbook: Failed Travel Capture

Use this runbook after the hosting platform and runtime/orchestration layer have been identified.

Goal: capture one controlled failed travel attempt with logs, process/listener state, and traffic evidence.

## 1. Set Case Variables

Run on: Linux host or Linux VM shell

```bash
DIRECTOR_SERVICE=
DESTINATION_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_MAP=
CLIENT_IP=
INSTANCE_PATH=
LOG_PATH=
```

Run on: Windows host PowerShell

```powershell
$env:DIRECTOR_SERVICE=""
$env:DESTINATION_SERVICE=""
$env:RABBITMQ_SERVICE=""
$env:DESTINATION_MAP=""
$env:CLIENT_IP=""
$env:INSTANCE_PATH=""
$env:LOG_PATH=""
```

Record in case notes:

```text
UTC start time:
Starting location:
Destination/action being tested:
Expected behavior:
Observed behavior:
Approximate failure time:
Client-side error or behavior:
```

## 2. Create a Capture Directory

Run on: Linux host or Linux VM shell

```bash
mkdir -p ~/dune-travel-capture-$(date -u +%Y%m%d-%H%M%S)
cd ~/dune-travel-capture-*
pwd
```

Run on: Windows host PowerShell

```powershell
$CaptureDir = "$env:USERPROFILE\dune-travel-capture-$(Get-Date -Format yyyyMMdd-HHmmss)"
New-Item -ItemType Directory -Path $CaptureDir | Out-Null
Set-Location $CaptureDir
Get-Location
```

## 3. Capture Service or Container State

Run on: Linux host or Linux VM shell

```bash
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|amp|docker|compose|rabbit' | tee services-before.txt || true
ps -ef | grep -Ei 'DuneSandbox|Rabbit|AMP|docker' | grep -v grep | tee processes-before.txt || true
```

Run on: Windows host PowerShell

```powershell
Get-Service | Where-Object { $_.Name -match 'dune|amp|docker|rabbit' -or $_.DisplayName -match 'dune|amp|docker|rabbit' } | Tee-Object -FilePath services-before.txt
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Sandbox|Docker|AMP|Rabbit' } | Tee-Object -FilePath processes-before.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps | tee docker-ps-before.txt
```

## 4. Capture Control-Plane or Director Logs

Run on: Docker host shell, only if the director/control-plane is a container

```bash
docker logs -f "$DIRECTOR_SERVICE" 2>&1 | tee director-travel-capture.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
tail -F "$DIRECTOR_LOG_FILE" 2>&1 | tee director-travel-capture.log
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Get-Content -Path $env:DIRECTOR_LOG_FILE -Wait | Tee-Object -FilePath director-travel-capture.log
```

Run in: control panel UI

```text
Open the console or log viewer for the director/control-plane service.
Export or copy log output covering the test window.
```

## 5. Capture Destination Logs

Run on: Docker host shell, only if the destination is a container

```bash
docker logs -f "$DESTINATION_SERVICE" 2>&1 | tee destination-capture.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
find "$INSTANCE_PATH" -type f -name '*.log' 2>/dev/null | tee possible-log-files.txt
tail -F "$DESTINATION_LOG_FILE" 2>&1 | tee destination-capture.log
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Get-ChildItem -Path $env:INSTANCE_PATH -Recurse -Filter *.log -ErrorAction SilentlyContinue | Select-Object FullName | Tee-Object -FilePath possible-log-files.txt
Get-Content -Path $env:DESTINATION_LOG_FILE -Wait | Tee-Object -FilePath destination-capture.log
```

Run in: control panel UI

```text
Open the destination/server log viewer.
Export or copy log output covering the test window.
```

## 6. Capture Active Listeners

Run on: Linux host or Linux VM shell

```bash
sudo ss -uapn | tee udp-listeners-before.txt
while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  sudo ss -uapn || true
  sleep 2
done | tee udp-listeners-during-test.log
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-before.txt
while ($true) {
  "===== $(Get-Date -Format u) =====" | Tee-Object -FilePath udp-listeners-during-test.log -Append
  Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-during-test.log -Append
  Start-Sleep -Seconds 2
}
```

## 7. Capture Traffic

Run on: Linux host or Linux VM shell if `tcpdump` is installed and `CLIENT_IP` is known

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP}" | tee tcpdump-client-test.log
```

Run on: Linux host or Linux VM shell if `CLIENT_IP` is unknown

```bash
sudo tcpdump -ni any -vv 'udp or tcp' | tee tcpdump-server-test.log
```

Run on: Windows host PowerShell

```powershell
pktmon start --capture --pkt-size 0
```

Stop Windows packet capture after the test:

```powershell
pktmon stop
pktmon format PktMon.etl -o pktmon-test-capture.txt
```

## 8. Run One Controlled Reproduction

```text
1. Start from a known working location or action.
2. Perform the single failing travel/action once.
3. Do not retry repeatedly during the same capture.
4. Let the attempt succeed, hang, disconnect, or return.
5. Record the exact UTC start and failure time.
```

## 9. Collect Final State

Run on: Linux host or Linux VM shell

```bash
sudo ss -uapn | tee udp-listeners-after.txt
ps -ef | grep -Ei 'DuneSandbox|Rabbit|AMP|docker' | grep -v grep | tee processes-after.txt || true
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-after.txt
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Sandbox|Docker|AMP|Rabbit' } | Tee-Object -FilePath processes-after.txt
```

Run on: Docker host shell, only if RabbitMQ runs in Docker

```bash
docker exec "$RABBITMQ_SERVICE" rabbitmqctl list_queues name messages messages_ready messages_unacknowledged consumers 2>/dev/null | tee rabbitmq-queues-after.txt || true
```

## 10. Package the Capture

Run on: Linux host or Linux VM shell

```bash
cd ..
tar -czf dune-travel-capture-$(date -u +%Y%m%d-%H%M%S).tar.gz dune-travel-capture-*/
ls -lh dune-travel-capture-*.tar.gz
```

Run on: Windows host PowerShell

```powershell
Compress-Archive -Path $CaptureDir -DestinationPath "$CaptureDir.zip"
Get-Item "$CaptureDir.zip"
```
