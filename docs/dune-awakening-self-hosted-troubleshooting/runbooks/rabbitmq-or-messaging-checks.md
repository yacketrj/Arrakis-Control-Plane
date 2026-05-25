# Runbook: RabbitMQ or Messaging Checks

Use this runbook only after RabbitMQ or another messaging component is discovered or suspected.

Goal: determine whether messaging connections, queues, consumers, or stale per-player/service queues correlate with the reported failure.

## 1. Identify the Messaging Service

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'rabbit|beam|epmd|mq|amqp' | grep -v grep
systemctl list-units --type=service 2>/dev/null | grep -Ei 'rabbit|mq|amqp' || true
```

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'rabbit|beam|epmd|mq|amqp' } | Select-Object ProcessName, Id, Path
Get-Service | Where-Object { $_.Name -match 'rabbit|mq|amqp' -or $_.DisplayName -match 'rabbit|mq|amqp' } | Select-Object Name, DisplayName, Status
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | grep -Ei 'rabbit|mq|amqp' || true
```

Record:

```text
RABBITMQ_SERVICE=
Messaging runtime:
Messaging host/IP:
Messaging ports:
```

## 2. List Queues

Run on: Docker host shell, if RabbitMQ is in Docker

```bash
docker exec "$RABBITMQ_SERVICE" rabbitmqctl list_queues -p / name durable auto_delete arguments consumers messages messages_ready messages_unacknowledged state
```

Run on: Linux host or Linux VM shell, if RabbitMQ is installed directly

```bash
sudo rabbitmqctl list_queues -p / name durable auto_delete arguments consumers messages messages_ready messages_unacknowledged state
```

## 3. List Consumers and Connections

Run on: Docker host shell, if RabbitMQ is in Docker

```bash
docker exec "$RABBITMQ_SERVICE" rabbitmqctl list_consumers -p /
docker exec "$RABBITMQ_SERVICE" rabbitmqctl list_connections pid user peer_host peer_port state name
```

Run on: Linux host or Linux VM shell, if RabbitMQ is installed directly

```bash
sudo rabbitmqctl list_consumers -p /
sudo rabbitmqctl list_connections pid user peer_host peer_port state name
```

## 4. Check Logs

Run on: Docker host shell, if RabbitMQ is in Docker

```bash
docker logs --since 30m "$RABBITMQ_SERVICE" 2>&1 | tee rabbitmq-recent.log
```

Run on: Linux host or Linux VM shell, if RabbitMQ is installed directly

```bash
journalctl -u rabbitmq-server --since "30 minutes ago" --no-pager | tee rabbitmq-recent.log
```

Run on: Windows host PowerShell, if RabbitMQ is installed as a Windows service

```powershell
Get-WinEvent -LogName Application -MaxEvents 300 | Where-Object { $_.Message -match 'rabbit|amqp|mq' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message
```

## 5. Interpret Results

```text
Queue has many messages_ready:
  Messages are waiting and may not be consumed.

Queue has many messages_unacknowledged:
  A consumer may have received messages but is not acknowledging them.

Queue has active consumers during cleanup/deletion errors:
  Do not delete the queue while it is still in use.

Repeated client reconnects are visible:
  Correlate reconnect times with the reported gameplay failure window.

No RabbitMQ/messaging component is present:
  Return to environment discovery and runtime-specific logs.
```
