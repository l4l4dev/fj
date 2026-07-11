---
id: TASK-2.3
title: Define the configuration model
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-1
dependencies: []
references:
  - ROADMAP.md
modified_files:
  - internal/application/config/model.go
  - internal/application/config/model_test.go
parent_task_id: TASK-2
ordinal: 17000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Represent settings required for multiple Forgejo instances.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Named instances, endpoints, and credential references are represented
- [x] #2 Storage concerns do not leak into core configuration concepts
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add a storage-agnostic configuration model under internal/application/config for named Forgejo instances, dedicated endpoint values, and credential references.
2. Add focused tests that construct multiple named instances using the dedicated Endpoint type and verify their represented values.
3. Run formatting and all Go tests.
4. Check all acceptance criteria and finalize TASK-2.3 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
The configuration model belongs to the Application layer because it coordinates runtime settings while remaining independent of storage and infrastructure adapters. Validation, loading, persistence, selection, and authentication are outside TASK-2.3.

Added a storage-agnostic Application-layer configuration model with Configuration, Instance, and CredentialReference types.
The model represents multiple named Forgejo instances, endpoint values, and credential references without storage tags, file paths, environment variables, loaders, or persistence dependencies.
Validation passed: gofmt completed and go test ./... passed, including the new multiple-instance model test.

Review feedback received: introduce type Endpoint string, use it from Instance, and update the existing test. No other changes are approved.

Addressed review feedback by adding type Endpoint string, changing Instance.Endpoint to Endpoint, and updating test fixtures to construct Endpoint values explicitly. Validation passed: gofmt completed, go test ./... passed, and git diff --check passed.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Defined the minimal storage-independent configuration model with dedicated Endpoint and CredentialReference types for named Forgejo instances. Updated the focused representation test and verified formatting and the full Go test suite.
<!-- SECTION:FINAL_SUMMARY:END -->
