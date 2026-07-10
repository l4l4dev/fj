---
id: TASK-2.14
title: Implement the repository REST adapter
status: To Do
assignee: []
created_date: '2026-07-10 17:08'
labels: []
dependencies:
  - TASK-2.12
  - TASK-2.13
references:
  - ROADMAP.md
parent_task_id: TASK-2
ordinal: 80000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the Infrastructure-layer REST adapter for repository listing using the authenticated Forgejo HTTP client and the Repository Service Port. Use only net/http and the approved Forgejo endpoint.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The adapter retrieves repositories from GET /api/v1/user/repos.
- [ ] #2 Pagination sends only page and limit parameters; fetch-all behavior is not implemented.
- [ ] #3 Remote errors and empty responses are translated into outcomes the Application and Interface layers can present clearly.
<!-- AC:END -->
