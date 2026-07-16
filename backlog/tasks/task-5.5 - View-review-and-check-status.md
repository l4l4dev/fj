---
id: TASK-5.5
title: View review and check status
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-16 09:29'
labels: []
milestone: m-4
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
modified_files:
  - internal/application/pullrequest/pullrequest.go
  - internal/application/pullrequest/status.go
  - internal/application/pullrequest/status_test.go
  - internal/infrastructure/pullrequest/rest.go
  - internal/infrastructure/pullrequest/status_test.go
  - internal/interface/cli/pullrequest.go
  - internal/interface/cli/pullrequest_presenter.go
  - internal/interface/cli/pullrequest_status_test.go
  - internal/interface/cli/root.go
parent_task_id: TASK-5
priority: medium
ordinal: 20020
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Present review state, checks, and merge readiness clearly.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Pending, successful, and failed states are distinguishable
- [x] #2 Unavailable status is not presented as success
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add Application-owned status models, aggregate rules, viewer port, and validation use case.
2. Read Forgejo pull detail, reviews, and head commit status through the existing REST transport with component-level unavailable handling.
3. Add the non-interactive `pr status` command, minimal Presenter output, and composition wiring without changing `pr inspect`.
4. Add focused Application, Infrastructure, CLI, and Presenter coverage; run the required verification.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Human Decision — Approved MVP Contract:
- Add `fj pr status OWNER/NAME NUMBER` as a new command; preserve the existing `pr inspect` contract unchanged.
- Output only the PR number, review aggregate, check aggregate, and Forgejo-provided mergeable value.
- Review aggregate states are `failed` when any effective non-dismissed, non-stale review requests changes; otherwise `pending` when a requested reviewer remains; otherwise `success` when at least one effective approval exists; otherwise `unavailable`. Unknown review states or unavailable review data are `unavailable` and never success.
- Check aggregate states are `failed` when any status is failure, error, or warning; otherwise `pending` when any status is pending; otherwise `success` only when at least one status exists and every status is success; otherwise `unavailable`. Unknown states or unavailable check data are `unavailable` and never success.
- Display Forgejo `mergeable` only as `yes`, `no`, or `unavailable`; do not infer, interpret, or compute merge readiness.
- If an API or field is unavailable, mark only the affected component `unavailable` where safe to do so.

Deferred / Out of Scope:
- Review lists
- Check lists
- JSON output
- Branch-protection interpretation
- Refresh or polling
- Cross-version interpretation
- Independent merge-readiness calculation

Implementation Authorization:
- A later implementation phase may proceed only within this approved MVP contract. Acceptance Criteria completion, Final Summary, Finalization, commit, and push remain separately gated.

Model: GPT-5 — recording an approved public CLI status contract and conservative unavailable-state rules.

Pre-Implementation Check:
- The approved MVP remains within the M4 pull-request scope and TASK-14 Forgejo-client boundary.
- Existing Application/Infrastructure/CLI/Presenter dependency directions are preserved.
- The command is read-only, introduces no dependency, and does not alter `pr inspect`.
- Authentication failures and unsafe remote failures continue through the existing error boundary; missing component APIs or fields degrade conservatively to `unavailable`.
- No additional Human Decision is required before implementation.

Implementation:
- Added read-only `fj pr status OWNER/NAME NUMBER` without changing `pr inspect`.
- Added Application-owned status request/result types, conservative review/check aggregate states, mergeable tri-state, viewer port, and validation use case.
- Added Forgejo REST reads for pull detail, paginated reviews, and combined head-commit status. Missing component endpoints or fields produce `unavailable`; authentication and other remote failures retain the existing safe error boundary.
- Added minimal Presenter output for PR number, review aggregate, check aggregate, and Forgejo-provided mergeable value only.
- Added Application, Infrastructure aggregation/API, CLI, Presenter, unavailable-state, and regression coverage.

Verification:
- `go test ./...`: PASS
- `go run ./cmd/fj pr status --help`: PASS
- `make pre-commit`: PASS (`git diff --check`, `go vet ./...`, `go test ./...`).

Independent Review: pending.

Independent Review and Remediation:
- Initial review found two Major issues: aggregation ownership in Infrastructure and full-history review aggregation.
- Remediation moved Review / Check aggregation into the Application use case; Infrastructure now performs API access, pagination, and DTO translation only.
- Review aggregation deterministically selects the latest effective review per reviewer by review ID. REQUEST_CHANGES → APPROVED resolves to success; APPROVED → REQUEST_CHANGES resolves to failed; different reviewers remain independently effective.
- Focused tests cover pending, success, failed, empty, unknown, unavailable, dismissed/stale reviews, and same-reviewer transitions.
- Verification passed: focused package tests, git diff --check, go vet ./..., go test ./..., and make pre-commit.
- Independent Re-Review (GPT-5): previous Major findings resolved; new Critical none, Major none, Minor none; AC #1 Pass; AC #2 Pass; Ready for Finalization. This supersedes the earlier pending review marker.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj pr status OWNER/NAME NUMBER with PR number, Application-owned review/check aggregates, and Forgejo-provided mergeable state. Review aggregation uses each reviewer's latest effective review; unavailable data remains unavailable and is never reported as success. Focused tests and make pre-commit passed, and Independent Re-Review found no remaining findings.
<!-- SECTION:FINAL_SUMMARY:END -->
