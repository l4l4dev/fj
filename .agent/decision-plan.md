# Decision Plan

## Purpose

Use this plan to identify unresolved design decisions for an otherwise
executable Leaf Task and present them for human approval before implementation.
The plan is a design handoff, not an implementation authorization.

## When to use

Use this plan when any of the following applies:

- The task is a Leaf Task but its Decision is not approved.
- The CLI contract is not final.
- The Application Port is not final.
- The Infrastructure API contract is not final.
- Presenter or Error Boundary behavior is not final.
- A Major Change determination is required.

## Template

## Task

- Task ID:
- Title:
- Parent:
- Dependencies:
- Current task type:

## Current Status

- Backlog Status:
- Assignee:
- Decision status:
- Implementation started: Yes / No
- Ready for implementation: Yes / No

## Missing Decisions

### CLI Contract

- Command:
- Arguments:
- Flags:
- Validation:
- Output:

### Application Design

- Model:
- Request:
- Port:
- Dependency direction:

### Infrastructure Contract

- API endpoint:
- Request body:
- Response DTO:
- Error mapping:

### Behavior Semantics

- Idempotency:
- Edge cases:
- Invalid input:

### Presenter

- Success output:
- Error output:

### Scope

- In scope:
- Out of scope:

### Human Approval Required

- Decision requiring approval:
- Compatibility impact:
- Security impact:
- Major Change impact:

## Decision Status

- Ready for Approval
- Approved
- Rejected

## Rules

- Do not start implementation while the Decision is not approved.
- Do not decide a public contract by assumption.
- Explicitly list every item requiring human approval.
- Keep this plan separate from implementation and do not treat it as approval.
