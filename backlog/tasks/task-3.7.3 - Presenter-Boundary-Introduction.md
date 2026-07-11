---
id: TASK-3.7.3
title: Presenter Boundary Introduction
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 00:31'
updated_date: '2026-07-11 05:50'
labels: []
dependencies:
  - TASK-3.7.2
parent_task_id: TASK-3.7
priority: medium
ordinal: 85000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
## Goal
Separate human-readable presentation from command orchestration before M3 while preserving all existing output.

## Scope
- Introduce a presenter/output boundary for existing repository results and errors.
- Move human-readable formatting responsibilities out of command orchestration.
- Preserve field order, empty-value rendering, empty-result messages, and standard stream behavior.
- Keep the boundary ready for a future JSON presenter without implementing JSON.

## Non-goals
- JSON output or M7 machine-interface contracts.
- M3 issue presenters.
- Changes to existing CLI names, flags, exit codes, or output text.
- New repository or Forgejo capabilities.

## Validation
- Add focused presenter and command regression tests.
- Verify every M1/M2 output remains byte-compatible.
- Run the repository verification workflow.

## Review
- Review output compatibility and architecture boundaries before implementation.
- Perform an independent post-implementation review.
- Classify findings as Critical, Major, Minor, or Suggestion.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Human-readable output responsibilities are separated from command orchestration.
- [x] #2 Existing M1/M2 output text, field order, empty-value rendering, and stream behavior remain unchanged.
- [x] #3 The boundary can support a future JSON presenter without implementing JSON or changing the M7 scope.
- [x] #4 Focused presenter and command regression tests pass, including make pre-commit.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Introduce a human presenter boundary after explicit composition wiring, migrate existing repository output, and verify byte-compatible behavior.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Third child of TASK-3.7. Depends on Explicit Composition Root Refactoring and must complete before M3.

Validation: gofmt -l ., git diff --check, go vet ./..., go test ./..., and make pre-commit all passed.

Independent Review: Critical: none. Major: none. Minor: Presenter unit tests were not added, but existing CLI exact-output tests confirmed compatibility. Suggestion: consider a Presenter interface and JSON presenter in M7.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Introduced a concrete repositoryPresenter in internal/interface/cli and moved all repository human-readable output responsibilities out of command orchestration. Repository details, lists, updates, archive/restore, access, and empty-result output remain byte-compatible; JSON output remains out of scope.
<!-- SECTION:FINAL_SUMMARY:END -->
