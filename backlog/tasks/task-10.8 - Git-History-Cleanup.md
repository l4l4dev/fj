---
id: TASK-10.8
title: Git History Cleanup
status: Done
assignee: []
created_date: '2026-07-11 10:29'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-9
dependencies:
  - TASK-10.6
parent_task_id: TASK-10
priority: high
ordinal: 90000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The privacy audit (TASK-10.6) detected an absolute local home directory path inside one Backlog task file (task-10.2), present in three commits. Rewrite git history with git filter-repo --replace-text so the path prefix is replaced with $HOME, then update origin/main with a lease-protected force push. Detected values are masked here per the TASK-10.6 decision plan; details were reported in chat only.

Note: task-10.2 file content is modified by filter-repo directly, not via the backlog CLI. This is an accepted exception because the CLI cannot rewrite git history; the replacement is a mechanical string substitution that does not touch frontmatter or metadata.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 All occurrences of the detected absolute home directory path in git history are replaced with $HOME.
- [x] #2 A case-insensitive scan of every commit for the audit detection terms returns zero matches.
- [x] #3 Commits unaffected by the replacement retain their original hashes, and total commit count is unchanged.
- [x] #4 A full backup bundle of all refs exists outside the repository before the rewrite.
- [x] #5 origin/main is updated using --force-with-lease pinned to the expected pre-rewrite head SHA.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create backup bundle outside repo. 2. Write replace-text expressions file in session scratchpad. 3. Run git filter-repo --replace-text --force. 4. Verify history is clean and unaffected hashes preserved. 5. Re-add origin, fetch, force-with-lease push. 6. Finalize task.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: Claude Fable 5 (claude-fable-5). Reason: destructive history rewrite with force push requires careful verification, per AGENTS.md Section 14/15. The backlog CLI has no official field for model selection, so it is recorded here in notes.

Executed 2026-07-11. Backup: ../fj-backup-20260711.bundle (verified, all refs, old main=b801e9a, old origin/main=17e9285). Rewrite: git filter-repo --replace-text with one expression replacing the absolute home path prefix with $HOME. 3 target commits rewritten; unaffected ancestors (verified 777248a, 06b1704) kept hashes; commit count unchanged (67). One additional occurrence found during verification in the freshly created task-10.7 description; removed via backlog CLI and commit --amend before push. Full-history case-insensitive scan after rewrite: zero matches. Pushed with --force-with-lease pinned to 17e9285; forced update accepted (17e9285 -> 59f9e21).
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Rewrote git history with filter-repo to replace an absolute home directory path with $HOME in 3 commits, verified the full history scans clean, and force-pushed origin/main with a pinned lease. Backup bundle retained outside the repository.
<!-- SECTION:FINAL_SUMMARY:END -->
