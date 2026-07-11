---
id: TASK-4.4
title: Create an issue
status: Done
assignee: []
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 02:20'
labels: []
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-4
ordinal: 33000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Create an issue with explicit content and metadata.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The target repository and submitted fields are observable
- [ ] #2 Invalid metadata is rejected before submission
- [ ] #3 The CLI provides fj issue create OWNER/NAME with required --title, optional --body, and --instance; assignee, label, milestone, and state flags are not accepted.
- [ ] #4 Application owns CreateRequest and Creator in internal/application/issue; Create returns IssueDetail and existing Repository ports are unchanged.
- [ ] #5 The Use Case validates owner, name, and title, rejects whitespace-only title, and permits an empty body.
- [ ] #6 Infrastructure sends POST /api/v1/repos/{owner}/{repo}/issues with title and body JSON fields only, safely encodes path segments, and converts the response to IssueDetail.
- [ ] #7 The operation name is create issue; failures use the existing apperror boundary without exposing credentials, URL details, response bodies, or raw transport causes.
- [ ] #8 Successful creation reuses the Inspect Presenter format, displays `-` for an empty Body, and emits no create-specific success message.
- [ ] #9 Assignee, labels, milestone, state, confirmation, interactive mode, and JSON output remain out of scope.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design:
- CLI: `fj issue create OWNER/NAME --title "TITLE" [--body "BODY"]` with `--instance`; assignee, label, milestone, and state are out of scope.
- Application: add CreateRequest and Creator in internal/application/issue; creation returns IssueDetail; Repository ports remain unchanged.
- Infrastructure: POST /api/v1/repos/{owner}/{repo}/issues with JSON body containing title and body only; safely encode path segments; use the existing apperror boundary; operation is `create issue`; credentials, URL details, response bodies, and raw causes remain private.
- Presenter: reuse the Inspect Presenter format; empty Body displays `-`; no create-specific success message.
- Validation: validate owner, name, and title; reject whitespace-only title; allow empty body.
- Out of scope: assignee, labels, milestone, state, confirmation, interactive mode, and JSON.

Independent Review: Critical: none. Major: none. Minor: HTTP error boundary test expansion, JSON decode failure test expansion, path encoding test expansion, empty-body JSON submission verification, secret/raw cause redaction test expansion, and Composition Root Creator injection test expansion remain possible. Suggestion: reconsider the jsonTransport type assertion when the shared transport boundary is revisited; consider adding a CLI-to-Creator delegation test.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented fj issue create OWNER/NAME with title/body validation, Application-owned Creator Port, safe Forgejo POST translation, explicit Creator dependency injection, and Inspect Presenter-compatible output. Independent review found no Critical or Major issues; remaining Minor items are boundary-test expansion opportunities and future transport/delegation Suggestions.
<!-- SECTION:FINAL_SUMMARY:END -->
