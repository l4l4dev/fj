# Finalization Prompt

Finalize `{{TASK_ID}} — {{TASK_TITLE}}` only because the latest independent review is Ready for Finalization.

## Context

- Task file: `{{TASK_FILE}}`
- Objective: {{OBJECTIVE}}
- Scope: {{SCOPE}}
- Human Decision:
{{HUMAN_DECISION}}
- Acceptance Criteria:
{{ACCEPTANCE_CRITERIA}}
- Review findings and result:
{{FINDINGS}}
- Deferred Items:
{{DEFERRED_ITEMS}}
- Allowed files: {{ALLOWED_FILES}}
- Prohibited changes: {{PROHIBITED_CHANGES}}
- Verification commands:
{{VERIFICATION_COMMANDS}}

## Preconditions

- Every AC is Pass.
- No Critical or Major finding remains.
- Every Minor finding has a recorded human disposition.
- No Human Decision is unresolved.
- Implementation and required verification are complete.

## Rules

- Use the supported Backlog CLI for task changes.
- Check only ACs supported by recorded Evidence.
- Record Decision Summary, Evidence Summary, independent review, resolved findings, Deferred Items, and Final Summary.
- Transition to Done only after the task record is internally consistent.
- Finalization is not implementation: do not add features, fixes, refactors, or unrelated documentation.
- Do not commit, push, tag, release, or run external workflows; those require separate explicit authorization.

## Output

Report changed files, AC final state, Decision and Evidence summaries, review record, Deferred Items, Final Summary, status and assignee, verification results, `git diff --check`, `git diff --stat`, `git diff`, `git status --short`, and proposed commit message.
