---
id: TASK-4.8
title: Issue metadata management
status: To Do
assignee: []
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 06:12'
labels: []
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-4
ordinal: 37000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Update labels, milestones, and assignments independently.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Each metadata category changes without overwriting the others
- [ ] #2 Unknown or inaccessible values produce actionable errors
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
This parent task is split into TASK-4.8.1 Manage issue labels, TASK-4.8.2 Manage issue milestone, and TASK-4.8.3 Manage issue assignment because labels, milestones, and assignments have distinct CLI contracts, Forgejo API contracts, Application Ports, validation, error boundaries, and Presenter output. Each child is a 30-90 minute task and follows Decision approval → Implementation → Verification → Independent Review → Acceptance Criteria completion → Done.
<!-- SECTION:NOTES:END -->
