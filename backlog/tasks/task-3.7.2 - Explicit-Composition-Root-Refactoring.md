---
id: TASK-3.7.2
title: Explicit Composition Root Refactoring
status: Done
assignee: []
created_date: '2026-07-11 00:31'
updated_date: '2026-07-11 01:23'
labels: []
dependencies:
  - TASK-3.7.1
parent_task_id: TASK-3.7
priority: high
ordinal: 84000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
## Goal
Make repository capability dependencies explicit before M3 and remove runtime type assertion wiring.

## Scope
- Define an explicit dependency structure for repository capabilities.
- Wire List, Getter, Creator, Updater, Archiver, and AccessViewer explicitly.
- Preserve existing command constructors, CLI contracts, and behavior through compatible migration wrappers where needed.
- Keep composition responsibilities at the application boundary.

## Non-goals
- M3 issue dependencies.
- New repository capabilities.
- JSON output.
- Presenter redesign beyond integration compatibility.

## Validation
- Test missing and complete dependency wiring.
- Run all M1/M2 command regression tests and the repository verification workflow.

## Review
- Review dependency direction and compatibility before implementation.
- Perform an independent post-implementation review.
- Classify findings as Critical, Major, Minor, or Suggestion.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Composition Root passes explicit capability dependencies without runtime type assertions.
- [x] #2 Existing repository commands and their public CLI contracts remain unchanged.
- [x] #3 All existing M1/M2 capabilities remain wired and tested.
- [x] #4 Focused composition tests and make pre-commit pass.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Introduce explicit repository dependencies, migrate command construction, retain compatibility wrappers if required, and verify all M1/M2 behavior.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Second child of TASK-3.7. Depends on Application Error Boundary Refactoring.

Design decisions approved:

Decision 1: Place RepositoryDependencies as a Composition-only type in internal/interface/cli. This keeps Composition Root responsibility in the Interface layer and avoids bringing CLI-dependent structures into Application.

Decision 2: Retain NewRootCommandWithRepositoryService as a migration wrapper. This limits impact on existing tests and internal callers while moving to explicit dependency injection.

Decision 3: Each command constructor receives only the Ports it requires. This removes runtime type assertions and makes required capabilities explicit at compile time.

Validation: gofmt, git diff --check, go vet ./..., go test ./..., and make pre-commit all passed.

Independent Review: Critical: none. Major: none. Minor: runtime assertions remain inside the legacy wrapper for staged migration compatibility. Suggestion: consider deprecating and removing the wrapper in a future task.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented explicit Composition Root dependency injection with RepositoryDependencies in internal/interface/cli. Each repository command receives only the required Application Port, command execution no longer performs runtime capability assertions, and the legacy NewRootCommandWithRepositoryService wrapper remains for staged migration compatibility. Existing CLI commands, flags, arguments, exit codes, and human-readable output remain unchanged.
<!-- SECTION:FINAL_SUMMARY:END -->
