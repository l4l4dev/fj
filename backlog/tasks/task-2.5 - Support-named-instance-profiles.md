---
id: TASK-2.5
title: Support named instance profiles
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
  - internal/application/config/profiles.go
  - internal/application/config/profiles_test.go
parent_task_id: TASK-2
ordinal: 19000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Allow users to identify separately configured Forgejo instances.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Profiles can be listed and inspected safely
- [x] #2 Profile output never reveals secrets
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add a safe Profile view containing only an instance name and endpoint.
2. Add Application-layer operations to list profiles and inspect one profile by name after configuration validation.
3. Add focused tests for listing, inspection, missing profiles, and absence of credential values from results and errors.
4. Run formatting and all Go tests.
5. Check all acceptance criteria and finalize TASK-2.5 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Profile operations remain in the Application-layer configuration boundary. The safe Profile result excludes CredentialReference; storage, CLI commands, and instance-selection precedence are outside TASK-2.5.

Added a safe Profile view containing only Name and Endpoint.
Added Configuration.ListProfiles and Configuration.InspectProfile; both validate configuration before returning profile data.
CredentialReference is excluded from profile results, and tests verify that list, inspection, and error output do not reveal secret values.
Validation passed: gofmt completed and go test ./... passed.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added minimal Application-layer operations to list and inspect validated named instance profiles through a credential-free Profile view. Added focused safety and behavior tests; formatting and all Go tests pass.
<!-- SECTION:FINAL_SUMMARY:END -->
