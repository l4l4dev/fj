---
id: TASK-2.2
title: Define global command behavior
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
  - README.md
  - internal/interface/cli/root.go
  - internal/interface/cli/root_test.go
parent_task_id: TASK-2
ordinal: 16000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Establish consistent shared flags, output, and failure behavior.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Global behavior is documented and tested
- [x] #2 Invalid global input produces a clear deterministic result
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Document the current global CLI contract in README.md: help flag, stdout, stderr, exit behavior, and invalid input.
2. Configure the root command to reject unexpected positional arguments and suppress usage output after errors.
3. Add focused interface tests for unknown flags and unexpected arguments while preserving the root-help test.
4. Run formatting, all relevant tests, direct CLI failure checks, and diff validation.
5. Check all acceptance criteria and finalize TASK-2.2 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Scope is limited to the global behavior defined by TASK-2.2. No feature commands, configuration, or work from later tasks will be added.

Defined and documented the minimal global command contract: -h/--help on stdout with success; failures on stderr with non-zero status and no repeated usage.
Configured the root command to reject unexpected positional arguments and suppress usage after errors.
Added deterministic tests for an unknown flag and an unexpected argument.
Validation passed: gofmt completed, go test ./... passed, both invalid-input CLI checks returned exit code 1 with the documented error, and git diff --check passed.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Defined the minimal global CLI behavior in README.md and enforced deterministic invalid-input handling in the Cobra root command. Added focused tests and verified formatting, the full Go test suite, actual failure output and exit status, and diff hygiene.
<!-- SECTION:FINAL_SUMMARY:END -->
