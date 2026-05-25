# Frontend Mutation Confirmation

## Purpose

Frontend Mutation Confirmation provides a shared user-interface foundation for high-risk and destructive operator workflows. It replaces ad hoc browser prompts with a consistent preview, warning, reason, and confirmation pattern.

This is a foundation slice only. It does not add new mutations and does not wire Player 360 quick actions.

## Current scope

The first slice provides:

- a reusable frontend confirmation dialog component
- typed confirmation payload metadata
- risk-aware display for medium, high, and destructive actions
- operator warning display
- recommended-path display
- rollback-hint display
- optional or required reason capture
- explicit cancel/confirm behavior

## Out of scope

Do not include these in this slice:

- new Player 360 quick actions
- new backend mutation routes
- automatic replacement of all existing mutation calls
- arbitrary SQL execution changes
- bypasses for audit or reason enforcement

## Expected data source

The component is intended to consume the existing mutation-safety classification response:

```text
GET /api/v1/mutation-safety/classify?method=POST&path=/api/v1/players/give-item
```

It should display:

- action name
- risk
- preview requirement
- reason requirement
- destructive flag
- rollback hint
- operator warnings
- recommended path

## Security rules

- Do not use this component as authorization.
- Do not treat frontend confirmation as a replacement for backend validation.
- Do not log secrets or raw payloads in confirmation text.
- Require operator reason when the backend classification says a reason is required.
- Keep final mutation submission responsible for sending `X-Admin-Reason`.

## Implementation status

Initial component added:

```text
web/src/components/MutationConfirmationDialog.tsx
```

The component is intentionally reusable and not yet wired into existing workflows. Existing mutation behavior remains unchanged until a later validated integration slice.

## Validation

```bash
cd web
npm run typecheck
npm run lint
npm run build
```

## Follow-up tasks

1. Add a hook or service adapter that classifies a mutation and opens the dialog.
2. Replace browser `confirm`/`prompt` fallback in a controlled slice.
3. Wire one low-risk existing mutation workflow first.
4. Confirm audit reason headers are preserved.
5. Only after that, allow Player 360 quick actions to use this shared flow.
