---
id: TASK-1.1
title: Verify project document hierarchy
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 05:50'
labels: []
dependencies: []
references:
  - ROADMAP.md
modified_files:
  - AGENTS.md
  - CLAUDE.md
parent_task_id: TASK-1
ordinal: 10000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Ensure governing documents define a consistent source-of-truth hierarchy.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The governing documents reference each other consistently
- [x] #2 Conflicting authority or workflow statements are identified for review
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Update AGENTS.md with the approved document precedence and an explicit prohibition against task definitions overriding PROJECT_CONSTITUTION.md.
2. Add ARCHITECTURE.md, ROADMAP.md, CLAUDE.md, and backlog/ to the AGENTS.md Repository Structure section.
3. Add a short CLAUDE.md introduction defining it as a supplement to AGENTS.md that cannot override PROJECT_CONSTITUTION.md or AGENTS.md.
4. Verify only the approved files changed, then finalize TASK-1.1 through the Backlog.md CLI.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Hierarchy review findings:

1. Ambiguous precedence: PROJECT_CONSTITUTION.md defines itself as the project's highest-level policy, while AGENTS.md places approved tasks and explicit human instructions above the constitution. The intended boundary between task-level authority and constitutional policy is not explicit.
2. Missing relationship: CLAUDE.md contains only Backlog.md workflow instructions and does not reference PROJECT_CONSTITUTION.md or delegate repository-wide behavior to AGENTS.md, although AGENTS.md states that it applies to Claude Code.
3. Incomplete document inventory: the Repository Structure section in AGENTS.md omits ARCHITECTURE.md, ROADMAP.md, CLAUDE.md, and the Backlog.md-managed area.
4. Ambiguous peer ordering: AGENTS.md groups all other approved repository documentation at one precedence level, so the relationship among ARCHITECTURE.md, ROADMAP.md, and CLAUDE.md is not explicit if they ever conflict.
5. Consistent relationships confirmed: ARCHITECTURE.md is explicitly subordinate to PROJECT_CONSTITUTION.md and follows AGENTS.md; ROADMAP.md explicitly references PROJECT_CONSTITUTION.md, AGENTS.md, and ARCHITECTURE.md; the Backlog.md workflow blocks in AGENTS.md and CLAUDE.md are identical.

No governing document was modified because TASK-1.1 authorizes verification only. Acceptance criterion 1 remains unmet pending an approved documentation change.

Human approval received to resolve the documented hierarchy and relationship gaps through the three specified documentation changes.

Approved documentation changes completed. Validation passed: the precedence order is explicit, task definitions cannot override PROJECT_CONSTITUTION.md, the Repository Structure includes all requested entries, CLAUDE.md delegates to AGENTS.md, referenced paths exist, and git diff --check reports no errors.

Historical note: This task was completed before the standard workflow was introduced. No Verification execution record or Independent Review record exists from that period.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Resolved the approved document hierarchy gaps in AGENTS.md and CLAUDE.md. Verified the requested precedence, repository inventory, supplement relationship, referenced paths, and Markdown whitespace; no code or unrelated files were changed.
<!-- SECTION:FINAL_SUMMARY:END -->
