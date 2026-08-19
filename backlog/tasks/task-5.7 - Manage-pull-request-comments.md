---
id: TASK-5.7
title: Manage pull request comments
status: Done
assignee:
  - '@claude'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:13'
labels: []
milestone: m-4
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-5
priority: medium
ordinal: 20040
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Add and inspect pull request discussion.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Comments remain associated with the correct pull request
- [x] #2 Empty content and permission failures are handled consistently
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (batch go-ahead for M4 remainder). Add fj pr comment list/add mirroring issue comments; the adapter verifies the number is a pull request before using the shared issue-comment endpoint so comments stay associated with a pull request.

Added fj pr comment list/add. The adapter issues GET /pulls/{index} before touching the shared /issues/{index}/comments endpoint so a non-pull-request number fails with 'pull request not found' instead of attaching a comment to an issue.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj pr comment list and add with pull-request existence verification, empty-body validation, and consistent permission error mapping. Verified with unit tests at all three layers and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
