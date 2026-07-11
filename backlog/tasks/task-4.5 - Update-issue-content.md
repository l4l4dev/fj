---
id: TASK-4.5
title: Update issue content
status: Done
assignee: []
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 05:21'
labels: []
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-4
ordinal: 34000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Modify supported issue fields without changing unspecified fields.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Only explicitly supplied fields are updated
- [ ] #2 The result identifies all changed fields
- [ ] #3 The CLI provides fj issue update OWNER/NAME NUMBER with --instance and optional --title/--body flags; issue number is a positive integer and at least one update flag is required.
- [ ] #4 Application owns UpdateRequest and Updater in internal/application/issue; nil means unspecified and pointer-to-empty means explicit empty update; Repository ports remain unchanged.
- [ ] #5 Only title and body are updated; state, assignee, labels, milestone, comments, metadata, and JSON remain out of scope.
- [ ] #6 Infrastructure calls PATCH /api/v1/repos/{owner}/{repo}/issues/{index}, sends only specified title/body fields, safely encodes path segments, and converts the DTO to IssueDetail.
- [ ] #7 The operation name is update issue; validation, authentication, remote, and internal failures use the existing apperror boundary without HTTP status classification in CLI.
- [ ] #8 Update output lists Changed fields in fixed title, body order and preserves existing list, filter, inspect, and create output byte compatibility.
- [ ] #9 The implementation preserves Interface → Application Port → Infrastructure adapter direction, explicit dependency injection, and adds no runtime type assertions.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design:
- CLI: `fj issue update OWNER/NAME NUMBER [--title TITLE] [--body BODY]` with `--instance PROFILE`; issue number is positive and at least one of --title or --body is required.
- Scope: title and body only. State, assignee, labels, milestone, comments, metadata, and JSON are out of scope.
- Application: add UpdateRequest{Owner, Name, Number, Title *string, Body *string} and Updater in internal/application/issue. nil means unspecified; ptr("") means update to empty. Repository ports remain unchanged.
- Infrastructure: PATCH /api/v1/repos/{owner}/{repo}/issues/{index}; send only specified title/body fields; safely encode path segments; convert DTO to IssueDetail; operation is `update issue`; use existing apperror boundary.
- Presenter: add update-specific output with Changed fields ordered title, body. Existing list, filter, inspect, and create output remains unchanged.
- Architecture: preserve Interface → Application Port → Infrastructure adapter, explicit DI, no new runtime type assertions, and no HTTP status classification in CLI.

Independent Review: Critical: none. Major: none. Minor: HTTP error boundary test expansion, JSON decode failure test expansion, path encoding test expansion, title/body combined JSON and Changed fields boundary test expansion, pointer semantics boundary test expansion, secret/raw cause redaction test expansion, Composition Root Updater injection test expansion, and CLI-to-Updater delegation test expansion remain possible. Suggestion: reconsider the jsonTransport type assertion when the shared transport boundary is revisited.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented fj issue update OWNER/NAME NUMBER with Application-owned Updater and pointer semantics, explicit-field PATCH requests, safe error translation, explicit dependency injection, and fixed Changed fields output while preserving existing list, filter, inspect, and create output. Independent review found no Critical or Major issues; remaining Minor items are boundary-test expansion opportunities and a future transport-boundary Suggestion.
<!-- SECTION:FINAL_SUMMARY:END -->
