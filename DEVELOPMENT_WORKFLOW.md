# Development Workflow

This workflow applies to human contributors and AI agents. It complements
`AGENTS.md`, `CONTRIBUTING.md`, and the Backlog task instructions.

## 1. Task selection

Select one executable task from Backlog according to status, dependencies,
priority, and approved scope. Confirm the task definition and relevant project
documents before starting.

## 2. Decision review

Review the task's Acceptance Criteria, Implementation Notes, dependencies, and
the applicable architecture and security boundaries. Identify unresolved
public-contract or compatibility decisions before implementation.

## 3. Approval

Obtain human approval for the implementation Decision whenever the task has an
approval gate. Record the approved decision in the Backlog task before coding.

## 4. Implementation

Assign an implementation owner, update the task to In Progress, and change only
the approved scope. Add or update tests for behavior changes.

## 5. Verification

Run the checks required by the task and project workflow. Record every command,
result, limitation, and unverified item in the Backlog task.

## 6. Independent review

After implementation and verification, obtain an independent review when
required by `AGENTS.md`. Classify findings as Critical, Major, Minor, or
Suggestion and resolve or explicitly defer them before completion.

## 7. Backlog finalization

Before marking a task Done:

- confirm the implementation is complete;
- check every Acceptance Criterion;
- update Implementation Notes and the review summary;
- record the Final Summary;
- confirm no unapproved scope remains;
- set Status to Done only after all completion conditions are satisfied.

## 8. Commit / Push

Commit and Push are separate operations and require explicit authorization.
Never commit or push as an implicit part of task completion.
