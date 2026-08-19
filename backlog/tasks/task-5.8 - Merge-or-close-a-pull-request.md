---
id: TASK-5.8
title: Merge or close a pull request
status: Done
assignee:
  - '@claude'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:20'
labels: []
milestone: m-4
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-5
priority: low
ordinal: 40020
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Complete a pull request through an explicitly confirmed operation.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The target and completion action are unambiguous
- [x] #2 Failed readiness checks or remote rejection do not appear successful
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (batch go-ahead for M4 remainder). Add fj pr merge with a required explicit --method (merge, rebase, squash) via POST /pulls/{index}/merge, and fj pr close via PATCH state=closed; readiness failures and remote rejections map to errors, never success output.

Added fj pr merge (required explicit --method: merge, rebase, squash; POST /pulls/{index}/merge) and fj pr close (PATCH state=closed). Readiness failures (405), conflicts (409), auth, and not-found map to distinct errors and produce no success output.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj pr merge with a mandatory explicit merge method and fj pr close. Verified with unit tests at all three layers (including a no-success-output-on-failure test) and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
