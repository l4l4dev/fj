---
id: TASK-3.7.3
title: Presenter Boundary Introduction
status: To Do
assignee: []
created_date: '2026-07-11 00:31'
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
- [ ] #1 Human-readable output responsibilities are separated from command orchestration.
- [ ] #2 Existing M1/M2 output text, field order, empty-value rendering, and stream behavior remain unchanged.
- [ ] #3 The boundary can support a future JSON presenter without implementing JSON or changing the M7 scope.
- [ ] #4 Focused presenter and command regression tests pass, including make pre-commit.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Introduce a human presenter boundary after explicit composition wiring, migrate existing repository output, and verify byte-compatible behavior.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Third child of TASK-3.7. Depends on Explicit Composition Root Refactoring and must complete before M3.
<!-- SECTION:NOTES:END -->
