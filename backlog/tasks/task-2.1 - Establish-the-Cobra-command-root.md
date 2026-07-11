---
id: TASK-2.1
title: Establish the Cobra command root
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-1
dependencies: []
references:
  - ROADMAP.md
modified_files:
  - go.mod
  - go.sum
  - cmd/fj/main.go
  - internal/interface/cli/root.go
  - internal/interface/cli/root_test.go
parent_task_id: TASK-2
ordinal: 15000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Provide the minimal command entry point and global help experience.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The CLI starts and displays valid root help
- [x] #2 Command construction respects ARCHITECTURE.md boundaries
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create the architecture-approved directories: cmd/fj, internal/domain, internal/application, internal/interface, internal/infrastructure, and docs.
2. Add cmd/fj/main.go as a minimal Go process entry point without business logic or configuration.
3. Add the minimal Cobra root command under internal/interface/cli and keep cmd/fj responsible only for composition and execution.
4. Add a focused test for the root help contract.
5. Run formatting, tests, root-help verification, and structure checks.
6. Check every acceptance criterion and finalize TASK-2.1 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Human instructions narrow this execution to the project skeleton only. Cobra initialization, CLI commands, root help, business logic, and configuration are explicitly out of scope.

Created the approved project skeleton directories: cmd/fj, internal/domain, internal/application, internal/interface, internal/infrastructure, and docs.
Added cmd/fj/main.go as the only new file. It contains only an empty Go main entry point.
No Cobra dependency, CLI commands, root help, business logic, or configuration was added.
Validation: gofmt completed; go test ./... passed after rerunning with build-cache permission; git diff --check passed.
The original acceptance criteria remain unchecked because root help and command construction are explicitly outside the human-approved scope for this execution.

Human approval received to continue TASK-2.1 beyond the completed skeleton and satisfy the remaining root-help and architecture-boundary acceptance criteria. TASK-2.2 remains out of scope.

Completed the approved continuation with Cobra v1.10.2.
The root command is constructed in internal/interface/cli; cmd/fj only composes and executes it, preserving the architecture boundary.
The root command displays help and defines no feature commands, configuration, or business logic.
Validation passed: gofmt completed, go test ./... passed, go run ./cmd/fj --help displayed the description, Usage, and help flag, and git diff --check passed.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Established the minimal Cobra root command and help experience while keeping Cobra in the Interface layer and the executable entry point limited to composition. Added a focused root-help test and verified formatting, all Go tests, live help output, and diff hygiene.
<!-- SECTION:FINAL_SUMMARY:END -->
