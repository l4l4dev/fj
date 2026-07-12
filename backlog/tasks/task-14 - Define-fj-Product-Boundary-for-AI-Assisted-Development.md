---
id: TASK-14
title: Define fj Product Boundary for AI-Assisted Development
status: Done
assignee:
  - '@codex'
created_date: '2026-07-12 08:40'
updated_date: '2026-07-12 12:54'
labels: []
dependencies: []
modified_files:
  - >-
    backlog/tasks/task-14 -
    Define-fj-Product-Boundary-for-AI-Assisted-Development.md
ordinal: 93000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Record the fj product boundary for AI-related features. fj remains a Forgejo/Gitea client that supports AI agents through stable commands, deterministic behavior, structured output, clear errors, and safe non-interactive execution; it does not own AI development workflow orchestration. This task records governance only and does not approve or implement commands, dependencies, workflow automation, cross-repository integration, or ROADMAP expansion.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 fj is explicitly defined as a Forgejo/Gitea client rather than an AI workflow orchestrator
- [x] #2 Criteria for including future features in fj core are documented
- [x] #3 Backlog, prompt, session, review, and finalization orchestration are explicitly outside fj core
- [x] #4 TASK-13 templates remain valid as fj repository development assets
- [x] #5 AI-friendliness is defined through structured output, deterministic behavior, clear errors, and safe non-interactive execution
- [x] #6 fj and agent-workflow-kit remain independently governed Sources of Truth
- [x] #7 No cross-repository dependency, automatic synchronization, import, or migration is approved
- [x] #8 A future comparison may occur only after a reviewed external M1 and a non-fj project trial, through a separate fj investigation and Human Decision
- [x] #9 Possible generic Forgejo/Gitea API access and fj configuration or connectivity diagnostics remain separate unapproved investigation candidates
- [x] #10 No TASK-20 or AI workflow assistance subsystem is created
- [x] #11 No code, command, workflow, dependency, external repository, ROADMAP, or existing task implementation scope is changed
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Record the approved fj product-boundary Human Decision in this task as the durable Source of Truth.
2. Record fj core inclusion criteria, explicit exclusions, TASK-13 treatment, cross-repository boundaries, future comparison gates, and allowed future fj investigation candidates.
3. Verify the Decision against PROJECT_CONSTITUTION.md, ARCHITECTURE.md, ROADMAP M7, TASK-8, and TASK-13 without changing those records.
4. Verify that only TASK-14 changed and report the diff and Git state without checking ACs, finalizing, committing, or pushing.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Human Decision — fj Product Boundary for AI-Assisted Development:

Product identity:
- fj is a Forgejo/Gitea client.
- fj may support AI agents through stable commands, deterministic behavior, structured output, clear errors, and safe non-interactive execution.
- fj does not orchestrate AI development workflows and does not own Backlog task lifecycle management, prompt generation, session recovery, review procedures, or finalization procedures.

Product inclusion criteria:
A future fj core feature should normally satisfy at least one of these conditions:
1. It accesses or operates on Forgejo/Gitea resources.
2. It supports fj authentication, instance configuration, connectivity, or Forgejo/Gitea API access.
3. It provides generally useful CLI contracts such as structured output, deterministic behavior, stable errors, or safe non-interactive execution.
A feature being useful only when an AI agent is present is not sufficient by itself to justify inclusion in fj core.

Explicitly outside fj core unless changed by a future explicit Human Decision:
- Reading `backlog/templates/` from the fj binary.
- Template list or show commands.
- Prompt generation.
- Backlog task lifecycle diagnosis.
- Session recovery.
- AI review orchestration.
- Automatic finalization assistance.
- AI provider API calls.
- Provider-specific prompt formatting.
- Autonomous decisions.
- Automatic Acceptance Criteria completion.
- Automatic status transitions.
- Automatic commit, push, tag, or release operations.
These capabilities may remain fj repository documentation, reusable prompt templates, agent skills, external tooling, or part of an independently governed project such as `agent-workflow-kit`. This Decision does not approve fj adoption of any external workflow kit.

TASK-13 template system:
- TASK-13 remains valid.
- The existing templates are approved fj repository development-process assets, part of fj’s current Source of Truth, Markdown-based, provider-neutral, and not runtime fj features.
- fj development continues to use the currently approved fj governance and workflow without waiting for `agent-workflow-kit`.
- This Decision does not replace, import, synchronize, or migrate the TASK-13 templates.

Cross-repository boundary:
- fj and `agent-workflow-kit` are independently governed repositories with separate Sources of Truth.
- No runtime dependency, development dependency, submodule, automatic synchronization, or generated-file relationship is approved.
- Assets must not be copied between repositories without analysis.
- Chat history is not a synchronization mechanism.
- Decisions in one repository do not automatically apply to the other.
- Future adoption or migration requires a separate fj Human Decision.
- The chat-provided current state of `agent-workflow-kit` is contextual only, is not fj repository Evidence, and creates no dependency or waiting condition for fj.

Future comparison candidate:
- After `agent-workflow-kit` reaches a reviewed M1 and has been tested against at least one non-fj project, fj may create a separate read-only investigation task.
- That investigation may compare duplicated workflow rules, generic versus fj-specific templates, possible upstream adoption, compatibility and migration impact, documentation ownership, synchronization policy, and whether no alignment is preferable.
- No comparison task or implementation is approved or created by this Decision. Any adoption or migration requires a separate fj Human Decision.

