# CLAUDE.md

This document provides Claude Code-specific guidance that supplements [AGENTS.md](AGENTS.md). It must never override [PROJECT_CONSTITUTION.md](PROJECT_CONSTITUTION.md) or `AGENTS.md`.

## Subagent Usage

- Use subagents when parallel or independent work adds value: research, design comparison, impact analysis, test analysis, and independent reviews.
- Do not launch subagents for simple work that the main session can complete directly.
- Never adopt subagent results unconditionally; the main agent must integrate and verify them before use.

## Model Selection

Follow AGENTS.md Section 14. In Claude Code, the session's default model is
fine for routine tasks; use the most capable available model (currently
Fable 5) for major changes and their reviews. Do not record model choices per
task.

## Major-Change Reviews

Only for major changes (AGENTS.md Section 15):

- Request the pre-implementation check and the post-implementation review from Fable 5; if the main session produced the implementation, use a separate Fable 5 subagent so the review stays independent.
- Classify findings per AGENTS.md Section 15 and report the results to the human. Never commit before human approval.

Routine tasks need neither review; they follow the single loop in DEVELOPMENT_WORKFLOW.md.

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
