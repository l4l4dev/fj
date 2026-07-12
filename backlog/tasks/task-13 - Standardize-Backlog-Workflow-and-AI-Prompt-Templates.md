---
id: TASK-13
title: Standardize Backlog Workflow and AI Prompt Templates
status: Done
assignee:
  - '@codex'
created_date: '2026-07-12 07:38'
updated_date: '2026-07-12 07:46'
labels: []
dependencies: []
modified_files:
  - AGENTS.md
  - backlog/templates/README.md
  - backlog/templates/workflow.md
  - backlog/templates/00-investigation.md
  - backlog/templates/01-decision.md
  - backlog/templates/02-implementation.md
  - backlog/templates/03-independent-review.md
  - backlog/templates/04-remediation.md
  - backlog/templates/05-independent-rereview.md
  - backlog/templates/06-finalization.md
  - backlog/templates/prompts/investigation.md
  - backlog/templates/prompts/decision-recording.md
  - backlog/templates/prompts/implementation.md
  - backlog/templates/prompts/independent-review.md
  - backlog/templates/prompts/remediation.md
  - backlog/templates/prompts/independent-rereview.md
  - backlog/templates/prompts/finalization.md
ordinal: 92000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Standardize the workflow established through TASK-11 and TASK-12 as reusable repository templates. The system must keep the repository as the Source of Truth, make Human Decision gates explicit, separate Implementation, Independent Review, and Finalization roles, separate project record templates from AI prompt templates, allow proportionate simplification for simple tasks, and keep commit and push as separately authorized actions.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The standard workflow and responsibilities of each stage are documented
- [x] #2 Human Decision gates, scope control, independent review, finding handling, and finalization conditions are explicit
- [x] #3 Reusable project record templates are provided under backlog/templates
- [x] #4 Reusable AI prompt templates are provided separately under backlog/templates/prompts
- [x] #5 Prompt placeholder naming is consistent across templates
- [x] #6 Conditions for simplifying the workflow for simple tasks are documented
- [x] #7 AGENTS.md links to the official template entry point with a minimal reference
- [x] #8 No unnecessary changes are made to existing tasks, code, workflows, root README, or LICENSE
- [x] #9 Commit and push remain prohibited until separately authorized
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add the reusable project workflow and stage record templates under backlog/templates.
2. Add copy-ready English AI prompt templates under backlog/templates/prompts using a consistent placeholder vocabulary.
3. Add one minimal AGENTS.md reference to the official template entry point.
4. Verify internal links, stage/template consistency, decision and review separation, finalization and commit boundaries, placeholder consistency, scope, and Markdown quality.
5. Run make pre-commit and obtain an independent post-implementation review without finalizing the task.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Human Decision — Backlog Workflow Template System:
- Introduce a reusable Backlog workflow template system for fj.
- Standard workflow: Investigation → Human Decision when required → Implementation → Independent Review → Remediation when required → Independent Re-Review when required → Finalization → Commit only after explicit authorization.
- Separate project workflow and record templates from copy-ready AI agent prompt templates.
- Primary location: `backlog/templates/`; AI prompts: `backlog/templates/prompts/`.
- Permit one minimal AGENTS.md reference to the official template entry point.
- Allowed changes: this TASK-13 record, `backlog/templates/**`, and the minimal AGENTS.md reference.
- Prohibited changes: existing tasks, code, workflows, root README, LICENSE, Go modules, external dependencies, status Done, Acceptance Criteria checks, Final Summary, commit, push, tag, and release creation.
- Implementation begins only after this Decision is recorded.
- This is a project-governance Decision approved by a human.

Model: GPT-5 (Codex) — project-wide workflow standardization requires careful governance alignment, template design, consistency verification, and independent review.

Pre-Implementation Check:
- PROJECT_CONSTITUTION.md, ARCHITECTURE.md, ROADMAP.md, compatibility, security, and task granularity: Pass.
- `backlog/templates/` does not conflict with Backlog.md-managed task, document, draft, decision, milestone, or archive directories.
- Templates use ordinary Markdown without Backlog task frontmatter.
- Existing `.agent/` files are operational helpers; the new reusable records and copy-ready prompts have a distinct documented role.
- No additional Human Decision is required beyond the Decision recorded above.

