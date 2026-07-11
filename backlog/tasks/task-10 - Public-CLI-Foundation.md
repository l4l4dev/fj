---
id: TASK-10
title: Public CLI Foundation
status: To Do
assignee: []
created_date: '2026-07-11 09:08'
updated_date: '2026-07-11 17:15'
labels: []
dependencies: []
priority: low
ordinal: 10990
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Make fj usable by the developer on macOS with a local install, Forgejo Playground onboarding, early read-only user acceptance, version metadata, and a validated quickstart. Homebrew, release automation, and large-scale distribution are out of scope.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 A developer can install fj for the current macOS user and run it from PATH.
- [ ] #2 A developer can configure and connect to Forgejo Playground without exposing credentials.
- [ ] #3 Read-only user acceptance succeeds for repository, issue, and pull request commands.
- [ ] #4 Version and Quickstart guidance match the validated local workflow.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision — Parent Finalization:
- Option A (Strict Completion) adopted.
- TASK-10 cannot be marked Done while TASK-10.5 remains In Progress.
- TASK-10.5 acceptance exception remains an explicit incomplete record; its unmet Acceptance Criteria are not marked passed.
- TASK-10.5 is currently held by external Forgejo Playground permission/API constraints, not a confirmed fj code defect.
- Next milestones must not be blocked solely by TASK-10 completion; independently executable follow-up milestones may proceed while TASK-10.5 awaits re-acceptance.
<!-- SECTION:NOTES:END -->
