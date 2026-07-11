---
id: TASK-2.9
title: Standardize common CLI errors
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-1
dependencies:
  - TASK-1.3
  - TASK-1.4
references:
  - ROADMAP.md
modified_files:
  - README.md
  - cmd/fj/main.go
  - internal/interface/cli/errors.go
  - internal/interface/cli/errors_test.go
  - internal/interface/cli/root.go
  - internal/interface/cli/root_test.go
parent_task_id: TASK-2
ordinal: 23000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Present validation, authentication, remote, and internal failures consistently.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Common failure categories have distinct process outcomes
- [x] #2 Errors identify the failed operation without leaking secrets
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add a private Interface-layer CLI error category and safe presentation contract for validation, authentication, remote, and internal failures.
2. Map categories internally to process outcomes internal=1, validation=2, authentication=3, and remote=4; classify unwrapped Cobra input failures as validation.
3. Use developer-controlled fixed operation names, rely on Authentication/Infrastructure boundaries for secret removal, and never display underlying causes for internal failures.
4. Update main execution, focused Interface tests, and README category behavior without publishing numeric exit-code values or adding remote transport or JSON contracts.
5. Run required checks and obtain an independent GPT-5 post-implementation review, then finalize only TASK-2.9.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — public exit behavior, compatibility, architecture, and secret-safe error presentation require high-capability design and verification

Pre-implementation check (GPT-5): Critical none; Major findings require human approval for (1) exit-code mapping, (2) Interface-layer ownership of classification and process mapping, and (3) fixed operation names plus boundary redaction and hidden internal causes. Minor none. Suggestion: use internal=1, validation=2, authentication=3, remote=4; classify unknown CLI input as validation; define remote presentation synthetically without transport; defer JSON contracts. No code or documentation implementation has started pending approval.

Human approved the Major Change design, including internal=1, validation=2, authentication=3, remote=4, Interface-layer classification, fixed operation names, boundary redaction, hidden internal causes, and exclusion of remote transport and JSON. Additional approved constraint: keep numeric exit codes internal and do not publish them as a README contract.

Implemented Interface-owned private categories and safe presentation for internal, validation, authentication, and remote failures. Internal exit mapping is 1/2/3/4 respectively; Cobra unknown flags and unexpected arguments are validation. Fixed developer-controlled operation names and category-safe messages prevent underlying causes from reaching CLI output. Unclassified errors become safe internal failures. README documents category distinction but explicitly keeps numeric values an unpublished implementation detail. Remote transport and JSON contracts were not added.

Validation: gofmt -l . produced no output; git diff --check passed; focused go test ./internal/interface/cli passed; go vet ./... passed; go test ./... passed. Initial focused and full checks encountered the known sandbox restriction on the Go build cache and passed when rerun with approved elevated cache access. During focused testing, unexpected arguments initially mapped to internal; the Cobra adapter was corrected to classify unknown-command input as validation, then all tests passed.

Post-implementation review (independent GPT-5): Critical none; Major none; Minor none; Suggestion none. Confirmed approved scope only, both acceptance criteria satisfied, Interface-layer ownership and dependency direction preserved, causes and secrets hidden, internal exit mapping correct, no remote transport or JSON work, numeric values not published as a README contract, sufficient tests, and accurate Backlog records.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Standardized validation, authentication, remote, and internal CLI failures with distinct internal process outcomes and safe operation-aware messages. Cobra input failures are validation, unknown causes become redacted internal failures, and numeric exit codes remain unpublished implementation details. Updated focused tests and README; all required checks and the independent Major Change review passed.
<!-- SECTION:FINAL_SUMMARY:END -->
