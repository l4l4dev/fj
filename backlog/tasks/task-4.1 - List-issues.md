---
id: TASK-4.1
title: List issues
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 05:50'
labels: []
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-4
ordinal: 30000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Discover issues in the selected repository.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Results include stable issue identity and lifecycle state
- [x] #2 Empty collections and pagination state are clear
- [x] #3 The CLI provides fj issue list OWNER/NAME with --page (default 1), --limit (default 30), and --state open|closed|all (default open); --all, fetch-all, and additional filters are out of scope.
- [x] #4 Application owns a minimal Issue model with Number, Title, and State (StateOpen or StateClosed) in internal/application/issue, plus ListRequest, Page, and Lister; existing repository ports are unchanged.
- [x] #5 The Use Case validates OWNER/NAME, positive page/limit, and state, then delegates to Lister; Infrastructure uses GET /api/v1/repos/{owner}/{repo}/issues with page, limit, and state query parameters and safely encodes path segments.
- [x] #6 MorePages is calculated as len(issues) == limit without Link headers, total-count APIs, or Forgejo-specific pagination metadata in Application.
- [x] #7 Human output lists Issues as - #<number> <title> [<state>] followed by Page, Limit, and More pages; empty output includes Issues:, No issues found., and pagination metadata.
- [x] #8 The operation name is list issues; HTTP, transport, and JSON failures use the Application error boundary without exposing credentials, URL details, response bodies, or raw causes; JSON output and detailed issue metadata are out of scope.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add the Application issue model, pagination contract, and validated list use case. 2. Add a Forgejo issue REST adapter with safe DTO conversion and application error mapping. 3. Add explicit CLI wiring and a presenter for the approved human-readable output. 4. Add focused tests and run the full verification suite.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design decisions:

1. CLI command
- `fj issue list OWNER/NAME`

2. Pagination
- Flags: `--page`, `--limit`
- Defaults: page=1, limit=30
- `--all` and fetch-all are out of scope.

3. State filter
- `--state open|closed|all`
- Default: `open`
- Additional filters such as assignee, label, and milestone are out of scope.

4. Human-readable output
- Issues: `- #<number> <title> [<state>]`
- Pagination fields: `Page`, `Limit`, `More pages`
- Author, created/updated timestamps, labels, assignee, and milestone are out of scope.

5. Pagination metadata
- Calculate More pages as `len(issues) == limit`.
- Do not use Link headers, total-count APIs, or Forgejo-specific pagination metadata in Application.

Approved design decisions:

- Command: `fj issue list OWNER/NAME`.
- Pagination: `--page` default 1 and `--limit` default 30; no --all/fetch-all.
- State: `--state open|closed|all`, default open; no assignee, label, milestone, or other filters.
- Application model: minimal Issue{Number, Title, State} with StateOpen/StateClosed in internal/application/issue.
- Port: ListRequest, Page, and Lister in internal/application/issue; Use Case validates target/page/limit/state; Infrastructure handles Forgejo access, DTO conversion, and pagination metadata; Interface handles CLI and Presenter. Existing repository ports remain unchanged.
- API: GET /api/v1/repos/{owner}/{repo}/issues with page, limit, and state; path segments are safely encoded; MorePages is len(issues) == limit; no Link header, total count, or fetch-all.
- Output: issue lines `- #<number> <title> [<state>]`, then Page, Limit, More pages; empty output includes `Issues:` and `No issues found.` with pagination metadata.
- Operation: `list issues`.
- JSON and detailed issue metadata are out of scope.

Implemented the approved issue listing flow with an Application-owned issue model and Lister port, Forgejo REST adapter, explicit CLI dependency wiring, and a human-readable presenter. Validation: gofmt -l ., git diff --check, go vet ./..., go test ./..., and make pre-commit all passed (GOCACHE=/tmp/fj-gocache used for sandbox compatibility).

Independent Review: Critical: none. Major: none. Minor: boundary test coverage can be expanded for HTTP status mapping, JSON decode failures, path encoding, empty output, pagination, and secret redaction. Suggestion: add focused adapter, presenter, and Composition Root boundary tests; clarify unknown Forgejo state handling in a future improvement.

Historical note: This task was completed before the standard workflow was introduced. No Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented the approved fj issue list OWNER/NAME flow with Application-owned issue contracts, explicit Lister dependency injection, a Forgejo REST adapter, and Presenter-based human-readable output. Independent review found no Critical or Major issues; only Minor boundary-test expansion opportunities and future Suggestions remain. Validation completed successfully.
<!-- SECTION:FINAL_SUMMARY:END -->
