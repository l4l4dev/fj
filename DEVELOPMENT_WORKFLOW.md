# Development Workflow

This workflow applies to human contributors and AI agents. It complements
`AGENTS.md`, `CONTRIBUTING.md`, and the Backlog task instructions.

## Task lifecycle

Every implementation task follows this single loop:

```text
Select task → Plan approval → Implementation + tests → Verification → Report → Stop
```

Work on one task at a time. Commit and Push are separate, explicitly
authorized operations and are not part of task completion.

Major changes (`AGENTS.md` Section 15) additionally require the
pre-implementation check before coding and the independent review before
finalization. Routine tasks — including new commands and flags that follow
existing patterns — do not.

## 1. Task selection

Select one executable task from Backlog according to status, dependencies,
and priority. Read the task, its Acceptance Criteria, and the code it touches.

## 2. Plan approval

Present a short implementation plan and wait for human approval. A plan-mode
approval or an explicit chat approval is sufficient. Record the approval in
the task in one line, set the Assignee, and move the task to In Progress.

Stop for human approval when the change is major per `AGENTS.md` Section 15.

## 3. Implementation

Change only the approved scope. Add or update tests for behavior changes.
Stop if the Acceptance Criteria are insufficient, an unapproved scope change
is needed, or the task will no longer remain a small, reviewable unit.

## 4. Verification

For Go changes, run `make pre-commit` (gofmt, `git diff --check`,
`go vet ./...`, `go test ./...`) before reporting completion.

## 5. Finalization and report

Before marking a task Done: check every Acceptance Criterion that verified
evidence proves, append a short implementation note (1–2 lines), and record a
one-line Final Summary. Then report to the human per `AGENTS.md` Section 16
and stop. Never commit or push without explicit authorization.

## Privacy-safe records

Use placeholders instead of real personal names, organization names,
hostnames, or repository owner names in all records by default:

- `example-owner`
- `example-repository`
- `https://forgejo.example.com`

Never record credential values, raw tokens, or credentials embedded in URLs.
Git history rewriting, force-push, and privacy cleanup require a separately
approved task.
