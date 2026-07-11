---
id: TASK-2.8
title: Add repository context selection
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
  - internal/application/target/context.go
  - internal/application/target/context_test.go
parent_task_id: TASK-2
ordinal: 22000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Resolve the repository targeted by a command explicitly.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The effective instance and repository are observable before execution
- [x] #2 Missing or conflicting context produces an actionable error
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Define an Application-layer repository target request and observable context containing a safe instance profile and repository owner/name.
2. Resolve the instance through existing deterministic selection and resolve explicit versus detected repository context.
3. Return actionable errors for missing repository fields, conflicting repository sources, and inherited missing or ambiguous instance selection.
4. Add focused tests for effective context, fallback, missing fields, conflicts, and ambiguous instances.
5. Run formatting and all Go tests.
6. Check all acceptance criteria and finalize TASK-2.8 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Repository context resolution is limited to Application-layer inputs and output. CLI flags, repository auto-detection, remote access, and common error presentation are outside TASK-2.8.

Added an Application-layer target Context containing a credential-free instance Profile and repository owner/name.
Added Resolve to reuse deterministic instance selection and reconcile explicit and detected repository context.
Missing repository context or fields, conflicting repositories, and inherited missing or ambiguous instance selection return actionable errors.
CLI flags, repository detection, and remote operations were not added.
Validation passed: gofmt completed and go test ./... passed.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added minimal Application-layer repository context selection with an observable credential-free target and actionable missing or conflicting context errors. Added focused tests; formatting and all Go tests pass.
<!-- SECTION:FINAL_SUMMARY:END -->
