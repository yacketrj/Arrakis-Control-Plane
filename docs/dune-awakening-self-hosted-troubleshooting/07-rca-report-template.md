# RCA Report Template

Use this template only after evidence has been collected. Do not write a final root cause before the working path, failing path, platform, runtime, and failure point are supported by logs or system evidence.

## 1. Executive Summary

```text
Issue summary:
Impact:
Affected environment:
Affected workflow:
Current status:
Most likely failure point:
Confidence level:
```

## 2. Scope

```text
Included systems:
Excluded systems:
Known evidence sources:
Known missing evidence:
```

## 3. Environment Confirmed by Evidence

```text
Hosting platform:
Runtime/orchestration layer:
Control panel:
Guest OS:
Container runtime, if applicable:
Cloud provider, if applicable:
Hypervisor, if applicable:
Public IP path:
Private IP/bind path:
Firewall/NAT/load-balancer path:
```

## 4. Incident Timeline

Use UTC.

```text
UTC time | Event | Evidence source
---------|-------|----------------
         |       |
```

## 5. Working Path

Document what works and why that matters.

```text
Known working workflow:
Evidence proving it works:
Systems involved:
Ports/listeners involved:
Logs proving success:
```

## 6. Failing Path

Document what fails and where the request stops.

```text
Known failing workflow:
Evidence proving failure:
Systems involved:
Ports/listeners involved:
Logs proving failure:
Last confirmed successful step:
First confirmed failed/missing step:
```

## 7. Evidence Analysis

```text
Evidence item:
What it proves:
What it does not prove:
Confidence:
```

## 8. Ruled-Out Causes

```text
Potential cause:
Evidence checked:
Conclusion:
```

## 9. Root Cause

Use this section only when evidence supports the conclusion.

```text
Root cause:
Evidence supporting root cause:
Why alternate causes were ruled out:
Confidence level:
```

If root cause is not proven, use:

```text
Root cause status: not yet proven
Most likely failure layer:
Next evidence required:
```

## 10. Corrective Actions

```text
Immediate action:
Owner:
Risk:
Rollback plan:
Validation step:
```

## 11. Preventive Actions

```text
Prevention item:
Owner:
Priority:
Due date:
Validation method:
```

## 12. Attachments

```text
[ ] Intake notes
[ ] Environment discovery
[ ] Platform guide output
[ ] Runtime guide output
[ ] Logs
[ ] Listener output
[ ] Packet capture
[ ] Cloud/firewall/NAT evidence
[ ] Screenshots
[ ] Change record
[ ] Backup record
```
