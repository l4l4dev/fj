---
id: TASK-11.2
title: Cross-platform build
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 17:33'
updated_date: '2026-07-11 18:29'
labels: []
dependencies:
  - TASK-11.1
modified_files:
  - .github/workflows/release-foundation.yml
parent_task_id: TASK-11
priority: high
ordinal: 11020
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Generate reproducible initial fj artifacts for darwin/arm64 and linux/amd64, with additional targets deferred to later decisions.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 A darwin/arm64 artifact is generated.
- [x] #2 A linux/amd64 artifact is generated.
- [x] #3 Artifact names identify version, operating system, and architecture.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Split version resolution into a resolve-version job and expose normalized version/tag as job outputs.
2. Add an explicit darwin/arm64 and linux/amd64 build matrix.
3. Cross-compile with matrix GOOS/GOARCH and CGO_ENABLED=0 while preserving the approved ldflags version injection.
4. Name binaries and workflow artifacts fj-<version>-<goos>-<goarch>.
5. Verify both target builds, formats, version injection, repository checks, and TASK-11.3 scope boundary.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 (Codex) — workflow responsibility separation, cross-platform build verification, and independent review required capable implementation and review reasoning.

Verification:
- git diff --check passed.
- YAML parse passed with Ruby/Psych.
- darwin/arm64 cross build passed; file identified Mach-O 64-bit arm64; fj version returned 0.0.0-test.
- linux/amd64 cross build passed; file identified ELF 64-bit x86-64 static binary; injected version and ldflags were confirmed from the binary. Direct execution was unavailable on the macOS/arm64 host.
- go test ./... passed.
- go vet ./... passed.
- gofmt -l . returned no files.
- make pre-commit passed.
- actionlint was unavailable and live GitHub Actions matrix/artifact execution remains unverified.

Independent Review:
- Result: Review Ready.
- Critical: none.
- Major: none.
- Minor: none.
- Suggestion: none.
- Confirmed approved scope only, all Acceptance Criteria, architecture and security boundaries, version injection continuity, and separation from TASK-11.3 checksum responsibilities.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented explicit darwin/arm64 and linux/amd64 cross-platform builds in release-foundation.yml. Version resolution now runs once and passes the normalized version to matrix builds, which preserve the approved ldflags injection and upload target-specific binaries named fj-<version>-<goos>-<goarch>. Verified both cross builds, binary formats, version injection, Go tests, vet, formatting, diff checks, and make pre-commit. Independent Review completed with no Critical, Major, Minor, or Suggestion findings. Checksum generation, archive management, and release publication remain deferred to TASK-11.3 or later tasks.
<!-- SECTION:FINAL_SUMMARY:END -->
