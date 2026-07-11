---
id: TASK-10.2
title: Local install on macOS
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 09:09'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-9
dependencies:
  - TASK-2.15
parent_task_id: TASK-10
priority: high
ordinal: 10050
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Install the fj binary for the current user at $HOME/.local/bin/fj and provide a matching uninstall flow.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 make install places the fj binary at $HOME/.local/bin/fj.
- [x] #2 make uninstall removes only $HOME/.local/bin/fj.
- [x] #3 Install and uninstall do not modify configuration, credentials, PATH, or shell settings.
- [x] #4 README documents install, PATH checking, PATH setup guidance, and uninstall.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed before TASK-10 was created.

Verification:
- gofmt -l . passed
- git diff --check passed
- go vet ./... passed
- go test ./... passed
- make pre-commit passed

Independent Review:
- Critical: none
- Major: none
- Minor: Task description should prefer $HOME/.local/bin/fj for portability; install and uninstall runtime checks are recorded.
- Suggestion: retain fj --help as user acceptance evidence.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented and accepted macOS user-local installation at $HOME/.local/bin/fj. make install, fj --help, and make uninstall workflow are documented and accepted; configuration, credentials, PATH, and shell settings remain untouched.
<!-- SECTION:FINAL_SUMMARY:END -->
