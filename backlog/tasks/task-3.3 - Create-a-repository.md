---
id: TASK-3.3
title: Create a repository
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
  - internal/application/repository/create.go
  - internal/application/repository/create_test.go
  - internal/infrastructure/forgejo/client.go
  - internal/infrastructure/repository/rest.go
  - internal/infrastructure/repository/rest_test.go
  - internal/interface/cli/root.go
  - internal/interface/cli/repository.go
  - internal/interface/cli/repository_test.go
parent_task_id: TASK-3
ordinal: 26000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Support safe repository creation with explicit settings.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Invalid input is rejected before the remote request
- [x] #2 The CLI provides fj repo create NAME with optional --description, --visibility public|private (default private), and --instance PROFILE flags; NAME must be non-blank.
- [x] #3 The command creates only the authenticated user's personal repository through POST /api/v1/user/repos and does not expose --owner, organization creation, prompts, or --dry-run.
- [x] #4 Application adds CreateRequest and dedicated Creator port without changing existing Service or Getter contracts; the existing Repository model represents the result and the REST adapter asserts Creator.
- [x] #5 Success output uses the inspect format: Repository, Description (or -), Private, Archived (false), and Default branch (or -), without extra success text or JSON.
- [x] #6 Invalid input, invalid flags, and invalid visibility are validation; 401/403 are authentication; 409 is remote with repository already exists; other HTTP, transport, and JSON failures are remote; unexpected failures are internal.
- [x] #7 Credential, sensitive URL parts, response bodies, and raw transport causes are never retained or displayed; tests cover Application, REST, CLI, visibility, description, error classes, and redaction.
- [x] #8 The selected instance is resolved by the existing --instance rules, and successful output makes owner/name and visibility observable.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — approved Major Change involving a mutating public CLI, POST transport body, new application port, and security-sensitive error classification.

Post-implementation review (GPT-5): Major finding resolved by aligning AC#1 with the approved output contract (instance selection is input behavior; instance is not printed). Minor operation-name fallback was corrected so REST transport failures use create repository for Create.

Validation: make pre-commit passed, including git diff --check, go vet ./..., and go test ./....

Independent post-implementation review (GPT-5) after correction: Critical none, Major none, Minor: typed transport operation may remain request internally while CLI operation is fixed and safe; Suggestion: add direct REST failure/redaction tests in a future hardening task. No scope expansion made.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented TASK-3.3 personal repository creation. Added a dedicated Creator port and use case, JSON POST transport support, Forgejo POST /api/v1/user/repos adapter, fj repo create with description/private visibility and instance selection, inspect-compatible output, and safe validation/authentication/remote error classification including repository already exists. Existing List and Getter contracts remain unchanged; organization creation, prompts, dry-run, JSON, and extra metadata remain out of scope.
<!-- SECTION:FINAL_SUMMARY:END -->
