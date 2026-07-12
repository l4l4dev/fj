---
id: TASK-11.5
title: Release acceptance
status: To Do
assignee: []
created_date: '2026-07-11 17:33'
updated_date: '2026-07-12 00:36'
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
- [ ] #5 A live workflow_dispatch run verifies resolve-version, darwin/arm64 and linux/amd64 builds, SHA-256 checksum generation, and the consolidated workflow artifact.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Design Review (pre-acceptance): see doc-1 - TASK-11.5 Release Acceptance Design Review. Result: Critical 0 / Major 3 (LICENSE missing; live workflow run unverified; workflow_dispatch tag-commit binding) / Minor 6 / Suggestion 4. No fixes applied; findings await human decisions.
Model: Fable 5 — long-term design review of the Public Release Foundation (Section 15 review-grade work).

Human Decisions — Design Review Findings:
- LICENSE (Major): accepted as separate work. Tracked by TASK-12. No license has been selected; exact license selection requires a separate explicit human decision.
- GitHub Actions live execution (Major): accepted in TASK-11.5. A successful workflow_dispatch run must verify resolve-version, both build targets, checksum generation, and consolidated artifact generation.
- workflow_dispatch arbitrary version (Major): deferred for the release foundation. TASK-11.5 does not change the current input policy. A future public Release workflow must define a separate tag-only policy.
- Minor and Suggestion findings from doc-1 remain future improvement candidates and do not expand TASK-11.5 scope.
<!-- SECTION:NOTES:END -->
