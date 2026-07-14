---
id: TASK-5.3
title: Create a pull request
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-14 01:18'
labels: []
milestone: m-4
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-5
priority: medium
ordinal: 20010
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Open a pull request with an explicit source and target.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Branch direction and repository are visible before submission
- [x] #2 Invalid or conflicting branches are reported clearly
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add the application create request, creator port, and validation use case for the approved same-repository MVP.
2. Add the Forgejo REST POST adapter with safe status-to-error mapping.
3. Add the non-interactive `pr create` command, success presenter, dependency wiring, and focused tests.
4. Verify the complete Go test suite, command help, and Git diff hygiene.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Human Decision — Approved MVP Contract:
- Purpose: create a pull request from a head branch to a base branch within the same repository.
- Required inputs: repository target, title, base branch, and head branch.
- Success output: created pull request number, title, head branch, and base branch.
- Execution is non-interactive. The MVP sends no body and does not support cross-fork or draft pull requests, labels, assignees, reviewers, milestones, or template support.
- Errors cover invalid repository targets; missing or invalid title, head, or base; identical head and base; duplicate or conflicting pull requests; authentication or authorization failures; and Forgejo API failures.
- Errors use the existing fj error boundary and must not expose raw HTTP responses, credentials, Authorization headers, or stack traces.

Deferred Items:
- Optional body
- Cross-fork pull requests
- Draft pull requests
- maintainer_can_modify
- Labels
- Assignees
- Reviewers
- Milestone
- Template support

Implementation Authorization:
- Implementation may proceed only within this approved MVP contract. Acceptance Criteria completion, Final Summary, Finalization, commit, and push remain separate later stages.

Model: GPT-5 — recording an approved public CLI contract with strict scope boundaries.

Implementation:
- Added non-interactive `fj pr create OWNER/NAME --title TITLE --head HEAD --base BASE` for same-repository pull requests.
- The request sends only `title`, `head`, and `base`; no deferred metadata is supported.
- Added application validation, REST creation, safe authentication/not-found/conflict/remote mapping, success output, dependency composition, and focused tests.
- Cross-fork head syntax is rejected locally.

Verification:
- `go test ./...`: PASS
- `go run ./cmd/fj pr create --help`: PASS
- `git diff --check`: PASS

Independent Review: pending.

Implementation verification completed:
- `make pre-commit`: PASS (`git diff --check`, `go vet ./...`, `go test ./...`).
- Implementation is complete and remains pending Independent Review.

Remediation — MAJ-1:
- Changed pull-request creation 404 handling to the safe, non-diagnostic message `repository or branch not found` while preserving the NotFound category and existing error boundary.
- Added infrastructure coverage for the PR creation 404 mapping.
- Verification: `go test ./internal/infrastructure/pullrequest/...`, `go test ./...`, and `git diff --check`: PASS.
- Independent Re-Review remains pending.

Evidence Summary:
- The approved same-repository, non-interactive `pr create` contract is implemented across Application, Forgejo REST adapter, CLI, Presenter, and composition boundaries.
- Required repository, title, head, and base inputs are validated; success output identifies the created PR and branch direction.
- Safe error mapping covers authentication, not-found, conflict, validation, and remote failures without exposing sensitive transport details.
- `go test ./...`, `go vet ./...`, command help verification, and `git diff --check` passed.

Independent Review and Re-Review:
- Initial review: Critical none; MAJ-1 identified ambiguous PR-creation 404 diagnosis; four non-blocking Minor findings recorded.
- Remediation changed the 404 message to `repository or branch not found` and added focused infrastructure coverage.
- Re-Review: MAJ-1 Resolved; AC #1 and #2 Pass; Critical none; Major none; Finalization Decision Ready for Finalization.

Deferred Items — non-blocking Minor findings:
- Consider distinguishing HTTP 409 from 422 when a stable Forgejo contract supports clearer categorization.
- Consider reducing duplicated CLI and Application validation without weakening boundary ownership.
- Consider simplifying the pull-request command constructor composition in a separate refactoring.
- Add explicit create-response decode-failure coverage during future adapter hardening.
- These items are outside TASK-5.3 completion scope.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented the approved same-repository `fj pr create` MVP with explicit title, head, and base inputs, deterministic success output, Application validation, Forgejo REST creation, safe error mapping, dependency wiring, and focused tests. Independent Review MAJ-1 was remediated and resolved; all Acceptance Criteria and verification pass. Four non-blocking Minor improvements remain explicitly deferred.
<!-- SECTION:FINAL_SUMMARY:END -->
