---
id: TASK-2.10
title: Load configuration from an XDG TOML file
status: To Do
assignee: []
created_date: '2026-07-10 17:08'
labels: []
dependencies:
  - TASK-2.3
  - TASK-2.4
references:
  - ROADMAP.md
parent_task_id: TASK-2
ordinal: 76000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Load named Forgejo instance profiles from an XDG Base Directory TOML configuration file using BurntSushi/toml. Keep loading separate from the existing configuration model and validation responsibilities.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The configuration file is resolved according to the XDG Base Directory convention.
- [ ] #2 TOML profiles are decoded with BurntSushi/toml and passed through existing configuration validation.
- [ ] #3 Missing, malformed, and invalid configuration files produce actionable errors without exposing credentials.
<!-- AC:END -->
