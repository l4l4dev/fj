---
id: TASK-11.4
title: Installation documentation
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 17:33'
updated_date: '2026-07-11 20:00'
labels: []
milestone: m-10
dependencies:
  - TASK-11.3
modified_files:
  - README.md
parent_task_id: TASK-11
priority: medium
ordinal: 11040
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Document how to obtain, verify, and install the initial fj release artifacts without adding Homebrew or system-wide distribution.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Installation instructions cover artifact retrieval and SHA-256 verification.
- [x] #2 Documentation covers version confirmation and uninstall considerations.
- [x] #3 Documentation uses placeholders and does not expose credentials or private information.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add a concise Install release artifacts section to README.md before the existing source-install guidance.
2. Document the supported darwin/arm64 and linux/amd64 targets and retrieval of the consolidated Release foundation workflow artifact.
3. Add platform-specific checksum verification, current-user installation, version/help confirmation, and limited uninstall procedures.
4. Clarify workflow-artifact, signing, notarization, and distribution limitations without duplicating existing configuration, Quickstart, or security guidance.
5. Rename the existing macOS section to Install from source on macOS, verify README-only scope and commands, then obtain independent review.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 (Codex) — user-facing installation safety, workflow/documentation consistency, command verification, and independent review required capable implementation and review reasoning.

Verification:
- README.md was the only implementation file changed.
- git diff --check passed.
- Markdown heading hierarchy and code-fence structure were checked.
- Workflow name, consolidated artifact name, binary names, checksum filename, normalized version, and supported targets match release-foundation.yml.
- macOS shasum -a 256 --check passed for both test binaries.
- sha256sum --check passed for both test binaries.
- Checksum tamper detection passed.
- Installation to a temporary HOME succeeded with executable permissions.
- The installed test binary reported version 1.2.3.
- Uninstall removed only the temporary binary and preserved the configuration fixture.
- Privacy and secret scan passed.
- make pre-commit passed, including go vet ./... and go test ./....
- Go module stat-cache writes emitted sandbox warnings, but cross-build and documentation command simulation succeeded.

Independent Review:
- Result: Review Ready.
- Critical: none.
- Major: none.
- Minor: none.
- Suggestion: none.
- Confirmed README-only scope, workflow naming and target consistency, safe checksum/install/uninstall commands, accurate security limitations, no duplicated Quickstart/configuration/security content, preserved source-install behavior, and separation from TASK-11.5 actual acceptance.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added a concise README guide for obtaining the consolidated Release foundation workflow artifact, verifying its darwin/arm64 and linux/amd64 binaries with the versioned SHA-256 manifest, installing fj for the current user, confirming the injected version and help output, and uninstalling only the installed binary. Clarified that public Releases, permanent distribution, signing, notarization, Homebrew, and system-wide installation are not yet provided. Existing source-install guidance remains available under a distinct heading, while configuration, Quickstart, and credential guidance are linked rather than duplicated. Verification and Independent Review completed successfully with no findings; actual release artifact acceptance remains TASK-11.5 responsibility.
<!-- SECTION:FINAL_SUMMARY:END -->
