---
id: TASK-6.7
title: Delete a release
status: Done
assignee:
  - '@claude-sonnet'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 10:15'
labels: []
milestone: m-5
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-6
ordinal: 52000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Remove a release only with explicit intent and target information.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The release, tag relationship, and expected effect are visible
- [x] #2 Cancellation and remote rejection are reported accurately
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (M5 go-ahead). fj release delete OWNER/NAME TAG --confirm: use case composes Inspector + a new Deleter, deletes by resolved id; success output names the tag and title and states the git tag is kept; missing --confirm and remote rejections are distinct errors. Implementation delegated to a Sonnet subagent; verified by the orchestrator session.

Implemented by Sonnet subagent; verified by orchestrator (make pre-commit 16 packages ok). --confirm is checked before any remote call; success output names the tag and title and states the git tag is kept; remote rejections map to distinct errors.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj release delete gated by --confirm, deleting by resolved id and reporting the deleted release and preserved git tag. Verified with unit tests at all three layers and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
