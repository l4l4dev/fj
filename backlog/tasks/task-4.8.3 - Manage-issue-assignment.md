---
id: TASK-4.8.3
title: Manage issue assignment
status: To Do
assignee: []
created_date: '2026-07-11 06:11'
labels: []
dependencies:
  - TASK-2.9
parent_task_id: TASK-4.8
ordinal: 88000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Manage issue assignment independently as a focused Issue metadata workflow.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The CLI contract for assigning and unassigning an issue is explicitly defined and supports --instance PROFILE.
- [ ] #2 OWNER/NAME and positive issue number validation are enforced before remote access.
- [ ] #3 Application defines an assignment model and dedicated Port without changing existing issue or repository Ports.
- [ ] #4 Infrastructure uses the approved Forgejo assignment API contract and translates failures through the Application error boundary.
- [ ] #5 Assignment changes do not overwrite labels, milestone, title, body, or state.
- [ ] #6 Human-readable assignment output is defined and does not change existing issue command output.
- [ ] #7 Tests cover validation, API mapping, error boundary, Presenter output, and explicit Composition Root injection.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision approval is required before implementation. Workflow: Decision approval → Implementation → Verification → Independent Review → Acceptance Criteria completion → Done.
<!-- SECTION:NOTES:END -->
