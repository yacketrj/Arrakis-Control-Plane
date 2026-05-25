# Runtime Guide: Docker or Docker Compose

Use this only after Docker, Docker Compose, or a containerized deployment has been discovered.

## 1. Confirm Docker Access

Run on: Docker host shell, Linux or Windows PowerShell

```bash
docker version
docker ps
```

If these fail, record the exact error and return to the platform guide.

## 2. List Containers

Run on: Docker host shell

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
```

Record:

```text
CONTAINER_NAME=
DIRECTOR_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_SERVICE=
```

## 3. Map Mounts, Network Mode, and Ports

Run on: Docker host shell

```bash
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format 'NetworkMode={{.HostConfig.NetworkMode}} PortBindings={{json .HostConfig.PortBindings}}'
  docker inspect "$c" --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
done
```

## 4. Capture Logs

Run on: Docker host shell

```bash
docker logs --since 30m "$CONTAINER_NAME" 2>&1 | head -300
```

For live capture:

```bash
docker logs -f "$CONTAINER_NAME" 2>&1 | tee container-live-capture.log
```

## 5. Continue to Focused Runbooks

- [Failed travel capture](../runbooks/failed-travel-capture.md)
- [Port and listener validation](../runbooks/port-and-network-listener-validation.md)
- [Log collection and redaction](../runbooks/log-collection-and-redaction.md)
