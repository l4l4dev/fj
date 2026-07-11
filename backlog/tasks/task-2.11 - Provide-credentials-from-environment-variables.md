---
id: TASK-2.11
title: Provide credentials from environment variables
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 17:08'
updated_date: '2026-07-11 05:50'
labels: []
dependencies:
  - TASK-2.7
references:
  - ROADMAP.md
modified_files:
  - internal/infrastructure/auth/environment.go
  - internal/infrastructure/auth/environment_test.go
parent_task_id: TASK-2
ordinal: 77000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Resolve configured credential references from environment variables only, without adding a credential store or embedding secret values in ordinary configuration output or errors.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 A credential reference resolves from its documented environment variable.
- [x] #2 Missing environment variables produce an authentication failure without exposing secret material.
- [x] #3 Credential values are not included in logs, diagnostics, or rendered configuration.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Implement the approved Infrastructure-layer environment Provider satisfying the existing auth.Provider interface.
2. Treat CredentialReference as the exact environment variable name; reject empty references and treat unset or empty variables as credential unavailable.
3. Return only the existing sanitized auth.ErrCredentialUnavailable-wrapped failure; never include the reference or credential value in errors, logs, or output.
4. Inject the provider through auth.NewResolver so Application code has no os.Getenv/os.LookupEnv dependency and future Keychain providers remain interchangeable.
5. Add focused success, missing, empty, redaction, and Application-import-boundary tests; run required checks and obtain an independent GPT-5 post-implementation review, then finalize TASK-2.11.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — credential acquisition changes the authentication/security boundary and requires high-capability design and verification

Pre-implementation check (GPT-5): Critical none. Major: exact CredentialReference-to-environment-variable mapping and empty-value semantics are user-facing configuration/security decisions requiring human approval. Minor none. Suggestion: Infrastructure env Provider implementing existing auth.Provider, exact reference as env name, unset/empty as sanitized ErrCredentialUnavailable, injected via auth.NewResolver, with success/missing/empty/redaction and import-boundary tests. No code implementation started.

Human approved the Major Change detail: CredentialReference is used as the exact environment variable name; empty references and unset/empty variables produce the same sanitized credential-unavailable outcome. The Infrastructure provider implements auth.Provider, Application remains free of os environment access, and future Keychain providers can replace it through the same interface.

Implemented the approved Infrastructure EnvironmentProvider using the existing application auth.Provider boundary. CredentialReference is used as the exact environment variable name; empty references, unset variables, and empty values return the sanitized auth.ErrCredentialUnavailable outcome. Application code does not access os environment APIs, and the provider is injected through auth.NewResolver so future Keychain providers can implement the same interface. Credential values and references are absent from errors and diagnostics.

Validation passed: gofmt -l . (no output); git diff --check; focused go test ./internal/infrastructure/auth; go vet ./...; go test ./.... The first focused test run exposed only a test assertion issue for an empty reference (empty text matches every string); the assertion was corrected and the focused and full suites then passed with approved elevated cache access.

Post-implementation review (independent GPT-5): Critical none; Major none; Minor none; Suggestion none. Review confirmed approved scope only, all acceptance criteria, Architecture dependency direction, existing Provider/Resolver use, no Application os dependency, full redaction, Keychain substitutability, sufficient tests, and accurate Backlog records.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added an Infrastructure EnvironmentProvider that resolves exact CredentialReference environment variables and delegates through the existing auth.Provider/Resolver boundary. Empty, missing, and empty-value credentials fail safely without exposing references or credential values; Application remains independent of os and future Keychain providers remain interchangeable. All checks and the independent Major Change review passed.
<!-- SECTION:FINAL_SUMMARY:END -->