Implementation and Verification:
- Added reusable record templates for Investigation, Decision, Implementation, Independent Review, Remediation, Independent Re-Review, and Finalization, plus the workflow overview and entry-point README.
- Added seven copy-ready English AI prompt templates using the approved placeholder vocabulary.
- Added one minimal AGENTS.md link to the official template entry point.
- Template structure, README internal links, placeholder vocabulary, stage boundaries, and explicit commit/push boundaries: Pass.
- `git diff --check`: Pass.
- `go vet ./...`: Pass.
- `go test ./...`: Pass.
- `make pre-commit`: Pass.
- No code, workflow, root README, LICENSE, existing task, dependency, external state, commit, or push change was made.

Evidence Summary:
- A reusable task workflow was documented from Investigation through explicit Human Decision, Implementation, Independent Review, optional Remediation and Independent Re-Review, Finalization, and separately approved Commit.
- Project record templates and copy-ready AI prompt templates were separated under `backlog/templates/` and `backlog/templates/prompts/`.
- Sixteen template files were added with consistent stage boundaries, responsibilities, outputs, prohibited actions, exit conditions, and next steps.
- The prompt templates use a consistent placeholder vocabulary and are written in English for direct use with AI agents.
- Repository-as-Source-of-Truth, Human Decision Gate, Scope Control, reviewer independence, Deferred Item boundaries, and explicit commit/push authorization are represented throughout the templates.
- Simple tasks may use a reduced workflow only when their scope is clear and they do not involve public behavior, compatibility, security, legal, governance, architecture, dependency, or migration decisions.
- `AGENTS.md` now links to `backlog/templates/README.md` as the official entry point without redefining existing governance rules.
- Existing tasks, source code, workflows, the root README, LICENSE, and Go module files were not changed.
- Verification completed successfully: expected template-file presence, Markdown link checks, placeholder consistency, stage-boundary checks, commit/push boundary checks, Finalization-only Done transition checks, `git diff --check`, `go vet ./...`, `go test ./...`, and `make pre-commit`.

Independent Review (post-implementation):
- Scope: TASK-13, `AGENTS.md`, `backlog/templates/**`, workflow-stage definitions, project record templates, AI prompt templates, placeholder consistency, simplified-task rules, and current Git scope.
- Acceptance Criteria #1-#9: Pass.
- Critical: none.
- Major: none.
- Minor: none.
- Suggestions: none.
- The review confirmed that the template system is compatible with the existing Backlog.md directory structure and does not require task frontmatter.
- The review confirmed separation between project records and AI prompts, Repository-as-Source-of-Truth rules, Human Decision Gate enforcement, implementation/reviewer role separation, scoped remediation, and Finalization-only status completion.
- The review confirmed that commit and push remain separate explicit authorization steps.
- Existing `.agent/` operational helpers remain distinct from the persistent workflow and prompt template system.
- Result: Ready for Finalization.

Deferred Items:
- Template usage should be refined only after repeated real-task usage reveals concrete gaps or unnecessary ceremony.
- Additional task-type-specific templates may be added later when supported by recurring project needs.
- Automation for validating placeholders, internal links, or template completeness may be considered separately if manual maintenance becomes burdensome.
- Existing `.agent/` helpers are not migrated into this system unless a separate review identifies clear duplication.
- These items are outside TASK-13 completion scope.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Standardized the fj project’s Backlog workflow and AI prompt process as a reusable repository-owned template system. TASK-13 added persistent record templates for Investigation, Human Decision, Implementation, Independent Review, Remediation, Independent Re-Review, and Finalization, together with copy-ready English AI prompts for each stage.

The system makes the repository the Source of Truth, requires Human Decisions to be recorded before implementation where applicable, separates implementation from independent review, prevents remediation scope expansion, and reserves Acceptance Criteria completion and Done transitions for Finalization after a successful review. Commit and push remain separate explicitly approved actions.

The templates also document when simple tasks may use a reduced workflow without weakening scope control, verification, Git-state reporting, or authorization boundaries. `AGENTS.md` now links to the template system as the official entry point without changing existing governance rules.

Verification and Independent Review passed with no Critical, Major, Minor, or Suggestion findings. TASK-13 is Ready for Finalization.
<!-- SECTION:FINAL_SUMMARY:END -->
