# Remediation Prompt

Remediate only the approved findings for `{{TASK_ID}} — {{TASK_TITLE}}`.

## Context

- Task file: `{{TASK_FILE}}`
- Objective: {{OBJECTIVE}}
- Scope: {{SCOPE}}
- Human Decision:
{{HUMAN_DECISION}}
- Findings to address:
{{FINDINGS}}
- Allowed files: {{ALLOWED_FILES}}
- Prohibited changes: {{PROHIBITED_CHANGES}}
- Verification commands:
{{VERIFICATION_COMMANDS}}
- Deferred Items:
{{DEFERRED_ITEMS}}

## Rules

- Confirm the human disposition for each finding before making changes.
- Address only the listed findings; do not add unrelated cleanup, refactoring, or improvements.
- Do not turn Deferred Items into current requirements or expand scope without approval.
- Preserve unrelated changes and verify the remediation and regression-sensitive areas.
- Do not check ACs, add a Final Summary, set Done, finalize, commit, push, tag, or release.
- Do not treat implementation verification as Independent Re-Review.

## Output

Report Findings Being Addressed, Explicitly Out of Scope, Root Cause, Remediation Plan, Changes Made, Verification, Remaining Findings, Git Diff Summary, `git status --short`, and Re-Review Request.
