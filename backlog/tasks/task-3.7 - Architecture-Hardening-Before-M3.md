---
id: TASK-3.7
title: Architecture Hardening Before M3
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 00:25'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-2
dependencies:
  - TASK-3.1
  - TASK-3.2
  - TASK-3.3
  - TASK-3.4
  - TASK-3.5
  - TASK-3.6
parent_task_id: TASK-3
priority: high
ordinal: 82000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
## Goal
Harden the architecture identified during the M2 completion review before M3 issue workflows begin.

## Scope
- Define an application-owned error classification boundary so HTTP status interpretation does not leak into the CLI layer.
- Consolidate duplicated command error mapping behind a consistent application/interface contract.
- Define a composition boundary that avoids runtime type assertions for repository capabilities.
- Establish a consistent validation error representation and mapping policy.
- Define the presenter/output boundary needed to preserve human-readable output while preparing for future JSON support.
- Record compatibility and migration expectations for the existing M1/M2 CLI contracts.

## Non-goals
- Implementing M3 issue workflows.
- Adding JSON output or changing the M7 machine-interface contract.
- Adding new Forgejo capabilities or changing remote API coverage.
- Introducing optional abstractions or framework dependencies without an approved design.
- Performance optimization unrelated to the identified architecture hardening.

## Validation
- Verify existing M1/M2 commands remain behaviorally compatible.
- Verify error categories, safe messages, and exit codes remain consistent.
- Verify composition and validation boundaries are covered by focused tests.
- Run the repository verification workflow after implementation.

## Review
- Perform an architecture and compatibility review before implementation.
- Perform an independent post-implementation review.
- Classify findings as Critical, Major, Minor, or Suggestion.
- Treat presenter abstractions and other future improvements as optional unless explicitly approved in the implementation scope.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 HTTP status classification is owned by the application boundary and the CLI receives safe categories without interpreting transport status codes.
- [x] #2 Duplicated command error mapping is consolidated without changing existing public error categories, messages, or exit codes.
- [x] #3 Composition Root wiring uses explicit capability dependencies and does not rely on runtime type assertions for repository operations.
- [x] #4 Validation failures use one consistent typed representation and preserve existing validation behavior.
- [x] #5 A presenter/output boundary is defined or introduced without changing existing human-readable output; future JSON support remains out of scope.
- [x] #6 Existing M1/M2 CLI contracts, credential safety, and repository operation behavior remain compatible and are covered by focused regression tests.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Review the current M1/M2 boundaries and contracts, approve a minimal hardening design, implement only the approved boundary changes, and validate compatibility before M3.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Created as M3 preparation work following the M2 completion review. Required hardening is separated from optional future presenter and JSON improvements.

Split into three implementation children before M3: TASK-3.7.1 Application Error Boundary Refactoring, TASK-3.7.2 Explicit Composition Root Refactoring (depends on 3.7.1), and TASK-3.7.3 Presenter Boundary Introduction (depends on 3.7.2). The parent remains the M3 preparation hardening gate; all children are intentionally To Do.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
M3 preparation Architecture Hardening is complete. Application Error Boundary Refactoring completed with semantic application errors, consolidated CLI mapping, and unified validation. Explicit Composition Root Refactoring completed with RepositoryDependencies and explicit capability injection while retaining the migration wrapper. Presenter Boundary Introduction completed by moving repository human-readable output into internal/interface/cli without changing output compatibility. M3 may begin after the established hardening contracts are used by subsequent workflows.
<!-- SECTION:FINAL_SUMMARY:END -->
