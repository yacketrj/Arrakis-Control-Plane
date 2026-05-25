# Platform Guide: Linux Local Host or Linux VM

Use this when the Dune: Awakening server runs directly on Linux or inside a Linux VM. This guide does not assume Docker, AMP, systemd, or a specific cloud provider.

## 1. Confirm Linux Host Context

Run on: Linux host or Linux VM shell

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
whoami
pwd
```

Record:

```text
Hostname:
Linux distribution:
Kernel:
Virtualization result:
Logged-in user:
```

## 2. Find the Management Layer

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'DuneSandbox|amp|cubecoders|docker|containerd|podman|compose|systemd|wine' | grep -v grep
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|awakening|amp|docker|podman|compose' || true
```

Use the result to choose a runtime guide:

```text
AMP or control panel visible -> AMP control panel runtime guide.
Docker/container runtime visible -> Docker or Docker Compose runtime guide.
systemd service visible -> Linux systemd runtime guide.
Dune process visible without known service -> Manual or custom script runtime guide.
```

## 3. Find Instance and Log Paths

Run on: Linux host or Linux VM shell

```bash
find /home /opt /srv -maxdepth 8 -type f \( -name '*Engine.ini' -o -name 'UserEngine.ini' -o -name '*.log' \) 2>/dev/null | grep -Ei 'dune|awakening|DuneSandbox' | head -100
```

Record:

```text
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
```

## 4. Capture Process and Listener State

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep DuneSandbox | grep -v grep | sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g'
sudo ss -tulpen | grep -Ei 'Dune|777|778|779|780|781|788|789|790|791|792|31982|31983' || true
```

## 5. Next Runbooks

- Runtime layer found as AMP: [AMP control panel](../runtimes/amp-control-panel.md)
- Runtime layer found as Docker: [Docker or Docker Compose](../runtimes/docker-or-compose.md)
- Runtime layer found as systemd: [Linux systemd](../runtimes/linux-systemd.md)
- Unknown or custom launch: [Manual or custom script](../runtimes/manual-or-custom-script.md)
