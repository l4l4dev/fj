---
id: TASK-6.4
title: Update release metadata
status: Done
assignee:
  - '@claude-sonnet'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:41'
labels: []
milestone: m-5
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-6
ordinal: 49000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Modify supported fields without replacing unspecified values.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Only supplied fields are changed
- [x] #2 Draft and pre-release changes are observable
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (M5 go-ahead). fj release update OWNER/NAME TAG with pointer-based partial updates (--title/--notes/--prerelease) and changed-field output mirroring pr update; resolves the release by tag first, then PATCHes by id. Implementation delegated to a Sonnet subagent; verified by the orchestrator session.

Implemented by Sonnet subagent; verified by orchestrator (make pre-commit 16 packages ok). PATCH payload contains only supplied fields; output names changed fields and the resulting draft/prerelease state.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj release update with pointer-based partial updates for title/notes/prerelease, changed-field output, and resulting state display. Verified with unit tests at all three layers and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
