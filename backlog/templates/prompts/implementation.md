# Implementation Prompt

Implement the approved scope for `{{TASK_ID}} — {{TASK_TITLE}}`.

## Context

- Task file: `{{TASK_FILE}}`
- Objective: {{OBJECTIVE}}
- Scope: {{SCOPE}}
- Human Decision:
{{HUMAN_DECISION}}
- Acceptance Criteria:
{{ACCEPTANCE_CRITERIA}}
- Allowed files: {{ALLOWED_FILES}}
- Prohibited changes: {{PROHIBITED_CHANGES}}
- Verification commands:
{{VERIFICATION_COMMANDS}}
- Deferred Items:
{{DEFERRED_ITEMS}}

## Rules

- Confirm dependencies, Decision record, assignee, and In Progress state before mutation.
- Follow repository instructions and perform any required pre-implementation check.
- Make the smallest change that satisfies the approved scope; do not perform opportunistic cleanup.
- Preserve unrelated user changes and keep Deferred Items outside current scope.
- Verify proportionately using the listed commands and repository-required checks.
- Do not claim your work is independently reviewed.
- Do not check ACs, add a Final Summary, set Done, commit, push, tag, release, or run external workflows unless separately authorized.

## Output

Report the approved Decision, files changed, implementation, verification and results, limitations, Deferred Items, full diff summary, `git status --short`, and proposed commit message.
