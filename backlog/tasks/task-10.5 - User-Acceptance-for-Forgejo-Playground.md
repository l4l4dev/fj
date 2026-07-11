---
id: TASK-10.5
title: User Acceptance for Forgejo Playground
status: To Do
assignee: []
created_date: '2026-07-11 09:09'
updated_date: '2026-07-11 09:16'
labels: []
dependencies:
  - TASK-10.3
parent_task_id: TASK-10
priority: high
ordinal: 10030
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Use the installed fj command against Forgejo Playground and validate read-only repository, issue, and pull request workflows, including safe failures.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 fj --help and fj version are verified on macOS.
- [ ] #2 Repository list/inspect read-only commands succeed against Forgejo Playground.
- [ ] #3 Issue list/inspect read-only commands succeed against Forgejo Playground.
- [ ] #4 Pull request list/inspect read-only commands succeed against Forgejo Playground.
- [ ] #5 Invalid input, authentication/remote failures, and secret redaction are checked.
<!-- AC:END -->
