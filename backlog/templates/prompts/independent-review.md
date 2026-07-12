# Independent Review Prompt

Perform a read-only independent review of `{{TASK_ID}} — {{TASK_TITLE}}`.

## Context

- Task file: `{{TASK_FILE}}`
- Objective: {{OBJECTIVE}}
- Scope: {{SCOPE}}
- Human Decision:
{{HUMAN_DECISION}}
- Acceptance Criteria:
{{ACCEPTANCE_CRITERIA}}
- Allowed files to review: {{ALLOWED_FILES}}
- Prohibited changes: {{PROHIBITED_CHANGES}}
- Deferred Items:
{{DEFERRED_ITEMS}}

## Rules

- Confirm reviewer independence from implementation.
- Read the repository, task record, diff, verification Evidence, and relevant governing documents.
- Assess every AC as Pass, Fail, or Insufficient Evidence with Evidence and Notes.
- Classify findings as Critical, Major, Minor, or Suggestion under repository rules.
- Check approved scope, architecture, compatibility, security, documentation, and Deferred Item boundaries.
- Do not change files, Backlog, status, ACs, Git state, or external systems.
- Do not commit or push.
- Ready for Finalization requires every AC to pass, no Critical or Major finding, and no unresolved Human Decision. Minor findings require human disposition before a fix.

## Output

Report Review Scope, Reviewer Independence, Files / Evidence Reviewed, AC Assessment, Findings, Scope Boundary Assessment, Deferred Items Assessment, Finalization Decision, Required Remediation, and—only when Ready—a Proposed Finalization Record.
