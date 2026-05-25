# Runbook: Client-Reported Symptoms

Use this runbook when the only available information is what the player sees in the Dune: Awakening client.

Goal: turn a vague player report into a precise troubleshooting path.

## 1. Capture the Player's Exact Experience

Ask the player or environment owner:

```text
Can the player see the server in the server list?
Can the player select the server?
Can the player enter the password, if one is required?
Does the client hang, disconnect, crash, return to menu, or show an error?
What is the exact error text?
What was the player doing immediately before the issue?
What map, destination, or action was involved?
Does the issue happen to one player, some players, or all players?
Does restarting the client change anything?
Exact UTC time of the issue:
```

## 2. Categorize the Symptom

```text
Server not visible:
  Use server visibility and listing checks.

Server visible but player cannot connect:
  Use login and authentication failure checks.

Player connects but fails during loading:
  Use login and destination lifecycle checks.

Player is in-game but travel hangs or disconnects:
  Use map travel and instancing failure checks.

Player disconnects after a delay:
  Use logs, listener checks, packet capture, and messaging checks.

Client crashes:
  Capture client-side crash details and server logs for the same UTC time.
```

## 3. Ask for a Screenshot or Exact Text

Run on: player/client side

```text
Ask the player to provide a screenshot or exact text of the client error.
Record the local time and convert or compare it to UTC.
```

## 4. Correlate With Server Time

Run on: Linux host or Linux VM shell

```bash
date -u
timedatectl
```

Run on: Windows host PowerShell

```powershell
Get-Date -Format u
w32tm /query /status
```

## 5. Choose the Next Runbook

```text
Server not visible -> server-visibility-and-registration.md
Login/password/token problem -> login-and-authentication-failure.md
Travel or map hang -> map-travel-and-instancing-failure.md
Startup issue -> server-startup-failure.md
Network or port suspicion -> port-and-network-listener-validation.md
Unclear -> failed-travel-capture.md or log-collection-and-redaction.md
```

## 6. Evidence to Escalate

```text
Player symptom in exact words
Screenshot or exact error text
Client local time and UTC time
Player action before issue
Known working behavior
Known failing behavior
Server-side logs for same UTC window
```
