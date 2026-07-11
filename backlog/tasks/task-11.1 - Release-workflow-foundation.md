---
id: TASK-11.1
title: Release workflow foundation
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 17:33'
updated_date: '2026-07-11 18:00'
labels: []
dependencies:
  - TASK-10.1
parent_task_id: TASK-11
priority: high
ordinal: 11010
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define the foundation workflow for reproducible fj builds: read a vMAJOR.MINOR.PATCH tag, strip the v prefix, inject the normalized version through ldflags, and prepare artifacts without creating or publishing public releases.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The workflow maps a vMAJOR.MINOR.PATCH tag to a normalized version.
- [x] #2 The normalized version is passed through the approved ldflags injection path.
- [x] #3 The workflow remains a release foundation and does not create or publish a public release.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Verification:
- git diff --check passed
- gofmt -l . passed (no Go changes)
- go vet ./... passed
- go test ./... passed
- make pre-commit passed
- Workflow YAML manually reviewed; actionlint unavailable.

Independent Review:
- Critical: none
- Major: none
- Minor: actionlint and real GitHub Actions execution remain unverified.
- Suggestion: add actionlint and runner smoke validation in future CI hardening.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented the release workflow foundation with strict vMAJOR.MINOR.PATCH validation, workflow_dispatch version input, TASK-10.1 ldflags injection, fixed Go environment, minimal permissions, and a safe workflow artifact. Public release creation, checksum generation, cross-platform matrix, and secrets remain out of scope.
<!-- SECTION:FINAL_SUMMARY:END -->
