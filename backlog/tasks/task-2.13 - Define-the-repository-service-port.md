---
id: TASK-2.13
title: Define the repository service port
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 17:08'
updated_date: '2026-07-11 05:50'
labels: []
dependencies: []
references:
  - ROADMAP.md
modified_files:
  - internal/application/repository/repository.go
  - internal/application/repository/repository_test.go
parent_task_id: TASK-2
ordinal: 79000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define the Application-layer Repository Service Port required by repository use cases. Keep the port transport-independent and limit this Task to the service contract.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The port expresses repository listing capability without depending on net/http, TOML, Cobra, or Forgejo SDK types.
- [x] #2 Repository models and remote failure outcomes crossing the boundary are explicit and testable.
- [x] #3 Repository list use-case orchestration is explicitly outside this Task and belongs to TASK-3.1.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add the Application-owned Repository model and ListRequest with Page and Limit fields.
2. Define the transport-independent Service interface with List(context.Context, ListRequest) ([]Repository, error).
3. Define a safe Application RemoteError carrying only operation/status classification, with no credential, URL, response body, or raw transport cause.
4. Add deterministic contract tests for the model, Service implementations, context propagation, and RemoteError safety; leave validation, use-case orchestration, and Infrastructure translation to their approved Tasks.
5. Run required checks and obtain an independent GPT-5 post-implementation review, then finalize TASK-2.13.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — introducing the Application Repository Service Port and boundary models is architecture-sensitive and requires high-capability design and verification

Pre-implementation check (GPT-5): Critical none. Major: Repository model ownership, Service API shape, page/limit validation ownership, and Application remote-error type/shape require human approval. Minor: naming details. Suggestion: add compile-time interface assertion in TASK-2.14 adapter tests. No code implementation started.

Human approved the Major Change design: Application owns Repository, Service uses List(context.Context, ListRequest) ([]Repository, error), ListRequest has Page/Limit, validation and orchestration belong to TASK-3.1, Application RemoteError is safe operation/status only, Infrastructure translation belongs to TASK-2.14, context.Context is required, and TASK-2.14 will add the compile-time assertion.

Implemented the approved Application repository boundary: Repository {Owner, Name}, ListRequest {Page, Limit}, and Service.List(context.Context, ListRequest) ([]Repository, error). Added safe Application RemoteError with only operation/status classification and no credentials, URLs, response bodies, raw cause, or Unwrap. Validation and use-case orchestration remain in TASK-3.1; Infrastructure translation and compile-time adapter assertion remain in TASK-2.14.

Validation passed: gofmt -l . (no output); git diff --check; focused go test ./internal/application/repository; go vet ./...; go test ./.... Initial focused run hit the known Go build-cache permission restriction and passed on approved elevated rerun.

Post-implementation review (independent GPT-5): Critical none; Major none; Minor none; Suggestion: future work may constrain NewRemoteError operation values to developer-controlled fixed strings. No fix adopted because the current approved scope requires only safe operation/status classification and no sensitive fields.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Defined the Application-owned Repository model, pagination request, transport-independent Service port, and safe RemoteError classification. The boundary uses context.Context, has no transport or presentation dependencies, excludes sensitive/raw failure data, and leaves validation/orchestration and Infrastructure translation to their approved Tasks. All checks and the independent Major Change review passed.
<!-- SECTION:FINAL_SUMMARY:END -->
