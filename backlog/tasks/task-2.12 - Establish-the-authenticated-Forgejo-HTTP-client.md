---
id: TASK-2.12
title: Establish the authenticated Forgejo HTTP client
status: To Do
assignee: []
created_date: '2026-07-10 17:08'
labels: []
dependencies:
  - TASK-2.10
  - TASK-2.11
references:
  - ROADMAP.md
parent_task_id: TASK-2
ordinal: 78000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the infrastructure boundary for authenticated Forgejo HTTP requests using net/http only. Compose the selected instance endpoint and environment-provided credential without introducing a Forgejo SDK.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Requests use net/http and attach the resolved credential without exposing it in diagnostics.
- [ ] #2 The client targets the selected instance endpoint and preserves remote failures for the Interface error presentation boundary.
- [ ] #3 Every request sends User-Agent in the form fj/<version>.
<!-- AC:END -->
