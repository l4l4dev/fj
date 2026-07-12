---
id: TASK-11.5
title: Release acceptance
status: Done
assignee: []
created_date: '2026-07-11 17:33'
updated_date: '2026-07-12 01:59'
labels: []
milestone: m-10
dependencies:
  - TASK-11.4
parent_task_id: TASK-11
priority: medium
ordinal: 11050
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define and execute acceptance for release artifacts, including install, version, help, read-only smoke checks, and secret redaction. TASK-10.5 results are inputs but not a required dependency.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Initial artifacts can be installed and verified.
- [x] #2 fj version and fj --help are confirmed after installation.
- [x] #3 Read-only smoke checks and secret redaction are evaluated.
- [x] #4 TASK-10.5 unresolved commands are not treated as successful.
- [x] #5 A live workflow_dispatch run verifies resolve-version, darwin/arm64 and linux/amd64 builds, SHA-256 checksum generation, and the consolidated workflow artifact.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Design Review (pre-acceptance): see doc-1 - TASK-11.5 Release Acceptance Design Review. Result: Critical 0 / Major 3 (LICENSE missing; live workflow run unverified; workflow_dispatch tag-commit binding) / Minor 6 / Suggestion 4. No fixes applied; findings await human decisions.
Model: Fable 5 — long-term design review of the Public Release Foundation (Section 15 review-grade work).

Human Decisions — Design Review Findings:
- LICENSE (Major): accepted as separate work. Tracked by TASK-12. No license has been selected; exact license selection requires a separate explicit human decision.
- GitHub Actions live execution (Major): accepted in TASK-11.5. A successful workflow_dispatch run must verify resolve-version, both build targets, checksum generation, and consolidated artifact generation.
- workflow_dispatch arbitrary version (Major): deferred for the release foundation. TASK-11.5 does not change the current input policy. A future public Release workflow must define a separate tag-only policy.
- Minor and Suggestion findings from doc-1 remain future improvement candidates and do not expand TASK-11.5 scope.

Release Acceptance Evidence Summary:

- Evidence source: user-provided execution records collected during TASK-11.5 release acceptance.
- AC #1: The release artifacts were verified through a macOS isolated installation on native macOS arm64 and a Linux amd64 runtime check.
- AC #2: The installed fj binary successfully reported its version and displayed command help.
- AC #3: The release artifact binary fj-0.1.0-artifacts/fj-0.1.0-darwin-arm64 successfully completed the Forgejo Playground read-only smoke commands: repository inspection, issue listing, and Issue #1 inspection. Failure-path verification confirmed secret redaction without recording credential values, endpoint details, authorization headers, or raw HTTP request/response data.
- AC #4: TASK-10.5 remains the source of truth for its unresolved repo list and pr list commands. Those commands remain unsuccessful, their acceptance criteria remain unchecked, and they are not included in TASK-11.5 successful release-acceptance evidence. TASK-11.5 read-only smoke acceptance covers only the successful repository inspection, issue listing, and issue inspection operations recorded above.
- AC #5: A live workflow_dispatch run verified resolve-version, the darwin/arm64 build, the linux/amd64 build, SHA-256 checksum generation, and consolidated workflow artifact generation.
- No fj implementation, README, workflow, LICENSE, or TASK-12 change was required for this acceptance task.

Independent Review (post-acceptance):
- Critical: none.
- Major: none.
- Minor: none.
- Suggestions: identify the evidence as user-provided execution records and explicitly preserve the TASK-10.5 unresolved-command boundary. Both suggestions are reflected in this evidence summary.
- Result: Ready for Finalization.
- The review confirmed that AC #1-#5 are supported, unsuccessful TASK-10.5 commands are not represented as successful, secret-redaction records contain no credential details, separate public-release policy and license decisions remain outside this task, and release acceptance is valid without implementation changes.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Completed release acceptance for the macOS arm64 and Linux amd64 artifacts. Verified isolated installation and runtime behavior, installed version and help output, Forgejo Playground read-only repository and issue operations, secret redaction, and the live workflow_dispatch build, checksum, and consolidated artifact pipeline. TASK-10.5 unresolved repo list and pr list commands remain explicitly unsuccessful and were not counted as release-acceptance successes. Independent review found no Critical, Major, or Minor findings and judged the task Ready for Finalization. No fj implementation changes were required.
<!-- SECTION:FINAL_SUMMARY:END -->
