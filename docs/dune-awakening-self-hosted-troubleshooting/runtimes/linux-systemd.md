# Runtime Guide: Linux systemd

Use this only after the Dune: Awakening server or related services are confirmed to be managed by Linux systemd.

## 1. Find Dune-Related Services

Run on: Linux host or Linux VM shell

```bash
systemctl list-units --type=service --all | grep -Ei 'dune|awakening|sandbox|rabbit|amp|docker|compose' || true
systemctl list-unit-files | grep -Ei 'dune|awakening|sandbox|rabbit|amp|docker|compose' || true
```

Record:

```text
SERVICE_NAME=
DIRECTOR_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_SERVICE=
```

## 2. Check Service Status

Run on: Linux host or Linux VM shell

```bash
systemctl status "$SERVICE_NAME" --no-pager
```

If multiple services exist, repeat the command for each one.

## 3. View Recent Logs

Run on: Linux host or Linux VM shell

```bash
journalctl -u "$SERVICE_NAME" --since "30 minutes ago" --no-pager
```

For live capture during a test:

```bash
journalctl -u "$SERVICE_NAME" -f | tee systemd-service-live.log
```

## 4. Identify the Launch Command

Run on: Linux host or Linux VM shell

```bash
systemctl cat "$SERVICE_NAME"
systemctl show "$SERVICE_NAME" -p ExecStart -p User -p WorkingDirectory -p Environment --no-pager
```

Record:

```text
Service user:
Working directory:
ExecStart command:
Environment file, if any:
```

## 5. Capture Process and Listener State

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep DuneSandbox | grep -v grep | sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g'
sudo ss -tulpen | grep -Ei 'Dune|777|778|779|780|781|788|789|790|791|792|31982|31983' || true
```

## 6. Stop Before Changing Anything

Do not restart or edit the service until logs, listener state, process arguments, and the user-defined failure window are captured.
