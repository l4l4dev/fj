---
id: TASK-10.3
title: Forgejo Playground configuration onboarding
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 09:09'
updated_date: '2026-07-11 09:27'
labels: []
dependencies:
  - TASK-10.2
parent_task_id: TASK-10
priority: high
ordinal: 10020
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Document config.toml, credential environment variables, instance selection, and --instance usage for Forgejo Playground.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The XDG config path and TOML schema are documented.
- [x] #2 Credential environment variable setup is documented without exposing token values.
- [x] #3 Forgejo Playground instance configuration is documented.
- [x] #4 Explicit --instance usage and single-instance selection behavior are documented.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Verification:
- git diff: passed
- git diff --check: passed
- gofmt -l .: passed
- go vet ./...: passed
- go test ./...: passed
- make pre-commit: passed

Independent Review:
- Critical: none
- Major: none
- Minor: none
- Suggestion: Quickstart全体はTASK-10.4で扱う。
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Documented XDG configuration paths, config.toml instance schema, Forgejo Playground placeholder onboarding, credential environment-variable setup, --instance selection, and secret redaction guidance. Verified repository checks and preserved the approved TASK-10.4 Quickstart boundary.
<!-- SECTION:FINAL_SUMMARY:END -->
