# Runbook: Permission and Ownership Errors

Use this runbook when startup, configuration generation, save writes, log writes, or file edits fail with permission or ownership errors.

Common error pattern:

```text
Permission denied
Access is denied
UnauthorizedAccessException
PermissionError: [Errno 13]
Failed to write file
Cannot create directory
```

## 1. Identify the File or Directory That Failed

Record the exact error and path.

```text
Error message:
Failed path:
Operation being attempted:
Process or service that attempted the write:
```

## 2. Check Linux File Ownership and Permissions

Run on: Linux host or Linux VM shell

```bash
ls -ld "$INSTANCE_PATH" "$SAVED_PATH" "$SAVED_PATH/UserSettings" 2>/dev/null || true
ls -la "$SAVED_PATH" 2>/dev/null || true
stat "$SAVED_PATH" 2>/dev/null || true
```

If the failed file is known:

```bash
stat "<FAILED_FILE_PATH>" 2>/dev/null || true
ls -l "<FAILED_FILE_PATH>" 2>/dev/null || true
```

## 3. Check Windows File Permissions

Run on: Windows host PowerShell

```powershell
Get-Item $env:INSTANCE_PATH | Format-List FullName, Attributes, Owner
Get-Acl $env:INSTANCE_PATH | Format-List
Get-Acl $env:SAVED_PATH | Format-List
```

If the failed file is known:

```powershell
Get-Item "<FAILED_FILE_PATH>" | Format-List FullName, Attributes, Length, LastWriteTime
Get-Acl "<FAILED_FILE_PATH>" | Format-List
```

## 4. Identify the Service or Container User

Run on: Linux host or Linux VM shell for host processes

```bash
ps -eo user,pid,ppid,cmd | grep -Ei 'DuneSandbox|AMP|Rabbit|docker' | grep -v grep
```

Run on: Linux systemd host if systemd is used

```bash
systemctl show "$SERVICE_NAME" -p User -p Group -p ExecStart -p WorkingDirectory --no-pager
```

Run on: Docker host shell if Docker is confirmed

```bash
docker inspect "$CONTAINER_NAME" --format 'User={{.Config.User}} Image={{.Config.Image}}'
IMG="$(docker inspect "$CONTAINER_NAME" --format '{{.Config.Image}}' 2>/dev/null || true)"
if [ -n "$IMG" ]; then
  docker run --rm --entrypoint sh "$IMG" -lc 'id; id dune 2>/dev/null || true'
fi
```

Run on: Windows host PowerShell if Windows service is used

```powershell
Get-CimInstance Win32_Service | Where-Object { $_.Name -eq $env:SERVICE_NAME } | Select-Object Name, StartName, PathName, State
```

## 5. Repair Pattern for Linux Bind Mounts or Host Paths

Use ACLs when both the host automation user and a container/service user need write access.

Run on: Linux host or Linux VM shell

```bash
sudo apt-get update
sudo apt-get install -y acl

HOST_USER="${SUDO_USER:-$(id -un)}"
CONTAINER_UID=<CONTAINER_UID_OR_SERVICE_UID>
CONTAINER_GID=<CONTAINER_GID_OR_SERVICE_GID>

sudo mkdir -p "$SAVED_PATH/UserSettings"
sudo setfacl -R -m "u:${HOST_USER}:rwx,u:${CONTAINER_UID}:rwx,g:${CONTAINER_GID}:rwx" "$SAVED_PATH"
sudo setfacl -R -d -m "u:${HOST_USER}:rwx,u:${CONTAINER_UID}:rwx,g:${CONTAINER_GID}:rwx" "$SAVED_PATH"
```

Do not repeatedly `chown` files back and forth unless the correct service user and host user are known.

## 6. Repair Pattern for Windows Paths

Run on: Windows host PowerShell as Administrator

```powershell
icacls $env:SAVED_PATH
```

After identifying the service account that needs access, grant only the needed permissions.

```powershell
icacls $env:SAVED_PATH /grant "<SERVICE_ACCOUNT>:(OI)(CI)M"
```

Use Modify (`M`) rather than Full Control unless Full Control is required.

## 7. Validate the Fix

Run on: Linux host or Linux VM shell

```bash
sudo -u "$HOST_USER" test -w "$SAVED_PATH" && echo "host user can write"
```

Run on: Windows host PowerShell

```powershell
Test-Path $env:SAVED_PATH
```

Then restart only the affected service or process according to the matching runtime guide, and capture logs during startup.
