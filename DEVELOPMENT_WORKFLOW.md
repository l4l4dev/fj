# Development Workflow

This workflow applies to human contributors and AI agents. It complements
`AGENTS.md`, `CONTRIBUTING.md`, and the Backlog task instructions.

## AI lifecycle

Every implementation task follows this lifecycle:

```text
Preflight → Decision review → Approval → Implementation → Verification → Independent review → Backlog finalization → Stop
```

Work on one task at a time. Commit and Push are separate, explicitly authorized
operations and are not part of task completion.

## Preflight

Before implementation:

1. Read `PROJECT_CONSTITUTION.md`, `ARCHITECTURE.md`, and `ROADMAP.md`.
2. Run `backlog instructions overview`.
3. Select an executable task whose dependencies are complete; choose the highest priority and use ordinal order as the tie-breaker.
4. Read the task, Acceptance Criteria, Implementation Notes, Decision, and dependencies.
5. Confirm Assignee, approved scope, non-goals, public compatibility, and Major Change status.
6. Stop for human approval when a public contract, security boundary, dependency direction, external dependency, or roadmap/milestone would change.

Record the selected task, selection reason, model and rationale, readiness decision,
and planned files in the task or work report.

### Privacy-safe records

Use placeholders instead of real personal names, organization names, hostnames,
or repository owner names in execution, verification, acceptance, and review
records by default. Even externally published information requires explicit
user approval or a clearly authorized publication purpose before being recorded.

Use:

- `example-owner`
- `example-repository`
- `https://forgejo.example.com`

Never record credential values, raw tokens, or credentials embedded in URLs.
Git history rewriting, force-push, and privacy cleanup require a separately
approved task and must not be performed as part of the current workflow.

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

During implementation, stop if the Acceptance Criteria are insufficient, an
unapproved scope change is needed, or the task will no longer remain a small,
reviewable unit. Do not start another task in the same loop.

## 5. Verification

Run the checks required by the task and project workflow. Record every command,
result, limitation, and unverified item in the Backlog task.

For Go changes, run at minimum:

```text
gofmt -l .
git diff --check
go vet ./...
go test ./...
make pre-commit
```

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

Use the Backlog CLI for all task updates. Before setting `Done`, confirm that
every Acceptance Criterion is checked, Verification and Independent Review are
recorded, Critical/Major findings are resolved, Implementation Notes and Final
Summary are present, Assignee is set, and no unapproved scope remains.

## 8. Commit / Push

Commit and Push are separate operations and require explicit authorization.
Never commit or push as an implicit part of task completion.
