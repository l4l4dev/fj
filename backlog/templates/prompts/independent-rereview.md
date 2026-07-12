# Independent Re-Review Prompt

Perform a read-only independent re-review of remediation for `{{TASK_ID}} — {{TASK_TITLE}}`.

## Context

- Task file: `{{TASK_FILE}}`
- Objective: {{OBJECTIVE}}
- Scope: {{SCOPE}}
- Acceptance Criteria:
{{ACCEPTANCE_CRITERIA}}
- Original findings:
{{FINDINGS}}
- Allowed files to review: {{ALLOWED_FILES}}
- Prohibited changes: {{PROHIBITED_CHANGES}}
- Deferred Items:
{{DEFERRED_ITEMS}}

## Rules

- Confirm reviewer independence; the implementer must not approve their own remediation.
- Classify each original finding as Resolved, Partially Resolved, or Not Resolved with Evidence.
- Reassess every affected AC as Pass, Fail, or Insufficient Evidence.
- Check regressions, scope boundaries, Decision preservation, and Deferred Items.
- Classify any new findings under repository severity rules.
- Do not change files, Backlog, status, ACs, Git state, or external systems.
- Do not commit or push.
- Ready requires all ACs to pass, no Critical or Major finding, and no unresolved Human Decision.

## Output

Report Review Scope, Original Findings, Remediation Reviewed, Finding Resolution, AC Reassessment, Regression / Scope Check, Findings, Finalization Decision, and—only when Ready—a Proposed Finalization Record.
