# Standard Task Workflow

This workflow supplements `AGENTS.md`; it never grants authority that the task, a human, or project rules have not granted.

```text
Investigation
    ↓
Human Decision (when required)
    ↓
Implementation
    ↓
Independent Review
    ├─ Ready → Finalization
    └─ Changes Required
           ↓
       Remediation
           ↓
       Independent Re-Review
           ├─ Ready → Finalization
           └─ Changes Required → Remediation
    ↓
Commit (only after explicit authorization)
```

## Shared Controls

- The repository is the Source of Truth; chat history is not durable Evidence.
- State the task, scope, allowed files, and prohibited changes before mutation.
- Do not mix Deferred Items or opportunistic improvements into the current completion scope.
- Keep implementer and independent reviewer roles distinct.
- Unless explicitly authorized, do not commit, push, tag, release, run external workflows, check ACs, add a Final Summary, or set Done.
- Classify findings as Critical, Major, Minor, or Suggestion according to `AGENTS.md`. Critical and Major findings block Finalization. Minor findings require a human adopt, defer, or reject decision before a fix; Suggestions are recorded when relevant.
- Finalization requires no Critical or Major findings, all ACs passing, and no unresolved Human Decision.

## Investigation

- **Entry Conditions:** A task or problem is identified; facts, options, or risks are not yet sufficiently established.
- **Activities:** Read governing records and relevant files; gather Evidence; compare options; identify risks and decision needs.
- **Outputs:** Investigation record, options, Evidence, risks, open questions, and proposed next step.
- **Prohibited Actions:** Implementation, irreversible or external changes, and automatic material decisions.
- **Exit Conditions:** The problem is understood well enough to request a Decision or plan implementation.
- **Next Step:** Human Decision when required; otherwise Implementation planning.

## Human Decision

- **Entry Conditions:** Investigation identifies a material choice or approval gate.
- **Activities:** Present options, trade-offs, recommendation, and consequences; obtain an explicit human choice; record authorization and prohibitions through Backlog.
- **Outputs:** Durable Human Decision and approved implementation boundary.
- **Prohibited Actions:** Treating an AI recommendation as approval or implementing before the Decision is recorded.
- **Exit Conditions:** The exact Decision, rationale, allowed scope, and implementation authorization are recorded.
- **Next Step:** Implementation.

## Implementation

- **Entry Conditions:** Dependencies are complete; required Decision is recorded; scope and owner are clear.
- **Activities:** Implement the smallest approved change, verify it proportionately, and record changes and limitations.
- **Outputs:** Reviewable diff, verification results, implementation record, and proposed commit message.
- **Prohibited Actions:** Scope expansion, self-declared independent review, AC checks, Final Summary, Done, commit, or push without authorization.
- **Exit Conditions:** Approved implementation and verification are complete with no hidden limitations.
- **Next Step:** Independent Review.

## Independent Review

- **Entry Conditions:** Implementation and verification are complete; a reviewer independent from implementation is available.
- **Activities:** Read the diff and durable Evidence; assess every AC; classify findings; evaluate scope and Deferred Items.
- **Outputs:** Pass, Fail, or Insufficient Evidence per AC; findings; Ready or Not Ready decision; required remediation or proposed Finalization record.
- **Prohibited Actions:** File, Backlog, status, or AC changes during review.
- **Exit Conditions:** Review decision and findings are reported to a human.
- **Next Step:** Finalization when Ready; otherwise human disposition and Remediation.

## Remediation

- **Entry Conditions:** A human has selected findings to adopt or a Critical finding requires an approved corrective response.
- **Activities:** Fix only the selected findings; verify the fix and affected scope; document remaining findings.
- **Outputs:** Focused remediation diff, verification, and re-review request.
- **Prohibited Actions:** Opportunistic improvements, unrelated refactors, Finalization, AC checks, Done, commit, or push.
- **Exit Conditions:** Approved remediation is complete and reviewable.
- **Next Step:** Independent Re-Review.

## Independent Re-Review

- **Entry Conditions:** Remediation is complete and an independent reviewer is available.
- **Activities:** Classify each original finding as Resolved, Partially Resolved, or Not Resolved; reassess ACs and regressions.
- **Outputs:** Finding resolution, new findings, and Ready or Not Ready decision.
- **Prohibited Actions:** Self-approval by the implementer or any mutation during review.
- **Exit Conditions:** No unresolved blocking finding, or required remediation is clearly stated.
- **Next Step:** Finalization when Ready; otherwise human disposition and Remediation.

## Finalization

- **Entry Conditions:** Latest independent review says Ready; all ACs pass; no Critical, Major, or unresolved Human Decision remains.
- **Activities:** Persist Decision and Evidence summaries, review result, resolved findings, Deferred Items, Final Summary, AC checks, and allowed status transition.
- **Outputs:** Internally consistent completed Backlog task and final Git-state report.
- **Prohibited Actions:** New implementation, unrelated cleanup, commit, push, tag, or release.
- **Exit Conditions:** Completion checklist in `AGENTS.md` is satisfied.
- **Next Step:** Commit only after explicit human authorization.

## Commit

- **Entry Conditions:** Finalization is complete and a human explicitly authorizes the commit.
- **Activities:** Stage only approved files, inspect the staged diff, commit with the approved message, and report SHA and status.
- **Outputs:** One focused commit.
- **Prohibited Actions:** Unapproved files, push, tag, or release unless separately authorized.
- **Exit Conditions:** Commit scope and repository state are verified.
- **Next Step:** Push confirmation when separately requested.
