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

1. Work on only one task at a time.
2. Read the relevant documentation, code, tests, and recent changes.
3. Confirm the task scope, constraints, and approval state.
4. Propose a design when required and wait for approval.
5. Implement only the approved scope in small, reviewable changes.
6. Update affected documentation and tests.
7. Run `make pre-commit` when the Makefile is available, and report the results,
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

1. Read `PROJECT_CONSTITUTION.md` first.
2. Read `ARCHITECTURE.md`.
3. Read `ROADMAP.md`.
4. Run `backlog instructions overview`.
5. Use Backlog.md to determine the next unfinished task.
6. Read only the documents required for that task.
7. Never continue work from chat history alone.
8. Treat the repository as the single source of truth for project state and continuity.
9. If the previous task is already Done, automatically continue with the next task unless the user explicitly requests otherwise.
10. Before implementing, summarize the task that will be worked on.
11. Select tasks according to Backlog.md status and dependencies, not task numbers.
12. When multiple tasks are available, choose the highest-priority task according to Backlog.md.
13. Treat task IDs as identifiers only; never use them as execution order.

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

- Before starting a task, choose a model appropriate for the work: use the most capable model available on the agent platform for design, significant reviews, and verification of major changes; small, well-defined implementation work may use a lighter model.
- Platform-specific model names and dispatch mechanics belong in platform files such as `CLAUDE.md`, not here.
- Record the model used and the reason for choosing it in the Backlog task.
- Before recording, check `backlog task edit --help` for an official model field. If one exists, use it. If none exists, do not invent custom metadata fields or edit task files directly; propose recording the model in the task's implementation plan or notes and wait for human approval of that method.
- Approved recording method (the current CLI has no model field): append one line to the task's Implementation Notes through the `backlog` CLI in the format `Model: <model name> — <reason>`, for example `backlog task edit TASK-123 --append-notes "Model: Fable 5 — architecture-sensitive design"`. Do not add custom Backlog fields. If a future CLI version introduces an official model field, prefer it per the rule above.

## 15. Major Changes

A change is a major change when it involves any of the following:

- modifying `PROJECT_CONSTITUTION.md`;
- changing important rules in `AGENTS.md`;
- modifying `ARCHITECTURE.md`;
- modifying `ROADMAP.md`;
- changing public CLI commands, flags, output formats, JSON contracts, or exit codes;
- changing a public API or compatibility guarantee;
- adding a new external dependency;
- changing package structure or dependency direction;
- changing authentication, credential handling, or security boundaries;
- affecting multiple milestones;
- large-scale refactoring.

Judge by user impact, compatibility, security, and design boundaries, not by line or file counts alone. When unsure whether a change is major, treat it as major and ask a human.

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

AI may present findings at any severity. For Major, Minor, and Suggestion findings, AI must not implement fixes on its own; a human decides for each finding whether it is adopted, deferred, or rejected before any fix is made. Critical findings are reported to a human with the highest priority. AI may propose a fix for a Critical finding and, where necessary, implement it, but it must then stop and wait for human confirmation instead of continuing the task or advancing to the next phase. Section 10 governs all commits and pushes regardless of severity.

## 16. Stop-and-Report Rule

After completing one task, stop. Do not start the next task automatically. Selecting the next task is governed by the Session Resume Workflow in Section 12 and begins with the next human work request.

Before stopping, report:

- the selected task and why it was selected;
- the model used and why;
- whether subagents were used;
- the files changed;
- what was implemented;
- the tests executed and their results;
- unverified items, constraints, and remaining work;
- whether the pre-implementation check and post-implementation review in Section 15 were performed, and their results;
- a recommended commit message.

Do not commit or push; wait for human confirmation, per Section 10.

## 17. Backlog Task Lifecycle

AI agents must follow this lifecycle for every implementation task:

```text
To Do → Decision approved → In Progress → Review → Done
```

- Human approval of the implementation Decision is required before implementation starts.
- Every Acceptance Criterion must be checked (`[x]`) before a task is changed to Done.
- Status and Acceptance Criteria must describe the same completion state.

### Ownership

- Assignee is optional while a task is To Do or Decision approved.
- An Assignee is required when implementation starts.
- The Assignee records the implementation owner or implementation agent.

### Completion Checklist

A task may be marked Done only when all of the following are true:

- Implementation completed
- Tests passed
- Required verification passed
- Independent review completed
- Review summary recorded
- All Acceptance Criteria checked
- Implementation Notes updated
- Final Summary recorded
- Status updated to Done
- No unapproved scope remains

## Task Lifecycle Workflow

After implementation, agents follow this workflow for the current task only:

```text
Verification → Independent Review → Review Ready → Finalization → Task Commit → Push Confirmation
```

- When Independent Review is `Review Ready`, meaning no unresolved Critical or
  Major finding and no missing human decision remains, proceed automatically to
  Backlog Finalization for the same task.
- Critical or Major findings require an immediate stop and human decision.
- Minor findings require a proposed fix or deferral and the human's adopt,
  defer, or reject decision before any follow-up change.
- Suggestions are recorded as future improvement candidates and do not block
  Finalization unless a human decision is still required.
- A commit must contain changes for one task only. Never mix unrelated changes
  into a task commit; inspect the staged file list and cached diff first.
- Before push, confirm the commit hash, commit scope, task status, and that no
  unrelated changes are included.
- Commit and push remain explicitly authorized operations under Section 10.
- Destructive operations, including history rewriting, force-push, deletion, or
  privacy cleanup, require separate human approval and an approved task.

## 18. Autonomous Workflow Enforcement

- An agent must not implement a task unless its dependencies are complete and its approved Decision is recorded in Backlog.
- When multiple tasks are executable, select the highest-priority task; use ordinal order as the tie-breaker.
- An agent must stop when no executable task exists, when an approval decision is unresolved, or when the approved scope cannot be maintained.
- Public CLI/API contracts, security boundaries, dependency direction, external dependencies, and roadmap or milestone changes require human approval before implementation.
- A task must not be marked Done with unchecked Acceptance Criteria, missing Verification, missing Independent Review, missing Final Summary, or unresolved Critical/Major findings.
- Minor findings must be fixed or explicitly recorded as deferred; Suggestions must be recorded when they affect future work.
- Model selection and its rationale must be recorded in the task's Implementation Notes.

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
