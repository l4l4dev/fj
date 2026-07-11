# Implementation Plan Template

## Task

- Task ID:
- Goal:
- Scope:
- Non-goals:
- Dependencies:

## Boundaries

- CLI contract:
- Application model / Port:
- Infrastructure adapter / API:
- Error Boundary:
- Presenter:
- Composition Root / DI:

## Changes

| File | Responsibility | Change |
| --- | --- | --- |
| | | |

## Tests

- Application:
- Infrastructure:
- CLI / Presenter:
- Regression:

## Verification

- [ ] gofmt -l .
- [ ] git diff --check
- [ ] go vet ./...
- [ ] go test ./...
- [ ] make pre-commit

## Workflow Handoff

After the Implementation Plan is complete:

1. Implement only the approved scope.
2. Execute the applicable Verification checks above.
3. Proceed to `.agent/independent-review.md` after Verification completes.
4. Proceed to `.agent/backlog-finalization.md` only after Independent Review completes.

### Independent Review Rules

- Classify findings as Critical, Major, Minor, or Suggestion.
- Critical findings require immediate stop and human decision.
- Major findings mean the implementation is not complete and require human decision.
- For Minor findings, AI presents a fix proposal or deferral proposal; adoption, deferral, or rejection requires human decision before any fix.
- Suggestions are recorded as future improvement candidates; adoption, deferral, or rejection requires human decision.
- Do not proceed to Backlog Finalization while Critical or Major findings remain unresolved or the required human decision is missing.

### Backlog Finalization Rules

Backlog Finalization is allowed only after Independent Review. Confirm all of
the following before synchronizing completion metadata:

- Acceptance Criteria
- Verification record
- Review record
- Final Summary
- Status change conditions

## Role Handoff

- `implementation-plan.md`: translate the approved specification into an actionable implementation plan.
- `independent-review.md`: evaluate the implementation result and classify findings.
- `backlog-finalization.md`: synchronize the task into its completed state.
