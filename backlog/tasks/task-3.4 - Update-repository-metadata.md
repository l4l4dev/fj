---
id: TASK-3.4
title: Update repository metadata
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-2
dependencies:
  - TASK-2.9
  - TASK-2.10
  - TASK-2.11
  - TASK-2.12
  - TASK-2.13
  - TASK-2.14
references:
  - ROADMAP.md
modified_files:
  - internal/application/repository/repository.go
  - internal/application/repository/update.go
  - internal/application/repository/update_test.go
  - internal/infrastructure/repository/rest.go
  - internal/infrastructure/repository/rest_test.go
  - internal/interface/cli/root.go
  - internal/interface/cli/repository.go
  - internal/interface/cli/repository_test.go
parent_task_id: TASK-3
ordinal: 27000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Modify supported metadata without changing unrelated settings.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Only explicitly supplied fields are changed
- [x] #2 The result identifies the repository and changed fields
- [x] #3 The CLI provides fj repo update OWNER/NAME with optional --description, --visibility public|private, and --instance PROFILE; OWNER/NAME uses inspect validation and at least one update field is required.
- [x] #4 Application adds UpdateRequest with Owner, Name, *string Description, *bool Private and a dedicated Updater port without changing Service, Getter, or Creator; nil fields are omitted and pointer fields are sent.
- [x] #5 The REST adapter uses PATCH /api/v1/repos/{owner}/{repo}, safely encodes path segments, and sends only explicitly supplied fields.
- [x] #6 Success output uses fixed order Repository, Changed fields, Description, Private, Archived, Default branch; changed fields are ordered description then visibility, and empty values display as -.
- [x] #7 Invalid input, no changes, and invalid visibility are validation; 401/403 authentication; 404 remote repository not found; 409 remote repository update conflict; other HTTP, transport, and JSON failures remote; unexpected failures internal.
- [x] #8 Credentials, sensitive URLs, request/response bodies, and raw transport causes are not displayed or retained; existing list, inspect, and create behavior remains compatible and is regression-tested.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — approved Major Change involving a mutating public CLI, PATCH transport body, dedicated Updater port, and compatibility-sensitive error/output contracts.

Validation: make pre-commit passed, including git diff --check, go vet ./..., and go test ./....

Independent post-implementation review (GPT-5): initial Major finding that unexpected Updater errors were classified as validation was resolved with an Application ValidationError; re-review Critical none, Major none. Minor: direct internal-error CLI coverage could be added later; no scope expansion.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented TASK-3.4 repository metadata updates. Added dedicated Updater port and Use Case with pointer-based explicit-field semantics, PATCH /api/v1/repos/{owner}/{repo} adapter, fj repo update OWNER/NAME with description and visibility, fixed changed-field output, safe 401/403/404/409/remote/internal classification, and regression coverage while preserving List, Getter, and Creator contracts.
<!-- SECTION:FINAL_SUMMARY:END -->
