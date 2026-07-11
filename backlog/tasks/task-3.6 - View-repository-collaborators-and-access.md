---
id: TASK-3.6
title: View repository collaborators and access
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 05:50'
labels: []
dependencies:
  - TASK-2.9
  - TASK-2.10
  - TASK-2.11
  - TASK-2.12
  - TASK-2.13
  - TASK-2.14
references:
  - ROADMAP.md
parent_task_id: TASK-3
ordinal: 29000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Make repository access relationships observable.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Collaborators and access information are presented consistently
- [x] #2 Permission failures are distinct from missing repositories
- [x] #3 The CLI provides fj repo access OWNER/NAME with optional --instance PROFILE and existing target/instance rules.
- [x] #4 Application owns RepositoryAccess and Collaborator models plus a dedicated AccessViewer port; existing Service, Getter, Creator, Updater, and Archiver contracts are unchanged.
- [x] #5 The REST adapter uses GET /api/v1/repos/{owner}/{repo}/collaborators, safely encodes path segments, performs no pagination, and normalizes permissions with admin > write > read.
- [x] #6 Success output is Repository: owner/name, Collaborators:, and - username: permission lines; empty results display No collaborators found.
- [x] #7 401/403 classify as authentication, 404 as remote repository not found, and other HTTP, transport, and JSON failures as remote operation failed; secrets and raw causes are not exposed; JSON and mutations are out of scope.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — new public read CLI, Application access model/port, REST endpoint, and permission normalization require architecture-aware implementation.

Validation: make pre-commit passed, including git diff --check, go vet ./..., and go test ./....

Independent post-implementation review (GPT-5): Critical none, Major none. Minor: direct REST ViewAccess and redaction coverage could be expanded in future hardening; no scope expansion made.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented repository access viewing. Added Application-owned RepositoryAccess/Collaborator/Permission models and dedicated AccessViewer port, Forgejo collaborators GET adapter with admin > write > read normalization, fj repo access command, fixed collaborator output and empty-result handling, and safe authentication/remote classification. Pagination, JSON, and mutations remain out of scope.
<!-- SECTION:FINAL_SUMMARY:END -->
