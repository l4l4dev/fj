---
id: TASK-11.1
title: Release workflow foundation
status: To Do
assignee: []
created_date: '2026-07-11 17:33'
labels: []
dependencies:
  - TASK-10.1
parent_task_id: TASK-11
priority: high
ordinal: 11010
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define the foundation workflow for reproducible fj builds: read a vMAJOR.MINOR.PATCH tag, strip the v prefix, inject the normalized version through ldflags, and prepare artifacts without creating or publishing public releases.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The workflow maps a vMAJOR.MINOR.PATCH tag to a normalized version.
- [ ] #2 The normalized version is passed through the approved ldflags injection path.
- [ ] #3 The workflow remains a release foundation and does not create or publish a public release.
<!-- AC:END -->
