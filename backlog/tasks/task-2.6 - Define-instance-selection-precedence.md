---
id: TASK-2.6
title: Define instance-selection precedence
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 05:50'
labels: []
dependencies: []
references:
  - ROADMAP.md
modified_files:
  - README.md
  - internal/application/config/selection.go
  - internal/application/config/selection_test.go
parent_task_id: TASK-2
ordinal: 20000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Make effective Forgejo instance selection deterministic.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Selection precedence is documented and tested
- [x] #2 Missing or ambiguous selection produces a clear error
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Document deterministic instance-selection precedence in README.md.
2. Add an Application-layer selection operation that prefers an explicit profile name and otherwise selects the sole configured instance.
3. Return clear errors for a missing explicit profile or ambiguous implicit selection.
4. Add focused tests for explicit selection, sole-profile fallback, missing profiles, and ambiguity.
5. Run formatting and all Go tests.
6. Check all acceptance criteria and finalize TASK-2.6 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Selection is limited to explicit profile name over sole-profile fallback. Environment variables, configured defaults, CLI flags, and authentication are outside TASK-2.6.

Documented selection precedence: an explicit profile name takes priority; otherwise the sole configured profile is selected.
Added Configuration.SelectInstance in the Application layer after configuration validation.
Missing explicit profiles and ambiguous multi-profile selection return clear deterministic errors.
Validation passed: gofmt completed and go test ./... passed.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Documented and implemented minimal deterministic Forgejo instance selection with explicit-profile precedence and sole-profile fallback. Added focused missing and ambiguous selection tests; formatting and all Go tests pass.
<!-- SECTION:FINAL_SUMMARY:END -->
