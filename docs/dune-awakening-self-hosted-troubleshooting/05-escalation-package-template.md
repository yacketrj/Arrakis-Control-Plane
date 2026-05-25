# Escalation Package Template

Use this template when preparing evidence for a developer, vendor, hosting provider, or senior escalation team.

Do not include assumptions as facts. Mark unknown values as `unknown`.

## 1. Summary

```text
Reported issue:
Impact:
First known failure time in UTC:
Current status:
Requested assistance:
```

## 2. Environment

```text
Hosting platform:
Runtime/orchestration layer:
Operating system observed from shell:
Operating system observed from game logs:
Control panel, if any:
Container runtime, if any:
Cloud provider, if any:
Hypervisor, if any:
Public IP handling:
Private IP handling:
Firewall/NAT/load-balancer path:
```

## 3. Working and Failing Paths

```text
Known working workflow:
Known failing workflow:
Source map/service, if relevant:
Destination map/service, if relevant:
Expected result:
Actual result:
Reproduction frequency:
```

## 4. Timeline

```text
UTC time | Event | Evidence file/source
---------|-------|---------------------
         |       |
```

## 5. Evidence Files Attached

```text
[ ] Intake notes
[ ] Environment discovery output
[ ] Platform guide output
[ ] Runtime/orchestration guide output
[ ] Control-plane logs
[ ] Source/destination logs
[ ] Process state before and after
[ ] Listener state before/during/after
[ ] Packet capture
[ ] Cloud/firewall/NAT evidence
[ ] Database/messaging evidence, if applicable
[ ] Change/backup record, if applicable
```

## 6. Key Log Excerpts

Paste short, relevant excerpts only. Keep enough surrounding lines to show context.

```text
<timestamp> <service> <log excerpt>
```

## 7. What Has Been Ruled Out

```text
Item checked:
Evidence:
Conclusion:
```

## 8. Open Questions

```text
Question:
Why it matters:
Evidence needed:
```

## 9. Redaction Confirmation

```text
[ ] Tokens removed
[ ] Passwords removed
[ ] Secrets removed
[ ] Personal names removed where needed
[ ] Player/account IDs redacted where needed
[ ] Cloud resource IDs redacted where needed
[ ] Operational values preserved when needed for troubleshooting
```
