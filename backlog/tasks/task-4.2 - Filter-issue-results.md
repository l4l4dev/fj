---
id: TASK-4.2
title: Filter issue results
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 11:55'
updated_date: '2026-07-11 06:06'
labels: []
dependencies:
  - TASK-2.9
references:
  - ROADMAP.md
parent_task_id: TASK-4
ordinal: 31000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Goal: Narrow issue discovery through documented filters.

Intended scope: approximately 30-90 minutes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Supported filter combinations are deterministic
- [x] #2 Invalid filter values are rejected locally
- [x] #3 The CLI accepts a single --assignee USER or --label LABEL filter in addition to the existing --page, --limit, and --state open|closed|all flags.
- [x] #4 Specifying both --assignee and --label, or repeating either filter, is rejected as a local validation error.
- [x] #5 Application owns an IssueFilter value and does not expose Forgejo query parameter names.
- [x] #6 Infrastructure maps IssueFilter.Assignee to assignee and IssueFilter.Label to labels in the Forgejo request.
- [x] #7 Existing fj issue list OWNER/NAME behavior, output, exit codes, pagination, and state defaults remain compatible with TASK-4.1; filter values are not displayed.
- [x] #8 Author, milestone, sort, keyword search, fetch-all, and JSON output remain out of scope.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add the Application-owned IssueFilter to the existing issue list request. 2. Add single-value assignee and label CLI flags with local validation while preserving TASK-4.1 behavior. 3. Map filters to Forgejo query parameters in the REST adapter. 4. Add regression and boundary tests, then run the full verification suite.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Approved filter contract: add single-value --assignee USER and --label LABEL filters while preserving --page, --limit, and --state open|closed|all. Only one filter may be specified per invocation; repeated assignee or label values are out of scope. Filter values are not shown in human-readable output. Existing TASK-4.1 command, output, and exit-code compatibility must remain unchanged. Application adds an IssueFilter value owned by the issue package; Forgejo query names do not cross the Application boundary. Infrastructure maps assignee to the assignee query and label to the labels query. Author, milestone, sort, keyword search, fetch-all, and JSON output are out of scope.

Implemented IssueFilter with single-value --assignee and --label flags. Added local validation for empty, combined, and repeated filters; mapped assignee to assignee and label to labels in the Forgejo request; preserved TASK-4.1 output, pagination, state, exit-code, error-boundary, Presenter, and explicit DI behavior. Validation passed: gofmt -l ., git diff --check, go vet ./..., go test ./..., and make pre-commit (GOCACHE=/tmp/fj-gocache for sandbox compatibility).

Historical note: This task was completed before the standard workflow was introduced. No Independent Review record exists from that period.

Verification:
Historical note:
- This was a task completed before workflow standardization; execution records were not preserved in the current format.

Independent Review:
Historical note:
- This was a task completed before workflow standardization; no Independent Review record was preserved in the current format.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Completed TASK-4.2 filter support with Application-owned IssueFilter, assignee/label CLI flags, deterministic single-filter validation, and Forgejo query mapping. Existing TASK-4.1 behavior remains compatible and filter values are not rendered.
<!-- SECTION:FINAL_SUMMARY:END -->
