---
id: TASK-1.3
title: Define development verification workflow
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
ordinal: 12000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Establish the checks expected before work is complete.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Formatting, static analysis, testing, and documentation checks are identified
- [x] #2 Handling of skipped or failed checks is documented
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Review the repository toolchain, existing contribution standards, and governing verification rules.
2. Expand CONTRIBUTING.md with the required formatting, static analysis, testing, and documentation checks, including how failures and skipped checks are handled.
3. Run the documented checks and verify the documentation diff.
4. Record verification results, complete the acceptance criteria, and finalize TASK-1.3.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — repository-wide verification policy requires careful consistency with governing documents

Expanded the contributor verification workflow with explicit formatting, static analysis, full-suite testing, and documentation checks. Documented that failed checks block completion and that skipped, inapplicable, or environment-blocked checks require the reason, unverified scope, and alternative verification to be recorded. Validation passed: gofmt -l . (no output); git diff --check; go vet ./...; go test ./...; governing-document link target existence check.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Defined the required development verification workflow in CONTRIBUTING.md, including concrete formatting, static analysis, testing, and documentation checks plus explicit handling for failed, skipped, inapplicable, and environment-blocked checks. All documented repository checks passed.
<!-- SECTION:FINAL_SUMMARY:END -->
