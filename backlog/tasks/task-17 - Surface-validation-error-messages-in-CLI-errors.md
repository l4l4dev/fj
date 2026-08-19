---
id: TASK-17
title: Surface validation error messages in CLI errors
status: Done
assignee:
  - '@claude'
created_date: '2026-08-19 07:59'
updated_date: '2026-08-19 07:59'
labels: []
dependencies:
  - TASK-5.6
type: bug
ordinal: 96000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
mapApplicationError dropped the message of apperror.ValidationError, so validation failures (e.g. an unsupported fj pr review --outcome) printed only 'invalid input'. Surface the specific message CLI-wide, keeping 'invalid input' as the fallback when no message exists. Found as a Minor finding in the TASK-5.6 independent review.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Validation errors with a message print that message; messageless ones still print 'invalid input'
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Plan approved by human in chat on 2026-08-19. One-line fix in mapApplicationError plus a unit test.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
mapApplicationError now passes ValidationError.Message through newCommandErrorWithMessage. Verified with the new TestMapApplicationErrorSurfacesValidationMessage and make pre-commit.
<!-- SECTION:FINAL_SUMMARY:END -->
