# Project Continuity Notes

Use this document when continuing the Dune: Awakening self-hosted troubleshooting documentation in a new chat, support session, or handoff.

## Active Work Label

```text
Active label: OIC
Purpose: reusable operations / incident / troubleshooting documentation for Dune: Awakening self-hosted server support.
Repository: yacketrj/dune-admin-fork
Primary path: docs/dune-awakening-self-hosted-troubleshooting/
Style: discovery-first, platform-neutral, entry-level support friendly
```

## Boundary From Incident Work

```text
Client-specific incident records belong under docs/dune-awakening-incidents/.
Reusable troubleshooting documentation belongs under docs/dune-awakening-self-hosted-troubleshooting/.
Do not place client-specific evidence, environment names, personal data, account identifiers, passwords, tokens, or incident-only conclusions in the reusable troubleshooting guide.
Promote only generalized lessons, neutral decision logic, and reusable capture steps into OIC documentation.
```

## Documentation Rules

```text
Confirm the hosting platform before running platform-specific steps.
Confirm the runtime or orchestration layer before running runtime-specific steps.
Commands must state where they run.
Use operational variables for paths, services, containers, maps, ports, and IPs.
Redact tokens, secrets, passwords, private keys, and personal identifiers before broad sharing.
Keep incident-specific evidence separate from reusable troubleshooting guidance.
```

## Current OIC State

```text
Reusable guide now includes platform guides, runtime/orchestration guides, focused runbooks, a QA checklist, and generalized lessons learned from incident work.
Incident-specific current-investigation files have been removed from the reusable troubleshooting folder.
The generalized Dynamic Instancing and Handoff Validation runbook has been added to the reusable troubleshooting workflow.
README links reflect the current OIC workflow.
```

## Handoff Prompt

```text
Continue work on the OIC Dune: Awakening Self-Hosted Troubleshooting docs in repository yacketrj/dune-admin-fork under docs/dune-awakening-self-hosted-troubleshooting/. Keep the docs discovery-first, platform-neutral, and written for entry-level support. Every command must state where it runs. Keep README links accurate, preserve operational variables and redaction guidance, and keep client-specific incident evidence under docs/dune-awakening-incidents/ only.
```
