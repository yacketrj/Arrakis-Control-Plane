# Route-Specific Audit Target Assertions

Date: 2026-06
Area: AppSec
Status: Validated partial remediation

## Summary

Expanded admin audit target metadata extraction and added route-specific audit target assertions for high-risk mutation requests.

This builds on the previous high-risk/destructive audit-event and blocked-mutation audit coverage by ensuring important target identifiers are captured in audit records, including player identity, item template, quantity, quality, command path, vehicle, guild, and rank fields.

## Commits

- `2731185f33b3395af84536f502ecd24829ad5ba7` — expanded audit target metadata extraction
- `c43c275281151752cfb82af230e1ed3c77dc602b` — added route-specific audit target assertions
- validation recorded after clean canonical local build/update path

## Validation

Validated from the canonical local update path:

```bash
./update.sh
```

Non-blocking build-performance warning observed:

```text
[PLUGIN_TIMINGS] Your build spent significant time in plugin `@tailwindcss/vite:generate:build`.
```

## Target metadata added or verified

- `player_id`
- `account_id`
- `actor_id`
- `controller_id`
- `fls_id`
- `item_id`
- `item_template`
- `item_template_id`
- `template_id`
- `quantity`
- `amount`
- `quality`
- `vehicle_id`
- `guild_id`
- `rank`
- `command`
- `command_path`

## Remaining work

- pre/post-change review verification where practical
- SAST/DAST/dependency evidence
- manual abuse-case validation

## Safety boundary

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
