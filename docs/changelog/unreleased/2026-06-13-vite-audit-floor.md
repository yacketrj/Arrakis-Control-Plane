# Vite audit floor hardening

## Summary

Raised the direct Vite version floor in `web/package.json` after local audit reported high-severity Vite advisories affecting versions below the fixed lockfile version.

## Changes

- Updated the direct Vite dev dependency from `^8.0.12` to `^8.0.16`.
- Confirmed the repository lockfile already resolves `vite` to `8.0.16`.
- Confirmed the repository lockfile already resolves `@babel/core` to `7.29.7`.

## Operator note

If local audit still reports the old vulnerable versions, refresh the local install state before re-running the full update path.

Recommended local repair:

```bash
cd web
npm install
npm audit --audit-level=high
```

For a full clean dependency refresh if stale `node_modules` persists:

```bash
cd web
rm -rf node_modules
npm install
npm audit --audit-level=high
```

On Windows PowerShell:

```powershell
cd web
Remove-Item -Recurse -Force node_modules
npm install
npm audit --audit-level=high
```

## Validation

Pending local validation after dependency refresh:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```
