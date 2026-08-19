# AGENTS.md

This file defines how AI agents must work in this repository. It applies to Codex, Claude Code, and future agents.

## 1. Mission

Help build and maintain `fj` safely, transparently, and under human direction. Follow the principles in [PROJECT_CONSTITUTION.md](PROJECT_CONSTITUTION.md); do not restate or reinterpret them here.

## 2. Repository Structure

- `PROJECT_CONSTITUTION.md`: highest-level project principles and policies.
- `ARCHITECTURE.md`: approved high-level architecture and dependency boundaries.
- `ROADMAP.md`: milestone-based direction for project evolution.
- `CLAUDE.md`: Claude Code-specific guidance that supplements this file.
- `README.md`: project introduction and user-facing overview.
- `go.mod`: Go module definition.
- `AGENTS.md`: operational rules for AI agents.
- `backlog/`: Backlog.md-managed project tasks and metadata; modify only through the `backlog` CLI.

Inspect the current tree before every task. Do not invent new structure or create files unless the approved task requires it.

## 3. Source of Truth

Use this precedence order:

1. Human instructions
2. `PROJECT_CONSTITUTION.md`
3. Approved design documents
4. `AGENTS.md`
5. Task definitions
6. Implementation

Task definitions must never override `PROJECT_CONSTITUTION.md`.

If sources conflict or requirements are unclear, stop and ask for clarification. Never resolve ambiguity by assumption.

## 4. Development Workflow

For the task workflow reference, see [`backlog/templates/README.md`](backlog/templates/README.md).

1. Work on only one task at a time.
2. Read the task and the documentation, code, and tests it touches.
3. Present a short implementation plan and wait for human approval before
   implementing. A plan-mode approval or an explicit chat approval is
   sufficient; record it in the task in one line.
4. Implement only the approved scope in small, reviewable changes, with tests.
5. Update affected documentation.
6. Run `make pre-commit` when the Makefile is available, and report the results,
   limitations, and remaining work.

## 5. Task Execution Rules

- Never implement work that is not part of the approved task.
- Do not perform opportunistic refactoring, cleanup, or feature work.
- Preserve unrelated user changes and do not modify out-of-scope files.
- Prefer the smallest change that fully satisfies the task.
- State assumptions only when they are explicitly confirmed; otherwise ask for clarification.
- Report blockers promptly and do not conceal uncertainty or failed verification.

## 6. Approval Gates

Human approval is required before:

- moving from a requested design phase to implementation;
- expanding or changing the approved scope;
- making destructive, irreversible, security-sensitive, or externally visible changes;
- changing public behavior, compatibility guarantees, or project-wide architecture;
- modifying `PROJECT_CONSTITUTION.md`.

Approval for one task or phase does not imply approval for another.

## 7. Documentation Rules

- Update documentation whenever behavior changes.
- Keep documentation consistent with implemented and verified behavior.
- Record important decisions with their rationale and impact.
- Do not treat chat history or agent memory as project documentation.
- Do not change project philosophy through incidental documentation edits.

## 8. Coding Rules

- Follow the Go version and module settings declared in `go.mod`.
- Prefer clear, explicit, idiomatic Go over speculative abstraction.
- Keep responsibilities narrow and dependencies minimal.
- Preserve deterministic behavior, structured interfaces, and actionable errors.
- Maintain backward compatibility unless an incompatible change is explicitly approved.
- Do not expose secrets, credentials, personal data, or non-public information.

## Repository Privacy and Recording Hygiene

- In Git history, commit messages, Backlog tasks, README, documentation, and
  `.agent` records, use placeholders instead of real personal names,
  organization names, hostnames, or repository owner names by default.
- Even externally published information may be recorded only with explicit
  user approval or a clearly authorized publication purpose.
- Use these placeholders in execution, verification, and acceptance records:
  - `example-owner`
  - `example-repository`
  - `https://forgejo.example.com`
- Never record credential values, raw tokens, or credentials embedded in URLs.
- Git history rewriting, force-push, and privacy cleanup are prohibited unless
  covered by a separately approved task.

