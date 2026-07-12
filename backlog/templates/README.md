# Backlog Workflow Templates

## Purpose

These templates standardize evidence-based, human-directed task work across sessions. They help maintainers, contributors, and AI agents preserve enough context in the repository to resume work without relying on chat history.

`AGENTS.md` and approved project documents remain authoritative. These templates organize work and records; they do not override approval gates, task lifecycle rules, or human instructions.

## Template Types

- **Record templates** in this directory define the durable information to record for each workflow stage.
- **Prompt templates** in [`prompts/`](prompts/) are copy-ready English instructions for an AI agent. Replace every `{{PLACEHOLDER}}` before use.

Do not treat an AI response or chat transcript as durable evidence. Record Decisions, Evidence, review results, Deferred Items, and Final Summaries in the relevant Backlog task or another approved repository record through the supported Backlog CLI.

## Standard Workflow

```text
Investigation
    ↓
Human Decision (when required)
    ↓
Implementation
    ↓
Independent Review
    ├─ Ready → Finalization
    └─ Changes Required
           ↓
       Remediation
           ↓
       Independent Re-Review
           ├─ Ready → Finalization
           └─ Changes Required → Remediation
    ↓
Commit (only after explicit authorization)
```

See [`workflow.md`](workflow.md) for stage entry conditions, outputs, prohibited actions, and transitions.

## Files and When to Use Them

| Stage | Record template | AI prompt | Use |
| --- | --- | --- | --- |
| Investigation | [`00-investigation.md`](00-investigation.md) | [`prompts/investigation.md`](prompts/investigation.md) | Gather facts, options, risks, and decision needs without implementation. |
| Human Decision | [`01-decision.md`](01-decision.md) | [`prompts/decision-recording.md`](prompts/decision-recording.md) | Persist an explicit human choice and implementation authorization. |
| Implementation | [`02-implementation.md`](02-implementation.md) | [`prompts/implementation.md`](prompts/implementation.md) | Plan, implement, and verify only the approved scope. |
| Independent Review | [`03-independent-review.md`](03-independent-review.md) | [`prompts/independent-review.md`](prompts/independent-review.md) | Evaluate Acceptance Criteria and findings without modifying files. |
| Remediation | [`04-remediation.md`](04-remediation.md) | [`prompts/remediation.md`](prompts/remediation.md) | Address only findings selected by a human. |
| Independent Re-Review | [`05-independent-rereview.md`](05-independent-rereview.md) | [`prompts/independent-rereview.md`](prompts/independent-rereview.md) | Verify finding resolution and regressions independently. |
| Finalization | [`06-finalization.md`](06-finalization.md) | [`prompts/finalization.md`](prompts/finalization.md) | Record completed Evidence and review, check ACs, and transition to Done. |

## Required and Optional Stages

- **Implementation, proportionate verification, review required by `AGENTS.md`, and Finalization** apply to implementation tasks.
- **Investigation** may be combined with planning for a small, obvious task when no meaningful uncertainty exists.
- **Human Decision** is required whenever an approval gate or material choice exists. Record it before implementation.
- **Remediation and Independent Re-Review** are used only when review findings require an approved response.
- **Commit, push, tag, release, and workflow execution** are not implied by any template. They require separate explicit authorization.

## Simplifying Simple Tasks

A typo, broken link, or similarly small and reversible documentation correction may omit standalone Investigation and Decision records when all of the following are true:

- the requested outcome and allowed files are unambiguous;
- no public behavior, compatibility, security, legal, governance, architecture, dependency, or migration decision is involved;
- no scope expansion is needed;
- project rules do not classify the change as major or require an independent review.

Even when simplified, preserve scope, verify the change proportionately, report the diff and Git state, and do not commit or push without authorization.

## Human Decision Gate

AI agents must not decide material matters such as licensing, public releases, breaking changes, architecture, security policy, API compatibility, data migration, legal or governance policy, new external dependencies, or major scope changes. Present options and trade-offs, obtain a human decision, and persist it through Backlog before implementation.

## Example

1. Copy the relevant prompt from `prompts/`.
2. Replace all placeholders with the current task, scope, ACs, constraints, and verification commands.
3. Run the stage and review its output.
4. Persist the approved result with the corresponding record template through the Backlog CLI.
5. Move to the next stage only when its entry conditions are satisfied.
