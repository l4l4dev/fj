---
id: TASK-3.2
title: Inspect a repository
status: Done
assignee: []
created_date: '2026-07-10 11:55'
updated_date: '2026-07-10 23:38'
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
modified_files:
  - internal/application/repository/repository.go
  - internal/application/repository/inspect.go
  - internal/application/repository/inspect_test.go
  - internal/infrastructure/repository/rest.go
  - internal/infrastructure/repository/rest_test.go
  - internal/interface/cli/errors.go
  - internal/interface/cli/root.go
  - internal/interface/cli/repository.go
  - internal/interface/cli/repository_test.go
parent_task_id: TASK-3
ordinal: 25000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Display the essential state and metadata of one repository.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The target instance and repository are explicit
- [x] #2 Missing and inaccessible repositories produce distinct errors
- [x] #3 The CLI provides fj repo inspect OWNER/NAME with an optional --instance PROFILE flag.
- [x] #4 OWNER/NAME is one required positional argument; exactly one slash is required and empty owner/name or extra slashes are validation errors.
- [x] #5 The application-owned Repository model includes Owner, Name, Description, Private, Archived, and DefaultBranch, and InspectUseCase depends only on a dedicated Getter port.
- [x] #6 The REST adapter uses GET /api/v1/repos/{owner}/{repo} with safely encoded path segments and translates the response into the application model.
- [x] #7 Success output uses the fixed order Repository, Description, Private, Archived, and Default branch; empty Description is rendered as -.
- [x] #8 404 remains the remote category with safe message repository not found; 401/403 are authentication and other HTTP, transport, and JSON failures are remote.
- [x] #9 The existing instance selection and Composition Root path are reused without exposing credentials, sensitive URL parts, response bodies, or raw transport causes; JSON, extra metadata, and fetch-all remain out of scope.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Pre-implementation check (GPT-5): Critical none. Major: implementation is blocked because the approved specification does not define the inspect command syntax/flags, repository detail fields, Get service port, output format, target syntax, or 404 versus 401/403 error classification. Existing repository.Service only provides List and Repository only has Owner/Name. Human design approval is required before implementation; no code changes made.

Model: GPT-5 — approved Major Change implementation involving a new public CLI command, application model/port, REST endpoint, and security-sensitive error classification.

Validation: make pre-commit passed, including git diff --check, go vet ./..., and go test ./....

Post-implementation review (GPT-5): Critical none, Major none. Minor: Backlog acceptance/status required finalization and was corrected. Suggestion: inspect transport operation remains internal while CLI uses fixed inspect repository; no change required.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented TASK-3.2 repository inspection without changing the existing List Service contract. Added the application-owned detail model and dedicated Getter port, InspectUseCase validation, Forgejo GET /api/v1/repos/{owner}/{repo} adapter with safe path encoding, fj repo inspect OWNER/NAME with optional --instance, fixed human-readable output, and safe 404/401/403/remote error classification. JSON, extra metadata, and fetch-all remain out of scope.
<!-- SECTION:FINAL_SUMMARY:END -->
