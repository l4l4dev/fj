---
id: TASK-3.1
title: List repositories
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 05:50'
labels: []
dependencies:
  - TASK-2.9
  - TASK-2.10
  - TASK-2.11
  - TASK-2.12
  - TASK-2.13
  - TASK-2.14
references:
  - ROADMAP.md
modified_files:
  - internal/application/repository/list.go
  - internal/application/repository/list_test.go
  - internal/interface/cli/root.go
  - internal/interface/cli/repository.go
  - internal/interface/cli/repository_test.go
parent_task_id: TASK-3
ordinal: 24000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Allow users to discover repositories on the selected Forgejo instance.

Scope: Implement only the repository-listing use case, its CLI command, the Composition Root wiring, and human-readable output.

Out of scope: XDG/TOML configuration loading, environment credential resolution, the authenticated HTTP client, the Repository Service Port, the REST adapter, JSON output, and fetch-all pagination behavior.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The CLI Composition Root wires the selected instance and dependencies, renders an explicit human-readable empty result, and presents remote failures through the common CLI error behavior.
- [x] #2 The use case requests a page and limit through the Repository Service Port and presents each returned repository with an unambiguous owner/name identity.
- [x] #3 The repository list command uses the approved command contract and explicit human-readable output.
- [x] #4 The command accepts page and limit flags with defaults page=1 and limit=30, and rejects non-positive values before the remote request.
- [x] #5 An empty repository result is rendered explicitly, and each repository identifies owner/name.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Pending human approval of the Major Change CLI contract:
1. Add an Application repository-list use case that validates page/limit and calls repository.Service.List.
2. Add the approved repo list CLI command, flags, and human-readable presenter.
3. Compose configuration loading, instance selection, credential resolution, Forgejo HTTP client, REST adapter, and use case without leaking endpoints or credentials.
4. Present an explicit empty result and map remote failures through the existing CLI error categories; do not add JSON or fetch-all behavior.
5. Add focused use-case/CLI/composition tests, run required checks, and obtain an independent GPT-5 post-implementation review.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Assessment only; implementation intentionally not started. Model: GPT-5 — determining whether this repository-listing Task can fit the approved architecture requires cross-layer boundary analysis.

Blocker: internal/infrastructure and internal/domain directories are empty; no Forgejo REST adapter, repository service port, repository application use case, CLI command, or composition path exists. Implementing TASK-3.1 now would require introducing multiple architectural boundaries and expanding beyond the Task's stated scope. Per the human stop condition, keep TASK-3.1 To Do and request Task decomposition or a separately approved foundation Task.

Baseline validation without code changes: gofmt -l . (no output); git diff --check passed; go vet ./... passed; go test ./... passed. No files were changed by this assessment.

Model: GPT-5 — TASK-3.1 introduces a public CLI command, flags, output contract, and composition boundary, requiring high-capability Major Change design and verification

Pre-implementation check (GPT-5): Critical none. Major: command name/flags, page/limit defaults and bounds, empty output wording, and instance display are compatibility-sensitive and require human approval. Minor none. Suggestion: recommended repo list command with page=1, limit=30, positive validation, owner/name output, safe remote category, no JSON/fetch-all. No code implementation started.

Status synchronized to To Do: implementation has not started. Human approval of the CLI contract remains pending; the approved contract is repo list with page/limit flags, page=1 and limit=30 defaults, positive-value validation, explicit empty-result output, and owner/name identity.

Pre-implementation check (GPT-5): Critical none, Major none, Minor none. Approved CLI contract satisfies the Major Change gate. Architecture boundaries remain Interface -> Application -> repository.Service -> Infrastructure; existing credential and HTTP safety boundaries are reused. Suggestion: add Application, CLI, and Composition boundary tests.

Model: GPT-5 — Major Change CLI and Composition Root implementation requiring architecture-aware compatibility and security review.

Validation: make pre-commit passed, including git diff --check, go vet ./..., and go test ./.... Initial sandboxed go test was blocked by Go cache permissions; escalated rerun passed.

Post-implementation review (GPT-5): Critical none, Major none, Minor none. Suggestion: add an explicit --instance Composition Root integration test in future; not required for this approved scope.

Review correction: Major inconsistency identified in repository.RemoteError classification. HTTP 401 and 403 now map to categoryAuthentication; all other repository RemoteError statuses remain categoryRemote. Operation names and safe messages are unchanged.

Validation update: added CLI tests for 401, 403, and 503 classification. make pre-commit passed (git diff --check, go vet ./..., go test ./...).

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Corrected TASK-3.1 remote error classification so authentication-related HTTP statuses 401 and 403 produce categoryAuthentication while other remote statuses produce categoryRemote. Added focused CLI coverage and revalidated the full repository.
<!-- SECTION:FINAL_SUMMARY:END -->
