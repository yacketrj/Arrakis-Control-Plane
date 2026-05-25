# Release Packaging Strategy

## Goals

Dune Admin releases should be installable without requiring users to clone the repository or build from source.

The release system should produce a repeatable artifact set with clear release tags, release notes, checksums, and enough packaging context for operators to understand exactly what changed between releases.

The release system should eventually produce:

- Windows 11 installer for normal users.
- Windows portable ZIP for operators who prefer manual installs.
- Generic Linux archive for any distribution.
- Linux native packages for common package managers.
- macOS ZIP or tarball for both Intel and Apple Silicon.
- Checksums for every artifact.
- Release notes generated from release tags and curated changelog entries.

## Current release baseline

The repository uses GitHub Actions plus GoReleaser on semantic version tags:

```text
vMAJOR.MINOR.PATCH
```

The current baseline artifact set is:

```text
dune-admin_<version>_windows_amd64.zip
dune-admin_<version>_linux_amd64.tar.gz
dune-admin_<version>_linux_arm64.tar.gz
dune-admin_<version>_darwin_amd64.tar.gz
dune-admin_<version>_darwin_arm64.tar.gz
checksums.txt
```

The release workflow also builds the React/Vite frontend before GoReleaser runs so `web/dist` can be packaged with the release archives.

## Release tag model

Use release tags for every stable release:

```text
v0.1.0
v0.1.1
v0.2.0
v1.0.0
```

Recommended versioning rule:

| Version part | Use when |
|---|---|
| MAJOR | Breaking config, API, packaging, or operator workflow changes. |
| MINOR | New features, new tabs, new endpoints, new release artifacts, or substantial UX changes. |
| PATCH | Bug fixes, documentation corrections, small safety improvements, or packaging repair. |

Before the first public stable release, use `v0.x.y` tags. The first stable checkpoint for this repo should be a `v0.1.0` release after build/test validation passes and the release archive contents are verified.

## Changelog and release notes policy

Every feature or packaging change should update both:

```text
CHANGELOG.md
PATCH_NOTES.md
```

`CHANGELOG.md` should keep the long-running `[Unreleased]` section and should add a version section when a release tag is cut:

```markdown
## [v0.1.0] - YYYY-MM-DD
```

`PATCH_NOTES.md` should describe the current operator-facing update in plain language:

- why the change was made
- what changed
- operator/security impact
- validation steps
- known limitations or follow-up work

GoReleaser will generate release notes from commits, but curated docs remain the operator source of truth.

## Target artifact matrix

| Platform | Artifact | Status | Notes |
|---|---|---|---|
| Windows 11 | `.zip` portable package | Required now | Contains `dune-admin.exe`, docs, `.env.example`, item data, and frontend build output when available. |
| Windows 11 | `.msi` or `.exe` installer | Planned | Prefer WiX/MSI or NSIS after install location, shortcuts, service behavior, and uninstall behavior are decided. |
| Linux generic | `.tar.gz` archive | Required now | Works everywhere. |
| Linux Debian/Ubuntu | `.deb` | Planned/desired | Good fit for Ubuntu server operators. |
| Linux RHEL/Fedora | `.rpm` | Planned/desired | Standard Linux package format. |
| Linux Alpine | `.apk` | Optional | Useful for lightweight deployments. |
| macOS Intel | `.tar.gz` or `.zip` | Required now | May require signing/notarization later. |
| macOS Apple Silicon | `.tar.gz` or `.zip` | Required now | May require signing/notarization later. |
| Homebrew | formula/tap | Future | Requires tap repository and maintenance process. |
| Winget | manifest | Future | Best after Windows installer is stable. |
| Scoop | manifest/bucket | Future | Good lightweight Windows install path. |

## Release process

1. Confirm repo is clean.
2. Pull latest `main` and run local validation:

```powershell
.\update.ps1
```

3. Confirm backend starts and frontend can authenticate locally.
4. Confirm Docker/Kubernetes runtime detection still reports only:

```text
docker
kubernetes
```

5. Update `PATCH_NOTES.md` and `CHANGELOG.md`.
6. Pick the next semantic version tag.
7. Tag the release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

8. GitHub Actions runs the release workflow.
9. Download artifacts from the GitHub release and test at least:

```text
Windows 11 portable ZIP
Ubuntu generic tar.gz
macOS archive, when available
checksums.txt verification
```

10. Add any release-specific notes to the GitHub release if manual operator warnings are needed.

## Packaging requirements

Every release package should include:

```text
dune-admin / dune-admin.exe
.env.example
README.md
CHANGELOG.md
PATCH_NOTES.md
item-data.json
web/dist/ when the frontend build is present
docs/ release and setup notes
```

Generated packages must not include:

```text
.env
SSH keys
database dumps
logs
admin-audit.jsonl
node_modules
local dist folders outside release output
```

## Windows installer plan

Windows installer work should be a separate implementation slice.

Recommended path:

1. Keep ZIP portable releases working first.
2. Decide install behavior:
   - per-user install or machine-wide install
   - desktop/start-menu shortcuts
   - whether to run backend as a service or manual app
   - where `.env` lives
   - where audit logs live
3. Add WiX/MSI or NSIS packaging.
4. Add code-signing support when certificates are available.
5. Add Winget or Scoop after installer behavior is stable.

## Linux package plan

Linux native packaging should use GoReleaser `nfpms` after archive releases are stable.

Recommended targets:

```text
deb
rpm
apk
```

Install layout should be:

```text
/usr/bin/dune-admin
/usr/share/dune-admin/item-data.json
/usr/share/doc/dune-admin/README.md
/etc/dune-admin/dune-admin.env.example
```

Systemd service packaging should wait until configuration location and security defaults are finalized.

## macOS plan

Initial macOS support should be archive-based.

Later work:

- sign binaries
- notarize archives or app bundle
- consider Homebrew tap
- document local execution and quarantine handling

## Open decisions

- Should the production binary serve the built frontend from `web/dist`, or should frontend and backend remain separate?
- Should Windows install as a background service, tray app, or portable CLI/server executable?
- Should Linux packages install a systemd service by default or only provide a helper unit file?
- Where should audit logs live for packaged installs?
- Should release packages include sample scripts for setup and update?
- Should `ADMIN_REQUIRE_REASON=true` become the default after all high-risk UI paths are fully wired?

## Immediate next implementation steps

1. Verify release workflow builds frontend before GoReleaser runs.
2. Verify release archives include frontend build output.
3. Add Linux native package definitions.
4. Keep Windows portable ZIP as the first Windows release artifact.
5. Add Windows installer work as a future roadmap item, not part of the immediate mutation-safety work.
6. Cut the first `v0.1.0` release only after a clean local update/build and one artifact smoke test.
