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

## Decision Classification

Classify every decision before deciding whether it requires human approval.

### Human Approval Required

Use this category for decisions that affect:

- Public CLI/API contracts, output, exit codes, or compatibility
- Architecture boundaries, dependency direction, or package structure
- Authentication, credentials, secret handling, or other security boundaries
- In-scope / out-of-scope behavior, task split, dependencies, or milestone scope
- New external dependencies
- Major Change status

Pending items in this category block implementation and require explicit human
approval.

### AI Decision Allowed

Use this category for non-public implementation choices that stay within the
approved scope and architecture, including:

- Private DTO names and internal field layout
- Helper functions and internal control flow
- Test fixtures and stubs
- Private formatting helpers that preserve the approved output
- Internal conversion and validation structure

The agent records the choice in the Implementation Plan and may proceed without
separate human approval.

### Convention Based

Use this category when an existing repository convention determines the choice,
including:

- Existing package placement and naming
- Existing `context.Context`, `apperror`, and DI patterns
- Existing Presenter boundaries
- Existing path encoding and credential redaction rules
- Standard verification commands
- Backlog CLI update procedures

Record the convention and its source. Escalate it if it conflicts with the
approved scope, public behavior, architecture, or security rules.

## Decision Escalation Rule

- Do not make assumptions for a Human Approval Required decision.
- If an AI Decision Allowed or Convention Based choice changes public behavior,
  architecture, security, compatibility, or scope, escalate it to Human
  Approval Required and stop implementation.
- If existing conventions conflict or are insufficient, classify the decision
  as Human Approval Required.
- Implementation may proceed only when all Human Approval Required items are
  approved, AI Decision Allowed items are recorded in the Implementation Plan,
  and Convention Based choices have documented sources.

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
- Options / proposal:
- Compatibility impact:
- Security impact:
- Scope impact:
- Major Change impact:

### AI Decision Allowed

- Internal decision:
- Proposed implementation:
- Constraint:

### Convention Based

- Decision:
- Convention:
- Source / evidence:

## Decision Status

- Ready for Approval
- Approved
- Rejected

## Rules

- Do not start implementation while the Decision is not approved.
- Do not decide a public contract by assumption.
- Explicitly list every item requiring human approval.
- Keep this plan separate from implementation and do not treat it as approval.
