# Decision Recording Prompt

Record the already approved Human Decision for `{{TASK_ID}} — {{TASK_TITLE}}` before implementation.

## Context

- Task file: `{{TASK_FILE}}`
- Objective: {{OBJECTIVE}}
- Scope: {{SCOPE}}
- Human Decision:
{{HUMAN_DECISION}}
- Allowed files: {{ALLOWED_FILES}}
- Prohibited changes: {{PROHIBITED_CHANGES}}

## Rules

- Record only the Human Decision supplied above; do not reconsider, expand, or replace it.
- Clearly separate any prior AI recommendation from the Human Decision.
- Record context, options considered, trade-offs, rationale, consequences, Deferred Items, implementation authorization, allowed files, and prohibited changes when supplied.
- Use the supported Backlog CLI. Do not edit Backlog-managed task files directly.
- Do not implement repository changes in this stage.
- Do not check ACs, add a Final Summary, set Done, commit, or push.

## Output

Report the Backlog record changed, the exact Decision recorded, implementation authorization state, status and assignee state, prohibited actions, Git diff summary, and `git status --short`.
