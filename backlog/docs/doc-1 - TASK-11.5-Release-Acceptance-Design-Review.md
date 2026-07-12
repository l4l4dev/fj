---
id: doc-1
title: TASK-11.5 Release Acceptance Design Review
type: other
created_date: '2026-07-11 20:19'
tags:
  - review
  - release
  - m-10
---
# TASK-11.5 Release Acceptance — Public Release Foundation Design Review

- Date: 2026-07-12
- Reviewer model: Fable 5 (main session; review only, no implementation produced by reviewer)
- Scope: design review of the Public Release Foundation (TASK-11.1–11.4 outputs and TASK-11.5 acceptance criteria) for long-term operability. No implementation changes were made.
- Reviewed artifacts: `.github/workflows/release-foundation.yml`, `README.md` (install/verify/uninstall sections), TASK-11 / TASK-11.1–11.5 Backlog records, TASK-10.5 acceptance-exception record, `internal/version` wiring, `go.mod`.

## Verdict

The foundation is sound for long-term operation at its current, pre-publication scope. Artifact naming, checksum policy, and installation documentation are internally consistent and honestly scoped. No Critical findings. Three Major findings require human decisions, primarily as gates before any public GitHub Release.

## Confirmed strengths

- Workflow uses least privilege (`permissions: contents: read`), strict `vMAJOR.MINOR.PATCH` validation, `CGO_ENABLED=0` static builds, and version injection through the approved ldflags path (`internal/version.Value`); `fj version` and the client `User-Agent: fj/<version>` share `version.Current()`.
- Checksum manifest is deterministic (fixed order, `LC_ALL=C`), sha256sum-compatible, self-verified in-workflow, and guarded against unexpected/missing files.
- Naming `fj-<version>-<goos>-<goarch>` plus `fj-<version>-checksums.txt` is unambiguous and matches between workflow and README (verified in TASK-11.4).
- README states integrity-vs-authenticity limits explicitly, keeps credentials out of examples, and does not overclaim distribution status.
- TASK-11.5 AC #4 correctly carries forward the TASK-10.5 Option A decision: unsuccessful `repo list` / `pr list` are never treated as passed.

## Findings

### Critical

None.

### Major (human decision required)

- **MAJ-1 — No LICENSE file in the repository.** Distributing fj binaries via a public GitHub Release without a declared license is a publication blocker. License selection is a human decision.
- **MAJ-2 — The workflow has never executed on real GitHub Actions.** TASK-11.1–11.3 all record live execution and actionlint as unverified; TASK-11.5 acceptance would be the first real run, but no acceptance criterion states that a real `Release foundation` run must succeed and produce `fj-<version>-artifacts`. Decide whether to add this explicitly to TASK-11.5 AC.
- **MAJ-3 — `workflow_dispatch` breaks tag↔commit binding.** Any ref can be built with an arbitrary version string, so a binary's reported version is not guaranteed to correspond to a tag. Acceptable for the foundation, but before public Releases a human must decide: restrict publication to tag-triggered runs, or formally scope dispatch as non-releasable verification builds.

### Minor

- **MIN-1** Go toolchain version is duplicated (`go-version: '1.26.5'` in the workflow vs `go 1.26.5` in `go.mod`); drift risk. `go-version-file: go.mod` would remove it.
- **MIN-2** Builds are described as reproducible but omit `-trimpath`; host paths are embedded, so reproducibility is currently best-effort.
- **MIN-3** The checksums job hardcodes both binary names in three places; each new target requires coordinated edits (matrix, checksum list, upload list).
- **MIN-4** Actions are pinned by tag (`@v4`, `@v5`) rather than commit SHA; weaker supply-chain guarantee for a release pipeline.
- **MIN-5** TASK-11.5 AC #1 does not define the acceptance platform. Only macOS arm64 is available for hands-on acceptance; define what "installed and verified" means for linux/amd64 (e.g., checksum + `file` type check, or container execution) so AC #1 stays verifiable.
- **MIN-6** TASK-11.5 AC #3 "evaluated" is vague; it should reference the TASK-10.5 confirmed read-only set (`version`, `repo inspect`, `issue list/inspect`) plus redaction checks as the concrete acceptance content.

### Suggestions

- **SUG-1** For public Releases, ship tar.gz archives (bundling LICENSE) instead of raw binaries.
- **SUG-2** Add authenticity on top of integrity: GitHub artifact attestations (e.g., `actions/attest-build-provenance`) and/or signing (minisign/cosign); decide macOS Gatekeeper handling (signing/notarization vs documented quarantine workaround).
- **SUG-3** Add actionlint to CI.
- **SUG-4** Define release-notes/CHANGELOG and artifact-retention/support policies before first publication.

## Gaps before moving to public GitHub Releases

1. LICENSE (MAJ-1) — blocker.
2. A separate publication workflow: `contents: write`, tag protection, tag↔commit guarantee (MAJ-3).
3. Authenticity/provenance decision (SUG-2).
4. Archive format and bundled files decision (SUG-1).
5. Reproducibility hardening: `-trimpath`, `go-version-file`, optionally `SOURCE_DATE_EPOCH` (MIN-1/2).
6. Real-run verification of the pipeline plus actionlint (MAJ-2, SUG-3).
7. Release-notes and retention policy (SUG-4).

## Handling

Per AGENTS.md Section 15, no fixes were implemented. Each Major/Minor/Suggestion finding awaits a human adopt/defer/reject decision. Nothing was committed.
