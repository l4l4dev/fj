# Task Preflight Template

## Task Classification

Classify every candidate before selecting it.

### Leaf Implementation Task

Conditions:

- Has no subtasks.
- Has a clear implementation scope.
- Its Decision is approved.
- Its Acceptance Criteria are implementable.

Handling:

- May become an implementation candidate.

If the Decision is not approved:

- Do not start implementation.
- Use `.agent/decision-plan.md` to organize unresolved decisions.
- Stop and wait for human approval.
- Proceed to implementation only after the Decision is approved.

### Parent / Epic Task

Conditions:

- Has subtasks.
- Exists to aggregate multiple tasks.

Handling:

- Exclude from normal implementation candidates.

### Finalization Task

Conditions:

- Is a Parent Task.
- All subtasks are Done.
- Its purpose is metadata synchronization rather than implementation.

Handling:

- Treat as a Backlog synchronization task, not an implementation task.

### Documentation / Process Task

Conditions:

- Its purpose is changing Markdown, workflow rules, templates, or other process documentation.

Handling:

- Keep it separate from source implementation tasks and apply documentation-specific verification.

## Executable Task Rule

A task is an implementation candidate only when all conditions below are true:

- It is a Leaf Implementation Task.
- All dependencies are complete.
- Its Decision is approved.
- Its Acceptance Criteria are sufficient.
- No Critical or Major review finding remains unresolved.
- Its scope is explicit.

## Selection Rule

Select one task in this order:

1. Critical fix.
2. Blocking Finalization Task.
3. Leaf Implementation Task.
4. Documentation / Process Task.

Parent / Epic Tasks must not be selected as implementation targets.

Within the same category, use priority first and ordinal as the tie-breaker.

## No Candidate

When no task satisfies the applicable rule:

- Do not start implementation.
- Report `Ready / Not Ready`.
- List the items requiring human approval.

## Selected Task

- Task:
- Selection reason:
- Priority / ordinal:
- Dependencies:
- Dependency status:
- Assignee:

## Decision and Scope

- Decision approved: Yes / No
- Scope:
- Non-goals:
- Acceptance Criteria sufficient: Yes / No
- Major Change: Yes / No
- Human approval required:

## Model

- Model:
- Selection rationale:

## Readiness

- Ready / Not Ready:
- Blockers:
- Planned files:
