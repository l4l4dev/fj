---
id: TASK-5.2
title: Inspect a pull request
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 08:48'
labels: []
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-5
ordinal: 39000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Display pull request metadata, branches, and current state.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Source, target, repository, and pull request identity are explicit
- [x] #2 Missing and inaccessible pull requests produce distinct errors
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision approved:
- CLI: fj pr inspect OWNER/NAME NUMBER
- PullRequestDetail: Number, Title, State, HeadBranch, BaseBranch, Body
- Output: fixed human-readable format
- Out of scope: review, comments, merge, labels, JSON
- AI Decision Allowed: private DTO, helper structure, test fixtures, internal implementation
- Convention Based: existing CLI pattern, Presenter boundary, Error Boundary, DI pattern
- Implementation gate: after Decision Approved, proceed to implementation-plan.md

Model: 現在のCodexモデル — Pull Request inspectのApplication境界、REST API、CLI互換性を確認するため、設計整合性を扱えるモデルを使用

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
created: 2026-07-11 08:48
---
Independent Review

Critical:
なし

Major:
なし

Minor:
今回は修正せず将来改善候補として記録:
- HTTP error boundary test
- JSON decode failure test
- path encoding test
- secret redaction test
- CLI delegation / Composition Root DI test
- empty field presenter test

Suggestion:
将来改善候補として記録:
- PullRequestDetail Presenter単体テストの拡充
- 既存pr listとinspectの回帰テスト共通化
---
<!-- COMMENTS:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
TASK-5.2の実装完了。fj pr inspect OWNER/NAME NUMBER、PullRequestDetail、Inspect専用Port、Forgejo REST adapter、Error Boundary、Presenter、明示DIを実装し、VerificationとIndependent Reviewを完了した。
<!-- SECTION:FINAL_SUMMARY:END -->
