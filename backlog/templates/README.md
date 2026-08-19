# Backlog Workflow Reference

`AGENTS.md` and `DEVELOPMENT_WORKFLOW.md` are authoritative. This file is a
quick reference for recording task work through the `backlog` CLI.

## Routine task loop

```text
Select task → Plan approval → Implementation + tests → Verification → Report → Stop
```

Record in the task, via the `backlog` CLI only:

- **Plan approval** — one note line, e.g.
  `backlog task edit TASK-123 --append-notes "Plan approved by human on YYYY-MM-DD."`
- **Completion** — check verified Acceptance Criteria (`--check-ac`), append a
  short implementation note (1–2 lines), and set a one-line
  `--final-summary "Changed X, verified with Y."` before `-s Done`.

Use privacy placeholders (`example-owner`, `example-repository`,
`https://forgejo.example.com`) in all records. Never record credentials.

## Major changes only

For changes classified as major by `AGENTS.md` Section 15, additionally
record:

- the pre-implementation check result (one note line: scope, risks, verdict);
- the independent review result (one note line: findings by severity, verdict).

Do not treat chat transcripts as durable evidence; the Backlog task and
repository documents are the record.
