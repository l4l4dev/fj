---
id: TASK-5.6
title: Submit a pull request review
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-08-19 07:55'
labels: []
milestone: m-4
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
modified_files:
  - internal/application/pullrequest/pullrequest.go
  - internal/application/pullrequest/review.go
  - internal/application/pullrequest/review_test.go
  - internal/infrastructure/pullrequest/rest.go
  - internal/infrastructure/pullrequest/rest_test.go
  - internal/interface/cli/pullrequest.go
  - internal/interface/cli/pullrequest_presenter.go
  - internal/interface/cli/pullrequest_test.go
  - internal/interface/cli/root.go
parent_task_id: TASK-5
priority: low
ordinal: 40010
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Support comment, approval, and change-request outcomes.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The intended review outcome is explicit
- [x] #2 Unsupported or unauthorized outcomes return actionable errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add Application review submission types, outcome mapping, validation, and tests.
2. Add Forgejo REST review submission payload/response/error mapping and tests.
3. Add fj pr review command, dependency composition, presenter output, and CLI tests.
4. Run focused tests and make pre-commit.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision: Human approved the TASK-5.6 public CLI contract, validation, success output, error mapping, and scope boundaries on 2026-07-16.

Model: GPT-5 — public CLI contract implementation requiring architecture-sensitive validation and verification

Pre-implementation check (GPT-5): Critical none; Major none; Minor none; additional Human Decision not required. The approved CLI contract is consistent with the Constitution, Architecture, and M4 roadmap. Submission outcome/event types remain separate from fetched review states.

Implemented Application-owned outcome mapping and validation, Forgejo review submission through the existing JSON transport, actionable HTTP error translation, fj pr review dependency composition, and the approved minimal success output. Inline comments, pending reviews, commit selection, dismissal, file input, interactive selection, and automation were not added.

Verification: focused go test ./internal/application/pullrequest ./internal/infrastructure/pullrequest ./internal/interface/cli passed; make pre-commit passed git diff --check, go vet ./..., and go test ./.... The initial sandboxed focused test attempt could not access the Go build cache; the approved elevated rerun passed.

Independent Review: not performed because the human explicitly excluded it from this implementation request. Task remains In Progress and is not finalized as Done.

Independent review (Fable 5 subagent, 2026-08-19): Review Ready. Critical/Major: none. Minor 1 (deferred, needs human decision): mapApplicationError drops ValidationError messages CLI-wide, so an unsupported --outcome prints only 'invalid input'; fixing it touches shared behavior beyond this task. Suggestions 2 recorded: outcome validation runs after credential resolution; presenter prints remote state verbatim. Scope, architecture, REST contract, security, and tests all confirmed.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj pr review submitting comment/approve/request-changes outcomes via the Forgejo REST API with application-owned validation and actionable error mapping. Verified with focused go tests and make pre-commit; independent review Review Ready.
<!-- SECTION:FINAL_SUMMARY:END -->
