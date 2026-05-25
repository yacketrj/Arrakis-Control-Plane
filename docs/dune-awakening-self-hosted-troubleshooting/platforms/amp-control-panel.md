# Platform Guide: AMP-Controlled Hosting

Use this when the environment owner says the server is hosted or managed primarily through AMP, or when AMP is the only interface available to support staff.

AMP can be both the visible hosting surface and the management layer. It may launch a direct process, a Docker container, or another runtime behind the scenes. Do not assume which mode is used until verified.

## 1. Confirm AMP Access

Run in: AMP web UI

```text
Open the Dune: Awakening instance.
Record the instance name, instance status, module type, install path, configuration path, log path, console output location, and start/stop/restart controls.
```

Record:

```text
AMP panel URL or environment name:
AMP instance name:
Instance status:
Install path shown by AMP:
Log path shown by AMP:
Startup command or launch method, if visible:
Docker/container mode shown: yes/no/unknown
```

## 2. Confirm Whether Shell Access Exists

Run in: AMP UI or hosting provider UI

```text
Check whether the panel provides a file manager, terminal, console, log viewer, or container shell.
```

If shell access exists, continue with the matching platform guide:

```text
Linux shell available -> Linux local or Linux VM guide.
Windows shell available -> Windows / Hyper-V guide.
Docker/container shell visible -> Docker runtime guide.
```

## 3. Gather AMP Evidence

Run in: AMP web UI

```text
Export or copy console output covering the failure window.
Capture instance status before and after a test.
Record recent restart history.
Record visible ports and startup options.
Record any file permission, startup, or task/job errors.
```

## 4. Continue to the Runtime Guide

After identifying how AMP launches the server, continue with the matching runtime guide:

- [AMP control panel runtime](../runtimes/amp-control-panel.md)
- [Docker or Docker Compose](../runtimes/docker-or-compose.md)
- [Linux systemd](../runtimes/linux-systemd.md)
- [Windows service](../runtimes/windows-service.md)
- [Manual or custom script](../runtimes/manual-or-custom-script.md)
