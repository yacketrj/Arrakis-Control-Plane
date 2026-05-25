# Documentation Maintenance and QA Checklist

Use this checklist before calling the troubleshooting documentation complete or before handing it to a support team.

## 1. Link and File QA

Check that every link in `README.md` points to an existing file.

```text
[ ] Core workflow links work.
[ ] Platform guide links work.
[ ] Runtime/orchestration guide links work.
[ ] Focused runbook links work.
[ ] Old folder README redirects to the current folder.
[ ] Duplicate files are either removed or converted to clear redirect files.
```

Known redirect files are acceptable when they prevent broken historical links.

## 2. Platform-Neutrality QA

For every guide and runbook, verify:

```text
[ ] The document does not assume Docker unless Docker has been discovered.
[ ] The document does not assume AMP unless AMP has been discovered.
[ ] The document does not assume Hyper-V, Proxmox, OCI, AWS, Azure, or GCP unless discovered.
[ ] The document tells the reader which platform path to use next.
[ ] The document separates host, hypervisor, guest VM, container, and cloud edge evidence.
```

## 3. Entry-Level Support QA

For every command block, verify:

```text
[ ] The command block says where it runs.
[ ] The command block is followed by what to record or how to interpret the result.
[ ] The document avoids unexplained acronyms where possible.
[ ] The document tells the reader what to do if the command fails.
[ ] The document tells the reader when to stop and escalate.
```

Examples of command location labels:

```text
Run on: Linux host or Linux VM shell
Run on: Windows host PowerShell
Run on: Hyper-V host PowerShell
Run on: Proxmox host shell
Run on: Docker host shell, only if Docker is confirmed
Run inside: container shell
Run in: control panel UI
Run in: cloud provider console
```

## 4. Evidence Quality QA

For every runbook, verify:

```text
[ ] It captures the user-defined symptom first.
[ ] It records exact UTC time.
[ ] It records known working behavior.
[ ] It records known failing behavior.
[ ] It captures logs from the correct component.
[ ] It captures process/service/container state where applicable.
[ ] It captures listener/port state where applicable.
[ ] It captures packet data where applicable.
[ ] It includes escalation evidence requirements.
```

## 5. Redaction QA

Verify that redaction guidance is consistent:

```text
[ ] Redact tokens.
[ ] Redact passwords.
[ ] Redact secrets.
[ ] Redact private keys.
[ ] Redact personal names before broad sharing.
[ ] Redact raw account/player IDs unless needed and approved.
[ ] Do not automatically redact operational values needed for internal troubleshooting.
```

Operational values that may be retained internally include:

```text
Instance paths
Log paths
Service names
Container names
Map names
Destination names
Port numbers
Cloud instance IDs, when shared with trusted support staff
```

## 6. RCA QA

Before an RCA is marked final, verify:

```text
[ ] Hosting platform is confirmed by evidence.
[ ] Runtime/orchestration layer is confirmed by evidence.
[ ] Working path is confirmed.
[ ] Failing path is confirmed.
[ ] First missing or failed step is identified.
[ ] Competing explanations were considered.
[ ] Unknowns are labeled as unknown.
[ ] Root cause is not overstated beyond the evidence.
```

## 7. Future Improvement Backlog

Potential additions:

```text
[ ] Add provider-specific screenshots or UI navigation paths.
[ ] Add printable PDF and DOCX output after Markdown stabilizes.
[ ] Add a quick-start one-page checklist for live support calls.
[ ] Add a sample completed escalation package using fake values.
[ ] Add command-output examples for common good/bad results.
[ ] Add diagrams for common deployment patterns.
```
