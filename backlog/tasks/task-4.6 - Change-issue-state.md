---
id: TASK-4.6
title: Change issue state
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 05:43'
labels: []
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-4
ordinal: 35000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Open, close, or reopen an issue safely.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The intended state transition is explicit
- [x] #2 Invalid or redundant transitions are predictable
- [x] #3 The CLI provides fj issue state OWNER/NAME NUMBER --state open|closed with --instance; issue number is a positive integer and invalid state values are rejected locally.
- [x] #4 Application owns ChangeStateRequest and StateChanger in internal/application/issue; existing issue and repository ports remain unchanged.
- [x] #5 Open and closed are supported, reopen is represented by --state open, and redundant transitions succeed idempotently.
- [x] #6 Infrastructure sends PATCH /api/v1/repos/{owner}/{repo}/issues/{index} with a state-only JSON body, safely encodes path segments, and converts the DTO to IssueDetail.
- [x] #7 The operation name is change issue state; validation, authentication, remote, and internal failures use the existing apperror boundary without HTTP status classification in CLI or secret leakage.
- [x] #8 The Presenter outputs Issue: #<number> and State: <state> for state changes; existing list, filter, inspect, create, and update output remains unchanged.
- [x] #9 Explicit dependency injection and Interface → Application StateChanger Port → Infrastructure adapter direction are preserved; no runtime type assertions are added and Application remains independent of HTTP, Cobra, and Forgejo DTOs.
- [x] #10 Tests cover CLI and state validation, redundant transitions, Application delegation, PATCH endpoint, state-only JSON, DTO conversion, error boundary, Presenter output, and Composition Root injection. Assignee, labels, milestone, comments, metadata, and JSON remain out of scope.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design:
- CLI: `fj issue state OWNER/NAME NUMBER --state open|closed` with `--instance PROFILE`; issue number is positive and state accepts only open or closed. No confirmation, interactive mode, or dedicated reopen command; reopen is expressed as --state open.
- Redundant transitions are idempotent success.
- Application: add ChangeStateRequest{Owner, Name, Number, State} and StateChanger in internal/application/issue. Existing Lister, Inspector, Creator, and Updater remain unchanged; Repository ports remain unchanged.
- Infrastructure: PATCH /api/v1/repos/{owner}/{repo}/issues/{index} with JSON body {"state":"open|closed"}; safely encode path segments, convert DTO to IssueDetail, and use the existing apperror boundary. Operation is `change issue state`. HTTP status is not classified in CLI; credentials, URL details, response bodies, and raw causes are not exposed.
- Presenter: add state-change output `Issue: #<number>` and `State: <state>`. Existing list, filter, inspect, create, and update output remains unchanged.
- Architecture: preserve Interface → Application StateChanger Port → Infrastructure adapter, explicit DI, no new runtime type assertions, and Application independence from HTTP/Cobra/Forgejo DTOs.
- Out of scope: assignee, labels, milestone, comments, metadata, and JSON.

Independent Review: Critical: none. Major: none. Minor: redundant-transition success case, HTTP error boundary, JSON decode failure, path encoding, state-only JSON body, secret/raw cause redaction, Composition Root StateChanger injection, CLI-to-StateChanger delegation, and existing list/filter/inspect/create/update regression tests have expansion opportunities. Suggestion: fix the idempotent same-state success contract in a future test.

Verification:
- gofmt -l .
- git diff --check
- go vet ./...
- go test ./...
- make pre-commit

All listed verification commands succeeded.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented fj issue state OWNER/NAME NUMBER with Application-owned StateChanger, open/closed validation, state-only PATCH requests, explicit StateChanger dependency injection, and fixed state-change output while preserving existing issue command output. Independent review found no Critical or Major issues; remaining Minor items are boundary-test expansion opportunities and a future idempotency-test Suggestion.
<!-- SECTION:FINAL_SUMMARY:END -->
