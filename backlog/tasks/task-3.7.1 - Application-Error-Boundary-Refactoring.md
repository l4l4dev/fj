---
id: TASK-3.7.1
title: Application Error Boundary Refactoring
status: Done
assignee: []
created_date: '2026-07-11 00:31'
updated_date: '2026-07-11 00:49'
labels: []
dependencies: []
parent_task_id: TASK-3.7
priority: high
ordinal: 83000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
## Goal
Move remote and validation semantics into the Application boundary before M3 while preserving M1/M2 behavior.

## Scope
- Define semantic application error categories independent of HTTP status codes.
- Translate Infrastructure failures into safe application errors.
- Consolidate CLI error mapping.
- Standardize ValidationError across existing use cases.
- Preserve existing exit codes, operation names, and safe messages.

## Non-goals
- M3 issue workflows.
- JSON output.
- New remote capabilities.
- Composition Root or presenter refactoring covered by sibling tasks.

## Validation
- Run focused error and validation tests.
- Run all existing M1/M2 regression tests and the repository verification workflow.

## Review
- Review error boundary and compatibility before implementation.
- Perform an independent post-implementation review.
- Classify findings as Critical, Major, Minor, or Suggestion.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Application errors expose semantic categories without requiring CLI interpretation of HTTP status codes.
- [x] #2 Infrastructure maps HTTP, transport, and JSON failures to safe application errors without secrets, response bodies, URLs, or raw causes.
- [x] #3 CLI error mapping is consolidated and preserves existing categories, messages, operation names, and exit codes.
- [x] #4 All M1/M2 use cases use one typed ValidationError while preserving validation behavior and messages.
- [x] #5 Focused and full regression tests pass, including make pre-commit.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Introduce the semantic error model, migrate Infrastructure and Use Cases, consolidate CLI presentation, and verify compatibility before M3.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
First child of TASK-3.7. Must complete before Composition Root and Presenter hardening.

Model: Fable 5 — Section 15 pre-implementation design review of the error boundary design. Result: 2 Major (error category package placement for M3 reuse; inconsistent unknown-error fallback category makes literal AC#3 preservation impossible), 3 Minor, 2 Suggestions. Reported to human for decisions before implementation.

Design decisions approved:

Decision 1: Place the shared Application error category model in internal/application/apperror rather than the repository package. This boundary is shared by M3 and later Application workflows and avoids coupling errors to a resource package.

Decision 2: Classify unknown errors as internal in the CLI fallback. This avoids treating unexpected failures as validation and aligns with the safe behavior established for update, archive, and access operations.

Decision 3: Classify HTTP 401/403 uniformly as Authentication to preserve the existing compatibility contract.

Decision 4: Preserve 404 and 409 as Application categories NotFound and Conflict, while mapping both to the existing remote CLI exit code for compatibility.

Validation: M1/M2 regression tests passed and make pre-commit passed, including git diff --check, go vet ./..., and go test ./....

Review: Independent post-implementation review (GPT-5) found Critical none and Major none. Minor: legacy repository.RemoteError remains for compatibility but is not referenced by CLI; future deprecation may be considered.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented the Application error boundary in internal/application/apperror, migrated Infrastructure failures to semantic categories, consolidated CLI error mapping without HTTP status interpretation, standardized ValidationError across M1/M2 use cases, and preserved existing CLI commands, flags, exit codes, and human-readable messages. Regression tests and make pre-commit passed.
<!-- SECTION:FINAL_SUMMARY:END -->
