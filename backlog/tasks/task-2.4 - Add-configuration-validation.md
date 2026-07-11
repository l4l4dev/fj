---
id: TASK-2.4
title: Add configuration validation
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
  - internal/application/config/validation.go
  - internal/application/config/validation_test.go
parent_task_id: TASK-2
ordinal: 18000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Reject incomplete or unsafe configuration before remote operations.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Invalid names, endpoints, and required values produce actionable errors
- [x] #2 Validation never exposes credentials
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add deterministic Application-layer validation for required and unique instance names, safe absolute HTTP(S) endpoints, and required credential references.
2. Return actionable field-specific errors without embedding endpoint or credential values.
3. Add focused tests for valid configuration, each invalid field category, deterministic first-error behavior, and sensitive-value redaction.
4. Run formatting and all Go tests.
5. Check all acceptance criteria and finalize TASK-2.4 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Validation remains storage-independent in internal/application/config. Loading, persistence, profile management, and credential resolution are outside TASK-2.4.

Added deterministic Configuration.Validate behavior in the Application layer.
Validation rejects empty configurations, missing or duplicate instance names, missing or unsafe endpoints, and missing credential references with field-specific errors.
Endpoint and credential reference values are never embedded in validation errors; tests include URL credentials and a sensitive credential reference.
Validation passed: gofmt completed and go test ./... passed.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added minimal storage-independent configuration validation with actionable deterministic errors for names, endpoints, and required values. Added focused tests for valid and invalid configurations and sensitive-value redaction; formatting and all Go tests pass.
<!-- SECTION:FINAL_SUMMARY:END -->