## 9. Testing Rules

- Add or update tests for every behavior change.
- Run the smallest relevant test set during development and the broader applicable suite before completion.
- Test success, failure, and boundary cases relevant to the task.
- Never claim tests passed unless they were executed successfully.
- Report skipped tests, failures, environment limitations, and unverified behavior.

## 10. Commit Rules

- Commit only when explicitly requested.
- Keep one commit focused on one purpose.
- Use Conventional Commits unless the human provides an exact commit message.
- Stage only approved files and inspect the staged diff before committing.
- Never include unrelated changes or generated artifacts accidentally.
- Never push unless explicitly requested.

## 11. Definition of Done

A task is done only when:

- the approved scope is fully implemented and no out-of-scope work was added;
- the change is small, reviewable, and consistent with `PROJECT_CONSTITUTION.md`;
- relevant tests and checks pass, or any limitations are clearly reported;
- behavior changes are reflected in documentation;
- changed files and user-owned differences are accurately reported;
- no required approval, decision, or follow-up remains unresolved.

## 12. Session Resume Workflow

When resuming work in this repository:

1. Run `backlog instructions overview`.
2. Use Backlog.md to determine the next unfinished task.
3. Read only the documents required for that task; consult
   `PROJECT_CONSTITUTION.md`, `ARCHITECTURE.md`, and `ROADMAP.md` when the
   task touches principles, architecture boundaries, or roadmap direction.
4. Never continue work from chat history alone.
5. Treat the repository as the single source of truth for project state and continuity.
6. If the previous task is already Done, automatically continue with the next task unless the user explicitly requests otherwise.
7. Before implementing, summarize the task that will be worked on.
8. Select tasks according to Backlog.md status and dependencies, not task numbers; when multiple tasks are available, choose the highest-priority task.

## 13. Human Approval Boundaries

- AI may choose the next task based on Backlog.md status, dependencies, and priority.
- AI must never change task priority.
- AI must never change dependencies.
- AI must never create, delete, or reorder roadmap items.
- AI must never modify `PROJECT_CONSTITUTION.md` without explicit human approval.
- AI must never modify `ROADMAP.md` without explicit human approval.
- AI must never modify Backlog priorities or milestones without explicit human approval.
- AI may implement, test, update Backlog task status, and stop after completing one task.
- AI may select the next executable task based on Backlog status, dependencies, and priority, and may update task plans, progress, verification results, and completion state through the `backlog` CLI.
- AI may request the pre-implementation checks and independent reviews required by Section 15.
- AI must never finalize the design of a major change (Section 15) without human approval.
- AI must never commit or push without an explicit human request, per Section 10.
- AI must never advance to a next phase based only on agreement between AI agents; phase transitions require the human approvals defined in Section 6.

## 14. Model Selection

- Use the platform's default capable model for routine work. Use the most
  capable model available on the platform for major changes (Section 15) and
  their reviews.
- Platform-specific model names and dispatch mechanics belong in platform files such as `CLAUDE.md`, not here.
- Do not record model choices per task.

## 15. Major Changes

A change is a major change when it involves any of the following:

- modifying `PROJECT_CONSTITUTION.md`;
- changing important rules in `AGENTS.md`;
- modifying `ARCHITECTURE.md`;
- modifying `ROADMAP.md`;
- breaking or incompatible changes to existing public CLI commands, flags,
  output formats, JSON contracts, exit codes, or other compatibility
  guarantees;
- adding a new external dependency;
- changing package structure or dependency direction;
- changing authentication, credential handling, or security boundaries;
- affecting multiple milestones;
- large-scale refactoring.

Adding a new command, subcommand, or flag that follows existing architecture
and CLI patterns is routine work, not a major change.

Judge by user impact, compatibility, security, and design boundaries, not by line or file counts alone. When unsure whether a change is major, treat it as major and ask a human.

The pre-implementation check and post-implementation review below apply only
to major changes. Routine tasks need neither; they follow the workflow in
Section 4.

### Pre-Implementation Check

Before starting a major change, request an assessment from an agent or subagent using the most capable model available on the platform. The assessment must evaluate:

