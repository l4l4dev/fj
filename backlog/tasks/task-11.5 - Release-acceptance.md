---
id: TASK-11.5
title: Release acceptance
status: To Do
assignee: []
created_date: '2026-07-11 17:33'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-10
dependencies:
  - TASK-11.4
parent_task_id: TASK-11
priority: medium
ordinal: 11050
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define and execute acceptance for release artifacts, including install, version, help, read-only smoke checks, and secret redaction. TASK-10.5 results are inputs but not a required dependency.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Initial artifacts can be installed and verified.
- [ ] #2 fj version and fj --help are confirmed after installation.
- [ ] #3 Read-only smoke checks and secret redaction are evaluated.
- [ ] #4 TASK-10.5 unresolved commands are not treated as successful.
<!-- AC:END -->
