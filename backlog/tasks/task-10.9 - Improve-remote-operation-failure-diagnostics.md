---
id: TASK-10.9
title: Improve remote operation failure diagnostics
status: Done
assignee:
  - '@codex'
created_date: '2026-07-11 11:44'
updated_date: '2026-07-11 11:44'
labels: []
dependencies: []
parent_task_id: TASK-10
priority: medium
ordinal: 89010
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Improve safe diagnostics for remote operation failures without changing the public CLI contract.

Scope:
- Infrastructure internal typed error metadata
- Application semantic error category conversion
- Pull Request 404 NotFound mapping
- Secret redaction preservation
- Boundary tests

Non-goals:
- CLI contract changes
- API specification changes
- JSON output
- Large-scale error architecture refactor
- Unrelated refactoring

Implementation status: completed before formal Task creation.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Infrastructure retains HTTP status only as internal typed error metadata.
- [x] #2 Remote failures are converted to Application semantic categories without exposing HTTP status, raw response, or transport cause to CLI output.
- [x] #3 Pull Request list 404 is classified as Application NotFound using the existing safe repository-not-found message.
- [x] #4 Credential, token, URL credential, and response body redaction is preserved.
- [x] #5 Boundary tests cover Pull Request list NotFound mapping and existing empty-result behavior.
<!-- AC:END -->



## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: Codex — targeted Infrastructure error mapping and boundary test implementation.

Verification:
- gofmt -l . passed
- git diff --check passed
- go vet ./... passed
- go test ./... passed
- make pre-commit passed.

Independent Review:
- Critical: none
- Major: none
- Minor: Pull Request list 401/403/other-remote boundary tests remain future expansion candidates.
- Suggestion: consider a shared Repository/Pull Request error-mapping helper in a future task.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Formalized the completed remote failure diagnostics improvement: internal status metadata remains private, Application semantic categories are preserved, Pull Request list 404 maps to NotFound with the existing safe message, and boundary coverage was added. No public CLI contract changed.
<!-- SECTION:FINAL_SUMMARY:END -->
