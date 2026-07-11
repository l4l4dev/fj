---
id: TASK-10.1
title: Version command and build metadata
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 09:08'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-9
dependencies:
  - TASK-10.5
parent_task_id: TASK-10
priority: medium
ordinal: 10030
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide fj version and define the build metadata and User-Agent version policy.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 fj version displays the approved version format.
- [x] #2 Build metadata has a documented default and injection path.
- [x] #3 User-Agent version policy is consistent with the displayed version.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision pending. Define fj version, build metadata, and User-Agent version policy after User Acceptance.

Model: Codex — implementation affects public CLI version output and User-Agent compatibility.

Independent Review:
- Critical: none
- Major: none; version/User-Agent injection path is consistent.
- Minor: ldflags real-build verification is deferred to TASK-11 Public Release Foundation.
- Suggestion: standardize release-time ldflags in TASK-11.

Verification:
- gofmt -l . passed
- git diff --check passed
- go vet ./... passed
- go test ./... passed
- make pre-commit passed.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented fj version and fj --version with version-only stdout output, default dev metadata, linker-injectable version, and consistent fj/<version> User-Agent propagation through the Composition Root. Existing CLI contracts remain unchanged. ldflags release-build verification is deferred to TASK-11.
<!-- SECTION:FINAL_SUMMARY:END -->
