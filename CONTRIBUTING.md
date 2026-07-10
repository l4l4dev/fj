# Contributing to fj

Thank you for helping improve `fj`. Contributions should be small, reviewable,
and consistent with the repository's governing documents.

## Before Starting

Read [PROJECT_CONSTITUTION.md](PROJECT_CONSTITUTION.md),
[ARCHITECTURE.md](ARCHITECTURE.md), [ROADMAP.md](ROADMAP.md), and
[AGENTS.md](AGENTS.md) before proposing or implementing a change. These
documents define the project's principles, architectural boundaries, direction,
and working rules. If they conflict with a task or implementation, stop and
raise the conflict for human resolution.

Use the existing Backlog task for the change. Keep work within its approved
scope and acceptance criteria. Do not combine unrelated cleanup, refactoring,
or feature work with the contribution. A scope expansion or a change requiring
an approval gate in `AGENTS.md` must be approved before implementation.

## Documentation

Update documentation whenever behavior changes. Describe important decisions,
their rationale, and their impact in the appropriate repository document or
Backlog record so the change does not depend on chat history or personal
knowledge. Do not modify `PROJECT_CONSTITUTION.md`, `ARCHITECTURE.md`, or
`ROADMAP.md` as an incidental part of another change.

## Testing and Verification

Add or update tests for every behavior change. Test the relevant success,
failure, and boundary cases at the narrowest effective architectural layer,
then run the broader applicable test suite before declaring the work complete.

Run the checks applicable to the contribution before requesting review:

1. **Formatting:** run `gofmt -l .` and require no output for Go sources. Run
   `git diff --check` to detect whitespace errors in all changed files.
2. **Static analysis:** run `go vet ./...` for all Go packages.
3. **Tests:** run focused tests while developing, then run `go test ./...` as
   the broader repository suite before completion.
4. **Documentation:** update user-facing and design documentation affected by
   the behavior, confirm changed Markdown links resolve to repository files,
   and verify that documented commands and behavior match the implementation.

Additional checks required by the task, affected subsystem, or review policy
remain mandatory. Record every command run and its result in the Backlog task
and completion report.

A failed required check blocks completion until the failure is fixed or a human
explicitly decides how to proceed. If a check is not applicable, is skipped, or
cannot run in the current environment, record the check, the reason, the scope
left unverified, and any alternative verification performed. Do not mark the
affected acceptance criteria complete or claim the unverified behavior is
complete.

## Compatibility

Treat public commands, flags, output formats, configuration, JSON contracts,
exit behavior, and exported APIs as compatibility-sensitive. Preserve existing
behavior unless an incompatible change has been explicitly approved. Any
approved incompatible change must document its user impact and provide
migration guidance.

## Review and Commit Scope

Before requesting review, confirm that the contribution satisfies its
acceptance criteria, follows the dependency direction in `ARCHITECTURE.md`, and
contains no unrelated files or generated artifacts. Follow the review and
approval requirements in `AGENTS.md`, including the additional checks required
for major changes.

Keep each commit focused on one purpose and use a Conventional Commit message
unless a reviewer provides an exact message. Do not commit or push on behalf of
another person without their explicit request.
