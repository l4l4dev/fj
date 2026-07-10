---
id: TASK-3.1
title: List repositories
status: To Do
assignee: []
created_date: '2026-07-10 11:55'
updated_date: '2026-07-10 17:09'
labels: []
dependencies:
  - TASK-2.9
  - TASK-2.10
  - TASK-2.11
  - TASK-2.12
  - TASK-2.13
  - TASK-2.14
references:
  - ROADMAP.md
parent_task_id: TASK-3
ordinal: 24000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Allow users to discover repositories on the selected Forgejo instance.

Scope: Implement only the repository-listing use case, its CLI command, the Composition Root wiring, and human-readable output.

Out of scope: XDG/TOML configuration loading, environment credential resolution, the authenticated HTTP client, the Repository Service Port, the REST adapter, JSON output, and fetch-all pagination behavior.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The CLI Composition Root wires the selected instance and dependencies, renders an explicit human-readable empty result, and presents remote failures through the common CLI error behavior.
- [ ] #2 The use case requests a page and limit through the Repository Service Port and presents each returned repository with an unambiguous owner/name identity.
<!-- AC:END -->







## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Assessment only; implementation intentionally not started. Model: GPT-5 — determining whether this repository-listing Task can fit the approved architecture requires cross-layer boundary analysis.

Blocker: internal/infrastructure and internal/domain directories are empty; no Forgejo REST adapter, repository service port, repository application use case, CLI command, or composition path exists. Implementing TASK-3.1 now would require introducing multiple architectural boundaries and expanding beyond the Task's stated scope. Per the human stop condition, keep TASK-3.1 To Do and request Task decomposition or a separately approved foundation Task.

Baseline validation without code changes: gofmt -l . (no output); git diff --check passed; go vet ./... passed; go test ./... passed. No files were changed by this assessment.
<!-- SECTION:NOTES:END -->
