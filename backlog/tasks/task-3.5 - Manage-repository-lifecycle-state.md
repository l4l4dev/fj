---
id: TASK-3.5
title: Manage repository lifecycle state
status: Done
assignee: []
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 00:07'
labels: []
dependencies:
  - TASK-2.9
  - TASK-2.10
  - TASK-2.11
  - TASK-2.12
  - TASK-2.13
  - TASK-2.14
references:
  - ROADMAP.md
parent_task_id: TASK-3
ordinal: 28000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Support approved archive, restore, and deletion workflows safely.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Consequential operations require explicit target and intent
- [x] #2 The CLI provides fj repo archive OWNER/NAME and fj repo restore OWNER/NAME with optional --instance PROFILE, using existing target validation and instance selection.
- [x] #3 Application adds a dedicated Archiver port without changing List, Getter, Creator, or Updater contracts; archive and restore request PATCH /api/v1/repos/{owner}/{repo} with archived true or false.
- [x] #4 The REST adapter safely encodes owner/name, sends only {archived:true|false}, and asserts the Archiver interface.
- [x] #5 Success output is exactly Repository: OWNER/NAME followed by Archived: true or false, with no delete, prompt, dry-run, or JSON behavior.
- [x] #6 Invalid input is validation; 401/403 authentication; 404 repository not found; other HTTP, transport, and JSON failures remote; credentials, sensitive URLs, response bodies, and raw causes are not exposed.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — architecture-sensitive lifecycle change with a new public CLI, dedicated port, PATCH semantics, and safety-sensitive error classification.

Validation: make pre-commit passed, including git diff --check, go vet ./..., and go test ./....

Independent post-implementation review (GPT-5): initial Major candidate on archive input classification was resolved with ValidationError; re-review Critical none, Major none. Minor: dedicated REST archive/restore failure and redaction tests could be expanded later.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented archive and restore lifecycle commands only. Added a dedicated Archiver port and use case, PATCH archived true/false adapter, fj repo archive/restore with existing target and instance rules, fixed Repository/Archived output, safe authentication/remote/internal classification, and regression tests. Repository deletion, prompts, dry-run, and JSON remain out of scope.
<!-- SECTION:FINAL_SUMMARY:END -->
