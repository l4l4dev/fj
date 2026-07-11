---
id: TASK-4.8.1
title: Manage issue labels
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 06:11'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-3
dependencies:
  - TASK-2.9
parent_task_id: TASK-4.8
ordinal: 86000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Manage issue labels independently through explicit add and remove workflows. This task is one 30-90 minute slice of Issue metadata management.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The CLI provides fj issue label add OWNER/NAME NUMBER LABEL and fj issue label remove OWNER/NAME NUMBER LABEL with --instance PROFILE.
- [x] #2 OWNER/NAME is validated and issue number must be a positive integer; label values must be non-empty.
- [x] #3 Application defines an issue-label model and dedicated Port contracts without changing existing issue or repository Ports.
- [x] #4 Infrastructure uses the approved Forgejo label endpoints and translates transport, permission, not-found, and remote failures through the Application error boundary.
- [x] #5 Adding or removing a label does not overwrite unrelated labels, milestone, assignment, title, body, or state.
- [x] #6 Human-readable output is defined and implemented for label add/remove without changing existing issue command output.
- [x] #7 Tests cover validation, add/remove delegation, API mapping, error boundary, Presenter output, and explicit Composition Root injection.
- [x] #8 Label replace, milestone, assignment, metadata bulk update, JSON, and pagination remain out of scope.
- [x] #9 The CLI provides fj issue label add OWNER/NAME NUMBER LABEL and fj issue label remove OWNER/NAME NUMBER LABEL with --instance; LABEL is single-valued and issue number is a positive integer.
- [x] #10 Empty and whitespace-only labels are rejected locally; label replace, bulk update, multiple-label operations, label listing, JSON, and pagination remain out of scope.
- [x] #11 Application owns Label{ID int64, Name string}, AddLabelRequest, RemoveLabelRequest, LabelAdder, and LabelRemover in internal/application/issue; Repository ports remain unchanged.
- [x] #12 Duplicate label add and missing label removal are idempotent successes.
- [x] #13 Infrastructure uses POST /api/v1/repos/{owner}/{repo}/issues/{index}/labels for add and DELETE /api/v1/repos/{owner}/{repo}/issues/{index}/labels/{id} for remove, resolving label name to ID for removal.
- [x] #14 The operations are add issue label and remove issue label; failures use the existing apperror boundary without exposing HTTP status, credentials, URL details, response bodies, or raw causes to CLI.
- [x] #15 Add output is Issue: #<number> followed by Label added: <label>; remove output is Issue: #<number> followed by Label removed: <label>.
- [x] #16 Explicit dependency injection and Interface → Application Port → Infrastructure adapter direction are preserved; Repository ports remain unchanged and no runtime type assertions are added.
- [x] #17 Infrastructure converts a CLI label name to POST body {"labels":["<label>"]}.
- [x] #18 For remove, Infrastructure may internally retrieve issue labels, match by name, obtain the label ID, and DELETE the matching label; no user-facing label listing command or Presenter is added.
- [x] #19 Duplicate add resolves existing labels and succeeds without mutation when the label already exists; missing remove resolves existing labels and succeeds without mutation when the label does not exist.
- [x] #20 Label ID resolution classifies issue/repository failures as remote, permission failures as authentication, and API failures as remote.
- [x] #21 Scope remains single-label add/remove only; label listing command, replace, bulk update, multiple-label operation, JSON, and pagination remain out of scope.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design:
- Scope: label add and label remove only. Label replace, bulk update, multiple-label operations, label listing, JSON, and pagination are out of scope.
- CLI: `fj issue label add OWNER/NAME NUMBER LABEL` and `fj issue label remove OWNER/NAME NUMBER LABEL`, both with `--instance PROFILE`. LABEL is single-valued, non-empty, non-whitespace; issue number is positive.
- Semantics: duplicate add is idempotent success; removing a missing label is idempotent success.
- Application: Label{ID int64, Name string}, AddLabelRequest, RemoveLabelRequest, LabelAdder, and LabelRemover in internal/application/issue. Repository ports remain unchanged.
- Infrastructure: add uses POST /api/v1/repos/{owner}/{repo}/issues/{index}/labels; remove uses DELETE /api/v1/repos/{owner}/{repo}/issues/{index}/labels/{id}. CLI accepts label name; Infrastructure resolves label name to ID for removal.
- Error Boundary: use existing apperror; operations are `add issue label` and `remove issue label`; HTTP status is not exposed to CLI.
- Presenter: add outputs Issue: #<number> and Label added: <label>; remove outputs Issue: #<number> and Label removed: <label>.
- Architecture: preserve Interface → Application Port → Infrastructure adapter, explicit DI, unchanged Repository ports, and no runtime type assertions.

Additional approved decisions:
- Add API body is {"labels":["<label>"]}; Infrastructure converts the CLI label name to this API shape.
- Remove accepts a label name. Infrastructure may internally GET issue labels, find the matching name, obtain its ID, and DELETE the label. Label listing remains out of scope for user-facing commands, Presenter output, and management features; the internal resolution GET is allowed.
- Idempotency: duplicate add first resolves existing labels and succeeds without an API mutation when the name exists; missing remove first resolves existing labels and succeeds without an API mutation when the name does not exist.
- Label ID resolution failures classify issue/repository failures as remote, permission failures as authentication, and API failures as remote.
- Scope remains limited to single-label add/remove; no label list command, replace, bulk update, multiple-label operation, JSON, or pagination.

Verification:
- gofmt -l .
- git diff --check
- go vet ./...
- go test ./...
- make pre-commit

All listed verification commands succeeded.

Independent Review:
Critical: none
Major: none
Minor: HTTP error boundary, idempotency suppression, path encoding, secret redaction, DI, and CLI delegation tests have expansion opportunities.
Suggestion: consider adding the remaining boundary tests.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented TASK-4.8.1 issue label management with single-label add/remove commands, Application-owned label ports, response DTO conversion with ID preservation, internal name-to-ID resolution, idempotent duplicate add and missing remove behavior, safe error mapping, explicit dependency injection, and fixed Presenter output. Existing issue commands and out-of-scope behavior remain unchanged.
<!-- SECTION:FINAL_SUMMARY:END -->
