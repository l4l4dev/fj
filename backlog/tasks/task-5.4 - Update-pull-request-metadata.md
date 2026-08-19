---
id: TASK-5.4
title: Update pull request metadata
status: Done
assignee:
  - '@claude'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:11'
labels: []
milestone: m-4
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-5
priority: medium
ordinal: 20030
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Change supported fields without altering unspecified state.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Only explicitly supplied fields are updated
- [x] #2 The result identifies changed fields
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (batch go-ahead for M4 remainder). Add fj pr update --title/--body mirroring issue update: application UpdateUseCase, REST PATCH adapter, CLI command with changed-field output.

Added fj pr update with pointer-based partial updates (only supplied flags sent as PATCH fields) and changed-field output. Consolidated the pr command constructors into one pullRequestDependencies struct while wiring the new dependency.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj pr update --title/--body updating only explicitly supplied fields via PATCH /pulls/{index} and reporting changed fields. Verified with unit tests at all three layers and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