Allowed future fj investigation candidates:
- Generic Forgejo/Gitea API access may be investigated separately because it could use fj authentication and instance configuration to access Forgejo/Gitea endpoints, comparable in purpose to `gh api`. No command name, syntax, contract, or implementation is approved.
- fj diagnostics for configuration, authentication, instance connectivity, and API compatibility may be investigated separately. Such diagnostics must not diagnose Backlog or AI workflow lifecycle state. No command name or implementation is approved.

ROADMAP impact:
- No ROADMAP change is required.
- M7 remains focused on structured output, deterministic behavior, command discoverability, stable errors, and safe non-interactive execution.
- This Decision does not create TASK-20, an AI Workflow Assistance subsystem, a new AI workflow milestone, or implementation tasks.

Authorization and prohibitions:
- This task authorizes product-boundary documentation only in TASK-14.
- It does not authorize source code, tests, commands, workflows, dependencies, README, LICENSE, ROADMAP, TASK-8, TASK-13, templates, external repository, cross-repository integration, task implementation, commit, push, tag, or release changes.
- The task remains In Progress with Acceptance Criteria unchecked and no Final Summary until a separate review and Finalization workflow.

Model: GPT-5 (Codex) — a cross-cutting product-governance boundary requires careful consistency review across fj’s constitution, architecture, roadmap, and existing AI-experience records.

Pre-Recording Assessment:
- PROJECT_CONSTITUTION.md, ARCHITECTURE.md, ROADMAP M7, TASK-8, and TASK-13 consistency: Pass.
- ROADMAP or existing governing-document change required: No.
- Runtime, compatibility, dependency-direction, and security impact: None; this is a documentation-only governance record.
- Task granularity: one dedicated independently reviewable task is appropriate.
- Additional Human Decision required: No; the complete Human Decision above was explicitly supplied by a human.

Decision Summary:
- fj remains a Forgejo/Gitea-compatible CLI, not an AI development workflow orchestrator.
- AI-agent support is provided through stable commands, structured output, deterministic behavior, clear errors, and safe non-interactive execution.
- Backlog lifecycle orchestration, prompt generation, session recovery, AI review orchestration, automatic finalization, provider integration, autonomous decisions, and automatic repository lifecycle actions remain outside fj core unless changed by a future explicit Human Decision.
- TASK-13 remains fj-owned repository development-process documentation and is not a runtime feature.
- fj and agent-workflow-kit remain independently governed with separate Sources of Truth and no approved dependency, synchronization, import, or migration relationship.

Evidence Summary:
- PROJECT_CONSTITUTION.md defines AI-first as explicit, structured, predictable assistance under human supervision, consistent with this boundary.
- ARCHITECTURE.md defines fj as a Forgejo CLI and supports automation and AI agents through stable command and output boundaries.
- ROADMAP M7 and TASK-8/TASK-8.1 through TASK-8.8 remain limited to structured output, deterministic non-interactive behavior, operation context, compatibility, errors, and command discovery.
- TASK-13 and backlog/templates/README.md confirm that the existing templates are repository-owned Markdown development-process assets, separate from fj runtime behavior.
- TASK-14 records core inclusion criteria, explicit exclusions, cross-repository separation, future comparison gates, and unapproved future fj investigation candidates.
- Git verification shows only the new TASK-14 file; no existing document, task, code, workflow, dependency, ROADMAP item, or external repository was changed.

Independent Review:
- Scope: PROJECT_CONSTITUTION.md, ARCHITECTURE.md, ROADMAP.md, AGENTS.md, README.md, TASK-8 and TASK-8.1 through TASK-8.8, TASK-13, backlog/templates/README.md, TASK-14, and current Git status and diff.
- External repository access: Not performed.
- Acceptance Criteria #1-#11: Pass.
- Critical findings: None.
- Major findings: None.
- Minor findings: None.
- Product boundary clarity, durability, flexibility, future-maintainer usability, and compatibility with fj identity: Pass.
- TASK-13 preservation and cross-repository governance separation: Pass.
- ROADMAP unchanged assessment: Pass.
- Result: Ready for Finalization.

Deferred Items:
- Generic Forgejo/Gitea API access may be considered only through a separate investigation and approval.
- fj configuration, authentication, connectivity, and API compatibility diagnostics may be considered only through a separate investigation and approval.
- Any future comparison with agent-workflow-kit requires a reviewed external M1, non-fj validation, a separate read-only fj investigation, and a separate fj Human Decision before adoption or migration.
- No alignment remains an acceptable outcome of any future comparison.
- Any change to the product boundary requires a future explicit Human Decision.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
TASK-14 established fj’s durable product boundary for AI-assisted development. fj remains a Forgejo/Gitea-compatible CLI and supports people, automation, and AI agents through stable commands, structured output, deterministic behavior, clear errors, and safe non-interactive execution. It does not own AI development workflow orchestration.

The Decision preserves TASK-13 as fj repository development-process documentation while keeping it outside runtime fj functionality. It also establishes independent governance and separate Sources of Truth for fj and agent-workflow-kit, with no approved dependency, synchronization, import, migration, or automatic decision inheritance.

Future comparison and possible fj capabilities such as generic Forgejo/Gitea API access or fj diagnostics remain separately gated, unapproved investigation candidates. ROADMAP M7 remains unchanged and no AI workflow subsystem or implementation task was created.

Independent Review confirmed Acceptance Criteria #1-#11 with no Critical, Major, or Minor findings. TASK-14 is Ready for Finalization.
<!-- SECTION:FINAL_SUMMARY:END -->
