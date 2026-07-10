# CLAUDE.md

This document provides Claude Code-specific guidance that supplements [AGENTS.md](AGENTS.md). It must never override [PROJECT_CONSTITUTION.md](PROJECT_CONSTITUTION.md) or `AGENTS.md`.

## Subagent Usage

- Use subagents when parallel or independent work adds value: research, design comparison, impact analysis, test analysis, and independent reviews.
- Do not launch subagents for simple work that the main session can complete directly.
- Never adopt subagent results unconditionally; the main agent must integrate and verify them before use.

## Model Selection

Follow the model-selection and recording rules in AGENTS.md Section 14. In Claude Code:

- Use Fable 5 for design work, important reviews, and verification of major changes as defined in AGENTS.md Section 15.
- Small, well-defined implementation work may use a lighter model.
- Decide the model before starting the task, and record the model and the reason in the Backlog task following AGENTS.md Section 14, including its check for an official CLI field.

## Major-Change Reviews

For major changes (AGENTS.md Section 15), perform both reviews with Fable 5:

- The pre-implementation check is requested from an agent or subagent running Fable 5.
- The post-implementation review is an independent review by an agent or subagent running Fable 5; if the main session produced the implementation, use a separate Fable 5 subagent so the review stays independent.
- Classify findings and handle fixes according to the Review Finding Classification in AGENTS.md Section 15.
- Report the results to the human. Never commit before human approval.

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
