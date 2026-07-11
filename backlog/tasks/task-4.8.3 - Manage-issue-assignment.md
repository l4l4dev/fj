---
id: TASK-4.8.3
title: Manage issue assignment
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 06:11'
updated_date: '2026-07-11 07:09'
labels: []
dependencies:
  - TASK-2.9
parent_task_id: TASK-4.8
ordinal: 88000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Manage issue assignment independently as a focused Issue metadata workflow.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The CLI contract for assigning and unassigning an issue is explicitly defined and supports --instance PROFILE.
- [x] #2 OWNER/NAME and positive issue number validation are enforced before remote access.
- [x] #3 Application defines an assignment model and dedicated Port without changing existing issue or repository Ports.
- [x] #4 Infrastructure uses the approved Forgejo assignment API contract and translates failures through the Application error boundary.
- [x] #5 Assignment changes do not overwrite labels, milestone, title, body, or state.
- [x] #6 Human-readable assignment output is defined and does not change existing issue command output.
- [x] #7 Tests cover validation, API mapping, error boundary, Presenter output, and explicit Composition Root injection.
- [x] #8 The CLI provides fj issue assign OWNER/NAME NUMBER USER and fj issue unassign OWNER/NAME NUMBER with --instance; USER is a non-empty username and none/empty cannot be used for clearing.
- [x] #9 The scope is single-user assign/unassign only; multiple assignee, team assignment, permission management, user listing, JSON, pagination, and bulk metadata update remain out of scope.
- [x] #10 Application owns Assignment{Username string}, AssignRequest, UnassignRequest, Assigner, and Unassigner in internal/application/issue; Repository ports remain unchanged.
- [x] #11 Assign replaces a different assignee, skips PATCH and succeeds for the same username; unassign skips PATCH and succeeds when no assignee exists; unknown users are remote errors.
- [x] #12 Infrastructure uses PATCH /api/v1/repos/{owner}/{repo}/issues/{index} with assignee username or null, safely encodes paths, and uses a private current-Issue DTO retaining only assignee username for idempotency.
- [x] #13 The operations are get issue assignee, assign issue, and unassign issue; validation, authentication, and remote failures use the existing apperror boundary without exposing HTTP status, credentials, URL details, response bodies, or raw causes.
- [x] #14 Assign output is Issue: #<number> followed by Assignee: <username>; unassign output is Issue: #<number> followed by Assignee cleared; existing issue outputs remain unchanged.
- [x] #15 Explicit dependency injection and Interface → Application Assignment Port → Infrastructure adapter direction are preserved; no runtime type assertions are added and Application remains independent of HTTP, Cobra, and Forgejo DTOs.
<!-- AC:END -->































## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design:
- CLI: `fj issue assign OWNER/NAME NUMBER USER` and `fj issue unassign OWNER/NAME NUMBER`, both with `--instance PROFILE`; issue identity is OWNER/NAME NUMBER. USER is a username; empty username is invalid. none/empty cannot clear; only unassign clears.
- Scope: single-user assign and single-user unassign only. Multiple assignee, team assignment, permission management, user listing, JSON, pagination, and bulk metadata update are out of scope.
- Semantics: assign replaces a different assignee, skips PATCH and succeeds for the same username; unassign skips PATCH and succeeds when no assignee exists. Unknown users are remote errors.
- Application: Assignment{Username string}, AssignRequest, UnassignRequest, Assigner, and Unassigner in internal/application/issue. Repository ports remain unchanged.
- Infrastructure: PATCH /api/v1/repos/{owner}/{repo}/issues/{index}; assign body {"assignee":"<username>"}; unassign body {"assignee":null}. Internal GET of the current Issue is allowed only for idempotency and is not a user-facing inspect command. Private DTO retains assignee username only and converts to Application Assignment.
- Presenter: assign outputs Issue: #<number> and Assignee: <username>; unassign outputs Issue: #<number> and Assignee cleared. Existing issue outputs remain unchanged.
- Error operations: `get issue assignee`, `assign issue`, and `unassign issue`. Categories are validation, authentication, and remote; user/repository/issue not found are remote, 401/403 authentication, and other HTTP/transport/JSON failures remote. HTTP status is not exposed to CLI.
- Architecture: preserve Interface → Application Assignment Port → Infrastructure adapter, unchanged Repository ports, Application independence from HTTP/Cobra/Forgejo DTOs, explicit DI, and no runtime type assertions.

Verification:
- gofmt -l . 成功
- git diff --check 成功
- go vet ./... 成功
- go test ./... 成功
- make pre-commit 成功
<!-- SECTION:NOTES:END -->

## Comments

<!-- COMMENTS:BEGIN -->
author: @codex
created: 2026-07-11 07:09
---
Independent Review

Critical:
なし

Major:
なし

Minor:
境界テスト拡充余地あり

Suggestion:
error boundary、DI、回帰テスト等の追加検討
---
<!-- COMMENTS:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
TASK-4.8.3の実装完了。Issue assignmentのApplication Port、REST adapter、CLI command、Presenter、明示的Dependency Injectionを追加し、承認済みのassign/unassign仕様と既存Issue architectureを維持した。
<!-- SECTION:FINAL_SUMMARY:END -->
