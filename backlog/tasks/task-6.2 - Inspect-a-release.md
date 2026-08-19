---
id: TASK-6.2
title: Inspect a release
status: Done
assignee:
  - '@claude-sonnet'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:32'
labels: []
milestone: m-5
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-6
ordinal: 47000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Display release metadata, tag relationships, and assets.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Release and tag identities are unambiguous
- [x] #2 Missing and inaccessible releases produce distinct errors
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (M5 go-ahead). fj release inspect OWNER/NAME TAG resolving by tag via /releases/tags/{tag}, showing metadata, tag, and assets; 404 vs 401/403 produce distinct errors. Implementation delegated to a Sonnet subagent; verified by the orchestrator session.

Implemented by Sonnet subagent; verified by orchestrator (make pre-commit 16 packages ok). Release resolved by tag via /releases/tags/{tag}; 404 maps to 'release not found', 401/403 to authentication errors.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj release inspect showing metadata, tag, state, notes, and assets, with distinct not-found and permission errors. Verified with unit tests at all three layers and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
