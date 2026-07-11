---
id: TASK-10.4
title: Public CLI Quickstart
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 09:09'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-9
dependencies:
  - TASK-10.1
  - TASK-10.2
  - TASK-10.3
parent_task_id: TASK-10
priority: medium
ordinal: 10040
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Document the validated macOS install, Forgejo Playground configuration, credential setup, help, and read-only command workflow.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Quickstart documents install and PATH verification.
- [x] #2 Quickstart documents config.toml and credential setup.
- [x] #3 Quickstart includes fj --help and read-only repository, issue, and pull request examples.
- [x] #4 Quickstart matches the validated User Acceptance workflow.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Independent Review:
- Critical: none
- Major: none
- Minor: repo list and pr list remain environment-dependent and are explicitly separated from guaranteed examples.
- Suggestion: revisit the environment-dependent section after TASK-10.5 re-acceptance.

Verification:
- git diff --check passed
- make pre-commit passed
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added a privacy-safe Quickstart covering install/config/version and validated read-only commands. Environment-dependent repository and pull request list commands are explicitly separated and not presented as guaranteed successes.
<!-- SECTION:FINAL_SUMMARY:END -->
