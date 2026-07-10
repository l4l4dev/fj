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
7. Run relevant verification and report the results, limitations, and remaining work.

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
