---
id: TASK-2.12
title: Establish the authenticated Forgejo HTTP client
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 17:08'
updated_date: '2026-07-10 17:35'
labels: []
dependencies:
  - TASK-2.10
  - TASK-2.11
references:
  - ROADMAP.md
modified_files:
  - internal/infrastructure/forgejo/client.go
  - internal/infrastructure/forgejo/client_test.go
parent_task_id: TASK-2
ordinal: 78000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the infrastructure boundary for authenticated Forgejo HTTP requests using net/http only. Compose the selected instance endpoint and environment-provided credential without introducing a Forgejo SDK.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Requests use net/http and attach the resolved credential without exposing it in diagnostics.
- [x] #2 The client targets the selected instance endpoint and preserves remote failures for the Interface error presentation boundary.
- [x] #3 Every request sends User-Agent in the form fj/<version>.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Define a private Infrastructure HTTP client using net/http with an injectable Doer/http.Client boundary and a default timeout.
2. Accept the selected Application config.Instance, resolved opaque auth.Credential, and version; safely join the Endpoint and API path and send Authorization: token <credential> plus User-Agent fj/<version>.
3. Prevent credential forwarding to unrelated redirect hosts and never include credentials in URL, query, errors, logs, or response-body diagnostics.
4. Return typed safe remote errors for transport and non-2xx failures while preserving status context without response bodies.
5. Add deterministic httptest and injected-Doer tests for request construction, slash handling, headers, timeout/redirect policy, and safe failures.
6. Run required checks and obtain an independent GPT-5 post-implementation review, then finalize TASK-2.12.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — authenticated HTTP transport changes the Infrastructure/auth/security boundary and requires high-capability design and verification

Pre-implementation check (GPT-5): Critical none. Major: Authorization scheme, endpoint path joining, remote failure representation/body redaction, injectable HTTP boundary/timeout, redirect credential policy, and version source for User-Agent fj/<version> require human approval. Minor: response body size/non-2xx details. Suggestion: private client with injectable Doer, Authorization token header, typed safe remote error, and httptest coverage. No code implementation started.

Human approved the Major Change details: Authorization: token <credential>, safe endpoint/path joining, header-only credential use, typed safe remote failures, injectable http.Client, timeout, redirect protection, User-Agent fj/<version>, constructor version argument, unchanged Application boundary using config.Instance plus resolved auth.Credential, and minimal private HTTP client API.

Post-implementation review (independent GPT-5): Critical none; Major none; Minor: completion records were pending at review time; Suggestion: RemoteError.Unwrap() can expose an injected Doer cause programmatically even though Error() is safe. Per AGENTS.md review classification, Suggestion fixes require human decision; task remains In Progress pending a decision to remove raw cause/Unwrap or retain it as an internal diagnostic boundary.

Human approved the post-review Suggestion: remove raw transport cause storage and Unwrap() from RemoteError. Preserve errors.As classification and safe operation/status access only; update related tests and finalize this Task.

Implemented the approved private Infrastructure HTTP client using net/http. It accepts config.Instance, resolved opaque auth.Credential, and constructor version; safely joins Endpoint and API path with url.JoinPath; sends Authorization: token <credential> and User-Agent fj/<version>; uses an injectable Doer and a 30-second default timeout; blocks redirects across host or scheme; and never places credentials in URL/query, errors, logs, or response-body diagnostics. RemoteError now stores only safe operation/status classification, with no raw cause or Unwrap; errors.As plus Operation() and StatusCode() preserve safe classification.

Validation passed: gofmt -l . (no output); git diff --check; focused go test ./internal/infrastructure/forgejo; go vet ./...; go test ./.... Initial focused run hit the known Go build-cache permission restriction and passed on approved elevated rerun. A compile error during the approved raw-cause removal was corrected before the final successful runs.

Final post-implementation review (independent GPT-5): Critical none; Major none; Minor none; Suggestion none. Confirmed Credential, URL sensitive portions, and injected Doer raw errors cannot be retrieved from RemoteError; errors.As classification and safe operation/status access remain; approved scope and architecture boundaries are preserved; and all verification passes.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Established the authenticated Forgejo HTTP client with safe endpoint/path joining, token Authorization, fj/<version> User-Agent, injectable net/http transport, timeout, redirect protection, and typed remote failures. RemoteError retains only safe operation/status data and exposes no raw causes, credentials, URLs, or response bodies. All acceptance criteria, full verification, and the final independent Major Change review passed.
<!-- SECTION:FINAL_SUMMARY:END -->
