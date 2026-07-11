---
id: TASK-11.3
title: Checksums and artifact management
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 17:33'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-10
dependencies:
  - TASK-11.2
modified_files:
  - .github/workflows/release-foundation.yml
parent_task_id: TASK-11
priority: medium
ordinal: 11030
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define SHA-256 checksum generation and artifact management for the initial fj build outputs.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 SHA-256 checksums are generated for each artifact.
- [x] #2 Checksum files map unambiguously to versioned artifacts.
- [x] #3 Artifact and checksum management avoids recording credentials or private information.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add a checksums job after resolve-version and the cross-platform build matrix.
2. Download and merge the two versioned binary artifacts, then validate the expected filenames and input count.
3. Generate fj-<version>-checksums.txt in deterministic sha256sum-compatible format and verify it with sha256sum --check.
4. Upload exactly the two binaries and checksum manifest as fj-<version>-artifacts.
5. Preserve read-only permissions and keep archives, signing, public releases, and installation documentation out of scope.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 (Codex) — checksum workflow semantics, artifact integrity validation, security boundaries, and independent review required capable implementation and review reasoning.

Verification:
- git diff --check passed.
- YAML parse passed with Ruby/Psych.
- darwin/arm64 and linux/amd64 test binaries were generated successfully.
- The two-entry versioned SHA-256 manifest was generated with lowercase 64-character digests, basenames only, and deterministic filename order.
- sha256sum --check passed for both binaries.
- Tampered binary detection passed.
- Missing binary detection passed.
- go test ./... passed.
- go vet ./... passed.
- gofmt -l . returned no files.
- make pre-commit passed.
- actionlint was unavailable and live GitHub Actions artifact download/merge/upload remains unverified.

Independent Review:
- Result: Review Ready.
- Critical: none.
- Major: none.
- Minor: none.
- Suggestion: none.
- Confirmed approved scope, all Acceptance Criteria, workflow semantics, unexpected file protection, checksum format, consolidated artifact composition, future Release asset compatibility, security/privacy, and TASK-11.2/TASK-11.4 boundaries.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added checksum and artifact aggregation to release-foundation.yml. The workflow downloads the two versioned binaries from TASK-11.2, validates the expected inputs, generates and verifies fj-<version>-checksums.txt in sha256sum-compatible format, and uploads exactly the two binaries plus the manifest as fj-<version>-artifacts. Verification covered successful checksums, deterministic manifest structure, tamper detection, missing-input detection, repository checks, and read-only security/privacy boundaries. Independent Review completed with no Critical, Major, Minor, or Suggestion findings. Archives, signing, public Release publication, and installation documentation remain out of scope.
<!-- SECTION:FINAL_SUMMARY:END -->
