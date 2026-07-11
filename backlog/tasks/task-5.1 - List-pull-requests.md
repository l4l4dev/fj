---
id: TASK-5.1
title: List pull requests
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 19:35'
labels: []
milestone: m-4
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-5
ordinal: 38000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Discover pull requests in the selected repository.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Results show stable identity and lifecycle state
- [x] #2 Filtering and pagination are predictable
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: 現在のCodexモデル — Pull RequestのCLI契約・Application Port・Infrastructure境界・Presenter分離を扱うため、設計整合性を確認できるモデルを使用

Verification:
- gofmt -l . 成功
- git diff --check 成功
- go vet ./... 成功
- go test ./... 成功
- make pre-commit 成功
<!-- SECTION:NOTES:END -->

## Comments

<!-- COMMENTS:BEGIN -->
author: @codex
created: 2026-07-11 08:06
---
Independent Review

Critical:
なし

Major:
なし（前回のPresenter責務不整合はpullrequest_presenter.goへの分離で解消）

Minor:
今回は修正せず将来改善候補として記録:
- HTTP 401/403/404/その他remote failure境界テスト
- JSON decode failure
- path encoding
- secret redaction
- CLI delegation / Composition Root DI
- default flagsとPresenter単体テスト
- 既存Repository/Issue command regression

Suggestion:
将来改善候補として記録:
- PullRequest Presenter単体テストの追加
- Infrastructure error mapping fixtureの共通化
- state query mappingの統合テスト固定
---
<!-- COMMENTS:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
TASK-5.1の実装完了。fj pr list OWNER/NAME、承認済みpage/limit/state契約、PullRequest専用Application Port、Forgejo REST adapter、明示DI、Presenter分離を実装し、全検証と独立レビューを完了した。Minor/Suggestionは将来改善候補として記録した。
<!-- SECTION:FINAL_SUMMARY:END -->
