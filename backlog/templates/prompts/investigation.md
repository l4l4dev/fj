# Investigation Prompt

Investigate `{{TASK_ID}} — {{TASK_TITLE}}` in read-only mode.

## Context

- Task file: `{{TASK_FILE}}`
- Objective: {{OBJECTIVE}}
- Scope: {{SCOPE}}
- Acceptance Criteria:
{{ACCEPTANCE_CRITERIA}}
- Deferred Items already known:
{{DEFERRED_ITEMS}}

## Constraints

- Allowed files to read: {{ALLOWED_FILES}}
- Prohibited changes: {{PROHIBITED_CHANGES}}
- Do not modify files, Backlog, Git state, or external systems.
- Do not commit or push.
- Treat the repository as the Source of Truth; do not use chat history as durable Evidence.
- Do not make a Human Decision. Identify whether one is required.

## Work

1. Read governing repository instructions and relevant durable records.
2. Establish current state and Evidence.
3. Compare viable options, trade-offs, risks, compatibility, security, and scope impact.
4. Identify open questions and decisions that only a human may make.
5. Recommend a next step without implementing it.

## Output

Report: Task, Objective, Scope, Current State, Files / Components Reviewed, Constraints, Options, Evidence, Risks, Open Questions, Human Decision Required, AI Recommendation, and Proposed Next Step.
