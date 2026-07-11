---
id: TASK-4.3
title: Inspect an issue
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
ordinal: 32000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Display the details and current state of one issue.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Repository and issue identity are explicit
- [x] #2 Missing and inaccessible issues produce distinct errors
- [x] #3 The CLI provides fj issue inspect OWNER/NAME NUMBER and accepts only a positive integer issue number.
- [x] #4 Application owns IssueDetail with Number, Title, State, and Body, plus InspectRequest and an Inspector interface in internal/application/issue; existing issue list and repository ports remain unchanged.
- [x] #5 Infrastructure calls GET /api/v1/repos/{owner}/{repo}/issues/{index}, safely encodes path segments, and converts the private DTO to IssueDetail.
- [x] #6 The operation name is inspect issue and failures use the existing Application error boundary without exposing HTTP status, credentials, URL details, response bodies, or raw causes to the CLI.
- [x] #7 The Presenter adds only issue inspect output; existing issue list output remains byte-compatible.
- [x] #8 Author, timestamps, labels, assignee, milestone, comments, repository metadata, JSON output, and markdown rendering are out of scope.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design:
- Command: `fj issue inspect OWNER/NAME NUMBER`
- Issue identity: positive integer issue number only.
- Application owns IssueDetail{Number, Title, State, Body}, InspectRequest, and Inspector in internal/application/issue.
- API: GET /api/v1/repos/{owner}/{repo}/issues/{index}; path segments are safely encoded and private DTOs convert to IssueDetail.
- Operation: `inspect issue`.
- Existing apperror boundary is used; HTTP status, credentials, URL details, response bodies, and raw causes must not leak to CLI.
- Presenter adds issue inspect output only; existing issue list output remains unchanged.
- Author, timestamps, labels, assignee, milestone, comments, repository metadata, JSON output, and markdown rendering are out of scope.

Independent Review: Critical: none. Major: none. Minor: Infrastructure error boundary test expansion, path encoding test expansion, empty-body Presenter test expansion, and Composition Root injection test expansion remain possible. Suggestion: organize Presenter responsibilities before future markdown rendering and consider adding secret-redaction boundary tests.

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
Implemented fj issue inspect OWNER/NAME NUMBER with Application-owned IssueDetail and Inspector, safe Forgejo REST translation, explicit Composition Root injection, and inspect-only Presenter output while preserving issue list behavior. Independent review found no Critical or Major issues; Minor test expansion opportunities and future Presenter/security Suggestions remain.
<!-- SECTION:FINAL_SUMMARY:END -->
