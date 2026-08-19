---
id: TASK-6.3
title: Create a draft release
status: Done
assignee:
  - '@claude-sonnet'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:36'
labels: []
milestone: m-5
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-6
ordinal: 48000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Create a release without publishing it immediately.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Tag, title, notes, and draft state are explicit
- [x] #2 Invalid tag or metadata input is rejected clearly
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (M5 go-ahead). fj release create always creates a draft (publishing is TASK-6.5); explicit --tag/--title required, --notes/--prerelease optional; invalid tag or duplicate release maps to clear errors. Implementation delegated to a Sonnet subagent; verified by the orchestrator session.

Implemented by Sonnet subagent; verified by orchestrator (make pre-commit 16 packages ok, whitespace-tag rejection exercised via the binary). Create always sends draft:true; publishing stays a separate explicit operation (TASK-6.5).
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj release create producing draft releases with explicit --tag/--title and optional --notes/--prerelease; invalid tags, duplicates, and remote rejections map to clear errors. Verified with unit tests at all three layers and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
