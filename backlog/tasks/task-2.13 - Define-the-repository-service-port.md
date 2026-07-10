---
id: TASK-2.13
title: Define the repository service port
status: To Do
assignee: []
created_date: '2026-07-10 17:08'
labels: []
dependencies: []
references:
  - ROADMAP.md
parent_task_id: TASK-2
ordinal: 79000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define the Application-layer Repository Service Port required by repository use cases. Keep the port transport-independent and limit this Task to the service contract.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The port expresses repository listing capability without depending on net/http, TOML, Cobra, or Forgejo SDK types.
- [ ] #2 Repository models and remote failure outcomes crossing the boundary are explicit and testable.
- [ ] #3 Repository list use-case orchestration is explicitly outside this Task and belongs to TASK-3.1.
<!-- AC:END -->
