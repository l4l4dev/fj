---
id: TASK-4.7
title: Manage issue comments
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-3
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-4
ordinal: 36000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Add and inspect discussion attached to an issue.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Comments remain associated with the correct issue
- [x] #2 Empty content and permission failures are handled clearly
- [x] #3 The CLI provides fj issue comment list OWNER/NAME NUMBER and fj issue comment add OWNER/NAME NUMBER --body BODY, both with --instance; issue number is a positive integer.
- [x] #4 Application owns Comment{ID int64, Body string}, ListCommentsRequest, AddCommentRequest, CommentViewer, and CommentCreator in internal/application/issue; Repository ports remain unchanged.
- [x] #5 The Use Cases validate owner, name, and positive issue number; add rejects missing, empty, and whitespace-only body.
- [x] #6 Infrastructure calls GET and POST /api/v1/repos/{owner}/{repo}/issues/{index}/comments, sends body-only JSON for POST, and converts private DTOs to Comment.
- [x] #7 The operations are list issue comments and add issue comment; failures use the existing apperror boundary without exposing HTTP status, credentials, URL details, response bodies, or raw causes.
- [x] #8 List output is Comments: followed by - #<id> <body>; empty output includes No comments found.; add output is Comment: followed by #<id> <body>.
- [x] #9 Comment update/delete, author, timestamp, reaction, markdown rendering, JSON, pagination, and fetch-all remain out of scope; existing issue command output remains unchanged.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design:
- CLI: `fj issue comment list OWNER/NAME NUMBER` and `fj issue comment add OWNER/NAME NUMBER --body BODY`, both with `--instance PROFILE`; issue number is a positive integer.
- Scope: comment list and comment add only. Comment update/delete, author, timestamp, reaction, markdown rendering, JSON, pagination, and fetch-all are out of scope.
- Application: Comment{ID int64, Body string}, ListCommentsRequest, AddCommentRequest, CommentViewer, and CommentCreator in internal/application/issue. Repository ports remain unchanged.
- Infrastructure: GET and POST /api/v1/repos/{owner}/{repo}/issues/{index}/comments. POST JSON contains body only. Use existing apperror boundary; operations are `list issue comments` and `add issue comment`.
- Validation: owner/name, positive issue number, and non-empty non-whitespace add body.
- Presenter: list uses Comments: with `- #<id> <body>` and empty results use `Comments:` plus `No comments found.`; add uses `Comment:` and `#<id> <body>`. Existing issue command output remains unchanged.

Independent Review:
Critical: none
Major: none
Minor:
- HTTP error boundary test expansion opportunity
- JSON decode failure test expansion opportunity
- path encoding test expansion opportunity
- empty response Presenter test expansion opportunity
- secret redaction test expansion opportunity
- Composition Root DI test expansion opportunity
- CLI delegation test expansion opportunity
- regression test expansion opportunity
Suggestion:
- Document multiline comment body human-readable output policy in a future task
- Consider boundary tests for abnormal Comment IDs and empty response handling
- Consider future cleanup of the jsonTransport type assertion

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
TASK-4.7 implementation completed: added issue comment list and add workflows with Application-owned comment ports, safe Forgejo REST adapters, explicit dependency injection, validation, and Presenter output. Existing issue command output and out-of-scope behavior remain unchanged. Independent review found no Critical or Major issues; remaining Minor items and Suggestions are recorded.
<!-- SECTION:FINAL_SUMMARY:END -->
