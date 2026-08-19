---
id: TASK-16
title: Simplify AI development workflow
status: Done
assignee:
  - '@claude'
created_date: '2026-08-19 07:51'
updated_date: '2026-08-19 07:55'
labels: []
dependencies: []
type: chore
ordinal: 95000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Adopt the human-approved Plan A workflow simplification: routine tasks (including new public CLI commands that follow existing patterns) run as plan-approval -> implementation+tests -> report -> commit approval. Reserve the heavy major-change process (pre-implementation check, independent review) for genuine major changes: constitution/architecture/roadmap/AGENTS.md edits, new external dependencies, package structure or dependency-direction changes, auth/credential/security boundaries, breaking changes to existing public contracts, multi-milestone impact, large refactoring. Remove per-task model recording. Slim Backlog records to AC checks plus 1-2 note lines. Quality invariants unchanged: tests required, commits only on explicit approval, PROJECT_CONSTITUTION.md and architecture boundaries respected, privacy placeholders kept.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 AGENTS.md major-change definition excludes pattern-following CLI additions and keeps the heavy process only for genuine major changes
- [x] #2 Per-task model recording and per-task independent review requirements are removed for routine tasks
- [x] #3 DEVELOPMENT_WORKFLOW.md, CLAUDE.md, and backlog/templates reflect the simplified single-loop flow
- [x] #4 Quality invariants (tests, commit approval, constitution, privacy placeholders) remain explicit
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Revise AGENTS.md sections 4, 14-18. 2. Rewrite DEVELOPMENT_WORKFLOW.md lifecycle. 3. Update CLAUDE.md model/review guidance. 4. Replace backlog/templates 7-stage pipeline with the simple loop. 5. Verify doc consistency.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision: Mika approved Plan A (large simplification) in chat on 2026-08-19 after reviewing options A/B/C. This task is itself a major change; the human approval above is the design approval.

Revised AGENTS.md (sections 4, 12, 14-18), rewrote DEVELOPMENT_WORKFLOW.md as a single loop, updated CLAUDE.md, replaced the 7-stage templates and prompts/ with one README reference. Grep confirmed no stale references to the removed lifecycle or templates.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Simplified the AI workflow per human-approved Plan A: heavy checks now apply only to genuine major changes; routine tasks use plan approval -> implement+tests -> report. Verified with make pre-commit and doc-reference grep.
<!-- SECTION:FINAL_SUMMARY:END -->
