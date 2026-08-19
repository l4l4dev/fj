---
id: TASK-6.1
title: List releases
status: Done
assignee:
  - '@claude-sonnet'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:28'
labels: []
milestone: m-5
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-6
ordinal: 46000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Discover releases in the selected repository.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Draft, published, and pre-release states are distinguishable
- [x] #2 Pagination and empty results are clear
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (M5 go-ahead). fj release list OWNER/NAME with draft/prerelease/published labels, pagination flags, and empty-result output, mirroring pr list across the three layers. Implementation delegated to a Sonnet subagent; verified by the orchestrator session.

Implemented by Sonnet subagent; verified by orchestrator (make pre-commit 16 packages ok, CLI arg validation exercised).
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj release list with draft/prerelease/published state labels, page/limit pagination, and empty-result output. Verified with unit tests at all three layers and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
