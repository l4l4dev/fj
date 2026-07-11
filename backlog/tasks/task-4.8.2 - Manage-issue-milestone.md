---
id: TASK-4.8.2
title: Manage issue milestone
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 06:11'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-3
dependencies:
  - TASK-2.9
parent_task_id: TASK-4.8
ordinal: 87000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Manage issue milestone independently as a focused Issue metadata workflow.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The CLI contract for setting and clearing an issue milestone is explicitly defined and supports --instance PROFILE.
- [x] #2 OWNER/NAME and positive issue number validation are enforced before remote access.
- [x] #3 Application defines a milestone model and dedicated Port without changing existing issue or repository Ports.
- [x] #4 Infrastructure uses the approved Forgejo milestone API contract and translates failures through the Application error boundary.
- [x] #5 Milestone changes do not overwrite labels, assignment, title, body, or state.
- [x] #6 Human-readable milestone output is defined and does not change existing issue command output.
- [x] #7 Tests cover validation, API mapping, error boundary, Presenter output, and explicit Composition Root injection.
- [x] #8 The CLI provides fj issue milestone set OWNER/NAME NUMBER MILESTONE and fj issue milestone clear OWNER/NAME NUMBER with --instance; milestone is name-only, non-empty, non-whitespace, and issue number is a positive integer.
- [x] #9 Setting replaces an existing milestone and is idempotent for the same milestone; clearing removes the milestone and is idempotent when already clear.
- [x] #10 Application owns Milestone{ID int64, Title string}, SetMilestoneRequest, ClearMilestoneRequest, MilestoneSetter, and MilestoneClearer in internal/application/issue; Repository ports remain unchanged.
- [x] #11 Infrastructure resolves milestone names internally through GET /api/v1/repos/{owner}/{repo}/milestones using exact name matching; no user-facing milestone listing command is added.
- [x] #12 Set sends PATCH /api/v1/repos/{owner}/{repo}/issues/{index} with {"milestone":<id>}; clear sends the same endpoint with {"milestone":null}; paths are safely encoded and DTOs remain private.
- [x] #13 The operations are set issue milestone and clear issue milestone; failures use the existing Application error boundary without exposing HTTP status, credentials, URL details, response bodies, or raw causes to CLI.
- [x] #14 Application remains independent of HTTP, Cobra, and Forgejo DTOs; Interface → Application Milestone Port → Infrastructure adapter, explicit DI, unchanged Repository ports, and no runtime type assertions are preserved.
- [x] #15 Set output is Issue: #<number> followed by Milestone set: <title>; clear output is Issue: #<number> followed by Milestone cleared, and existing issue command output remains unchanged.
- [x] #16 Milestone resolution keeps only private DTO id/title fields and converts them to Application Milestone{ID, Title}; set PATCH response milestone is converted and returned, while clear returns error only.
- [x] #17 Operations are set issue milestone and clear issue milestone; validation, authentication, and remote categories are used. Milestone/repository/issue not found are remote, 401/403 are authentication, and HTTP status is not exposed to CLI.
- [x] #18 Same milestone set and already-clear operations skip PATCH and succeed idempotently.
- [x] #19 Milestone names are matched exactly with case sensitivity.
- [x] #20 Infrastructure may internally GET the current Issue at GET /api/v1/repos/{owner}/{repo}/issues/{index} solely for idempotency; no new user-facing inspect command is added.
- [x] #21 Existing Application IssueDetail remains unchanged; private Infrastructure DTOs represent Issue number and nullable milestone with id/title fields.
- [x] #22 Set resolves milestones, matches name exactly, gets the current Issue, compares milestone IDs, skips PATCH on the same ID, and otherwise patches the selected ID.
- [x] #23 Clear gets the current Issue, skips PATCH when milestone is nil, and patches null when a milestone exists.
- [x] #24 Set response milestone object is converted to Application Milestone; clear accepts a null milestone response and returns error only; comparison uses IDs, not names.
- [x] #25 Error operations are resolve milestone, get issue milestone, set issue milestone, and clear issue milestone. Not found is remote, 401/403 authentication, and other HTTP/JSON failures remote.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved design:
- CLI: `fj issue milestone set OWNER/NAME NUMBER MILESTONE` and `fj issue milestone clear OWNER/NAME NUMBER`, both with `--instance PROFILE`. Milestone is specified by name only; ID input is not supported. Empty or whitespace-only milestone is invalid; issue number is positive.
- Semantics: set replaces an existing milestone and is idempotent for the same milestone. clear removes the existing milestone and is idempotent when already clear.
- Application: Milestone{ID int64, Title string}, SetMilestoneRequest, ClearMilestoneRequest, MilestoneSetter, and MilestoneClearer in internal/application/issue. Repository ports remain unchanged.
- Infrastructure: resolve milestone names internally using GET /api/v1/repos/{owner}/{repo}/milestones with exact name matching; no user-facing milestone listing command. Set uses PATCH /api/v1/repos/{owner}/{repo}/issues/{index} with {"milestone":<id>}; clear uses the same endpoint with {"milestone":null}. Safely encode paths, keep DTOs private, convert to Application models, and use the existing apperror boundary without exposing HTTP status to CLI.
- Operations: `set issue milestone` and `clear issue milestone`.

