---
id: TASK-2.15
title: Add shared Makefile development commands
status: Done
assignee: []
created_date: '2026-07-10 17:48'
updated_date: '2026-07-10 17:49'
labels: []
dependencies: []
modified_files:
  - Makefile
parent_task_id: TASK-2
priority: medium
ordinal: 81000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide common local, Codex, Claude, and future CI verification commands through a repository Makefile after M1 completion.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Makefile provides fmt, check-fmt, vet, test, and build targets.
- [x] #2 verify runs check-fmt, git diff --check, vet, and test.
- [x] #3 pre-commit runs verify.
- [x] #4 Commands use existing Go tooling without adding dependencies.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Add Makefile targets for formatting, formatting checks, vetting, tests, builds, verification, and pre-commit checks.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — small but cross-tool repository workflow task; selected to ensure command sequencing matches project verification policy.

Validation: make fmt, make check-fmt, make vet, make test, make build, make verify, make pre-commit, and git diff --check passed. go build ./... emitted a sandbox cache warning but exited successfully. Model: GPT-5 — selected for consistency review across repository development workflow commands; no sub-agent used.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added Makefile targets for shared formatting, static analysis, tests, build, verification, and pre-commit workflows. verify runs check-fmt, git diff --check, vet, and test in the required order. No README change was necessary because existing contribution guidance already defines the underlying commands.
<!-- SECTION:FINAL_SUMMARY:END -->
