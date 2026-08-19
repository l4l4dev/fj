---
id: TASK-6.5
title: Publish a release
status: Done
assignee:
  - '@claude-sonnet'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:45'
labels: []
milestone: m-5
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-6
ordinal: 50000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Transition a draft release to published state safely.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The target release and publication intent are explicit
- [x] #2 Failed publication does not report success
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (M5 go-ahead). fj release publish OWNER/NAME TAG: use case inspects the release first, rejects already-published releases, then flips draft=false by id; failures never print success. Implementation delegated to a Sonnet subagent; verified by the orchestrator session.

Implemented by Sonnet subagent; verified by orchestrator (make pre-commit 16 packages ok). Use case inspects first and rejects non-draft releases with Conflict before calling the publisher; failures produce errors, never success output.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj release publish transitioning drafts to published by resolved id, rejecting already-published releases. Verified with unit tests at all three layers and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
