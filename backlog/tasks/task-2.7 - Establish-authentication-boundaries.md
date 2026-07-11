---
id: TASK-2.7
title: Establish authentication boundaries
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 05:50'
labels: []
dependencies: []
references:
  - ROADMAP.md
modified_files:
  - internal/application/auth/auth.go
  - internal/application/auth/auth_test.go
parent_task_id: TASK-2
ordinal: 21000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Supply credentials to transports without exposing them elsewhere.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Authentication data stays outside domain and presentation models
- [x] #2 Errors and diagnostics redact credential material
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Define an Application-layer credential Provider port and Resolver boundary using configuration credential references.
2. Represent credential material with an opaque type whose diagnostic formatting is always redacted.
3. Sanitize provider failures into a stable application error without preserving credential-bearing messages.
4. Add focused tests for credential delivery, provider input, diagnostic redaction, and error redaction.
5. Run formatting and all Go tests.
6. Check all acceptance criteria and finalize TASK-2.7 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Authentication storage and HTTP transport implementations are outside TASK-2.7. This task defines only the Application-layer port and safe credential/error boundary; domain and presentation packages remain unchanged.

Added an Application-layer Provider port and Resolver for supplying credentials from configuration references.
Credential material is held in an opaque Credential type; standard and Go-syntax diagnostics render only [REDACTED].
Provider failures are converted to the stable ErrCredentialUnavailable error without retaining provider messages that may contain credential material.
Domain and presentation packages were not changed.
Validation passed: gofmt completed and go test ./... passed.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Established the minimal Application-layer authentication boundary with an opaque redacting Credential, provider port, resolver, and sanitized failure contract. Added focused delivery and redaction tests; formatting and all Go tests pass.
<!-- SECTION:FINAL_SUMMARY:END -->