- consistency with `PROJECT_CONSTITUTION.md`, `ARCHITECTURE.md`, and `ROADMAP.md`;
- scope of impact;
- backward compatibility;
- security;
- alternatives;
- whether the task granularity is appropriate;
- whether additional human approval is required.

If the assessment surfaces significant decisions, do not start implementation; wait for human approval. This check supplements, and never replaces, the approval gates in Section 6.

### Post-Implementation Review

After implementation and tests are complete, obtain an independent review by an agent or subagent using the most capable available model. The review must confirm:

- only the approved scope was changed;
- acceptance criteria are satisfied;
- architecture boundaries are preserved;
- no compatibility or security problems were introduced;
- tests are sufficient;
- documentation matches the implementation;
- Backlog records are accurate.

Report the review result to a human. Never commit before human approval.

### Review Finding Classification

Classify every review finding as one of four severities:

- **Critical:** must be resolved before the change can proceed.
- **Major:** a significant problem that requires a human decision.
- **Minor:** a small defect or inconsistency.
- **Suggestion:** an optional improvement.

Critical and Major findings require an immediate stop and a human decision;
AI may propose or implement a fix for a Critical finding but must then wait
for human confirmation. Minor findings may be fixed within the approved scope
without a separate decision, or recorded as deferred. Suggestions are recorded
when they affect future work. Section 10 governs all commits and pushes
regardless of severity.

## 16. Stop-and-Report Rule

After completing one task, stop. Do not start the next task automatically. Selecting the next task is governed by the Session Resume Workflow in Section 12 and begins with the next human work request.

Before stopping, report:

- the task worked on and what was implemented;
- the files changed;
- the tests executed and their results;
- unverified items and remaining work;
- for major changes, the Section 15 check and review results;
- a recommended commit message.

Do not commit or push; wait for human confirmation, per Section 10.

## 17. Backlog Task Lifecycle

AI agents must follow this lifecycle for every implementation task:

```text
To Do → In Progress (after plan approval) → Done
```

- Human approval of the implementation plan (Section 4) is required before implementation starts; record the approval in the task in one line.
- An Assignee is required when implementation starts.
- A task may be marked Done only when: the approved scope is fully implemented,
  tests pass, every Acceptance Criterion is checked (`[x]`), a short
  implementation note and final summary are recorded, and no unapproved scope
  remains. For major changes, the Section 15 review must also be complete with
  no unresolved Critical or Major finding.
- A commit must contain changes for one task only; inspect the staged diff
  first. Commit and push remain explicitly authorized operations under
  Section 10. Destructive operations, including history rewriting, force-push,
  deletion, or privacy cleanup, require separate human approval and an
  approved task.

## 18. Autonomous Workflow Enforcement

- An agent must not implement a task unless its dependencies are complete and its plan is approved.
- When multiple tasks are executable, select the highest-priority task; use ordinal order as the tie-breaker.
- An agent must stop when no executable task exists, when an approval decision is unresolved, or when the approved scope cannot be maintained.
- Major changes (Section 15) require human approval before implementation.
- A task must not be marked Done with unchecked Acceptance Criteria, failing or unrun tests, or unresolved Critical/Major findings.

<!-- BACKLOG.MD GUIDELINES START -->
<CRITICAL_INSTRUCTION>

## Backlog.md Workflow

This project uses Backlog.md for task and project management.

**For every user request in this project, run `backlog instructions overview` before answering or taking action.**

Use the overview to decide whether to search, read, create, or update Backlog tasks.

Use the detailed guides when needed:
- `backlog instructions task-creation` for creating or splitting tasks
- `backlog instructions task-execution` for planning and implementation workflow
- `backlog instructions task-finalization` for completion and handoff

Use `backlog <command> --help` before running unfamiliar commands. Help shows options, fields, and examples.

Do not edit Backlog task, draft, document, decision, or milestone markdown files directly. Use the `backlog` CLI so metadata, relationships, and history stay consistent.

</CRITICAL_INSTRUCTION>
<!-- BACKLOG.MD GUIDELINES END -->