Additional approved decisions:
- Presenter: set outputs `Issue: #<number>` and `Milestone set: <title>`; clear outputs `Issue: #<number>` and `Milestone cleared`. Existing issue command output remains unchanged.
- Milestone resolution GET keeps only private DTO fields id and title, then converts to Application Milestone{ID, Title}.
- Issue PATCH response uses a private DTO with milestone. Set converts the response milestone to Application Milestone and returns it. Clear returns error only; no Milestone result is required.
- Operations are `set issue milestone` and `clear issue milestone`. Error categories are validation, authentication, and remote. Milestone not found and repository/issue not found are remote; 401/403 are authentication; HTTP status is not exposed to CLI.
- Idempotency: same milestone set and already-clear both skip PATCH and succeed.
- Milestone names use exact, case-sensitive matching.

Additional approved decisions:
- Internal current Issue GET is allowed only for idempotency: GET /api/v1/repos/{owner}/{repo}/issues/{index}. It is not a new user-facing inspect command and is used to inspect the current milestone.
- Existing Application IssueDetail is unchanged. Infrastructure uses private IssueDTO{number, milestone *MilestoneDTO} and private MilestoneDTO{id, title}.
- Set order: resolve milestone list, exact case-sensitive name-to-ID match, get current Issue, compare current milestone ID. Same ID skips PATCH and succeeds.
- Clear gets current Issue; nil milestone skips PATCH and succeeds; an existing milestone sends PATCH with null.
- Milestone comparison uses ID equality, never name comparison.
- Set PATCH response contains milestone object {id, title}, which converts to Application Milestone. Clear permits a null milestone response and returns error only.
- Error operations are `resolve milestone`, `get issue milestone`, `set issue milestone`, and `clear issue milestone`. Not found is remote, 401/403 authentication, and other HTTP/JSON failures remote.

Verification:
- gofmt -l . 成功
- git diff --check 成功
- go vet ./... 成功
- go test ./... 成功
- make pre-commit 成功

Independent Review:

### Critical
なし

### Major
なし

### Minor
境界テスト拡充余地:
- error boundary
- path encoding
- secret redaction
- DI／委譲
- regression tests

### Suggestion
- same milestone idempotency test追加検討
- clear存在ケース追加検討
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
TASK-4.8.2の実装完了: Application-owned Milestone Port、milestone name-to-ID resolution、current IssueによるID比較、idempotentなset/clear、Forgejo PATCH、safe error boundary、明示DI、固定Presenter出力を実装した。
<!-- SECTION:FINAL_SUMMARY:END -->
