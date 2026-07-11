---
id: TASK-1.4
title: Define public compatibility policy
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
  - CONTRIBUTING.md
parent_task_id: TASK-1
ordinal: 13000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Clarify which interfaces are compatibility-sensitive.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Commands, configuration, JSON output, and exit behavior are addressed
- [x] #2 Incompatible changes require review and migration guidance
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Review the existing compatibility commitments and complete the required Major Change pre-implementation assessment.
2. Expand only the Compatibility section of CONTRIBUTING.md to identify compatibility-sensitive interfaces and the review and migration requirements for incompatible changes.
3. Run formatting, static analysis, tests, documentation checks, and inspect the scoped diff.
4. Obtain the required independent post-implementation review, then record results and finalize TASK-1.4.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — public compatibility policy is a Major Change requiring high-capability design and verification

Pre-implementation check: GPT-5 subagent found no Critical, Major, or Minor findings. It confirmed the minimal CONTRIBUTING.md-only approach is consistent with governing documents, changes no runtime behavior or security boundary, has appropriate task granularity, and needs no additional human approval. Suggestion: defer SemVer, fixed deprecation periods, and stability tiers to separate future work if needed.

Implemented the approved minimal policy in CONTRIBUTING.md: identified compatibility-sensitive commands and flags, configuration, documented output and JSON contracts, exit behavior, and exported APIs; required explicit human review and approval plus old/new behavior, user impact, actionable migration guidance, documentation updates, and compatibility tests where implementation exists. No versioning, deprecation-period, or stability-tier guarantee was introduced. Validation passed: gofmt -l . (no output); git diff --check; go vet ./...; go test ./...; governing-document link target existence check.

Post-implementation review: independent GPT-5 subagent found no Critical, Major, Minor, or Suggestion findings. It confirmed approved scope only, both acceptance criteria satisfied, architecture boundaries preserved, no compatibility or security regression, sufficient checks, consistent documentation, and accurate Backlog records.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Defined the public compatibility policy in CONTRIBUTING.md for commands, configuration, JSON and other documented output, exit behavior, and exported APIs. Documented the review, approval, impact, migration, documentation, and testing requirements for incompatible changes. All checks and the required independent Major Change review passed.
<!-- SECTION:FINAL_SUMMARY:END -->
