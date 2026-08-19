---
id: TASK-6.6
title: Manage release assets
status: Done
assignee:
  - '@claude-opus'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 08:52'
labels: []
milestone: m-5
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-6
ordinal: 51000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Support inspection, upload, and removal of release assets.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Assets remain associated with the correct release
- [x] #2 Duplicate, missing, and transfer failures are clear
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19 (M5 go-ahead). fj release asset list/upload/delete addressing assets by release tag + asset name; adapter resolves tag to id before every operation so assets stay associated with the right release; multipart upload via a new content-type-aware DoRaw method on the shared forgejo client. Implementation delegated to an Opus subagent (heavier scope); verified by the orchestrator session.

Implemented by Opus subagent; verified by orchestrator (make pre-commit 16 packages ok). Both asset operations resolve the tag to a release id via a shared helper before acting; upload is multipart via a new content-type-aware DoRaw on the shared client; delete requires --confirm and maps duplicates (409) and missing assets to distinct errors.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj release asset list/upload/delete with tag-to-id resolution, multipart upload, --confirm-gated deletion, and clear duplicate/missing/transfer errors. Verified with unit tests at all four affected packages and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
