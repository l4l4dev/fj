---
id: TASK-4.8
title: Issue metadata management
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-3
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-4
ordinal: 37000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Update labels, milestones, and assignments independently.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Each metadata category changes without overwriting the others
- [x] #2 Unknown or inaccessible values produce actionable errors
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
This parent task is split into TASK-4.8.1 Manage issue labels, TASK-4.8.2 Manage issue milestone, and TASK-4.8.3 Manage issue assignment because labels, milestones, and assignments have distinct CLI contracts, Forgejo API contracts, Application Ports, validation, error boundaries, and Presenter output. Each child is a 30-90 minute task and follows Decision approval → Implementation → Verification → Independent Review → Acceptance Criteria completion → Done.

Verification:
- 子Task TASK-4.8.1〜TASK-4.8.3の完了状態を確認済み
- 子TaskごとのVerification（gofmt -l .、git diff --check、go vet ./...、go test ./...、make pre-commit）を確認済み
- 親Taskのmetadata同期に伴うソースコード検証は適用外
<!-- SECTION:NOTES:END -->

## Comments

<!-- COMMENTS:BEGIN -->
author: @codex
created: 2026-07-11 07:28
---
Independent Review

Critical:
なし

Major:
なし

Minor:
子Taskごとの境界テスト拡充余地は各SubtaskのReviewへ記録済み

Suggestion:
M4開始前にM3全体の回帰検証を継続検討
---
<!-- COMMENTS:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
TASK-4.8のFinalization完了。TASK-4.8.1 labels、TASK-4.8.2 milestone、TASK-4.8.3 assignmentの完了内容を親Taskへ同期し、Issue metadata managementのM3スコープを完了状態へ集約した。
<!-- SECTION:FINAL_SUMMARY:END -->
