---
id: TASK-18
title: Add fj auth login for token registration
status: To Do
assignee: []
created_date: '2026-08-19 08:37'
labels: []
dependencies: []
type: feature
ordinal: 97000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide a first-run route to register a token from the command line (option C, minimal gh-style login). fj auth login --instance <name> reads a token (hidden input on a TTY, --token-stdin for pipes), verifies it against GET /api/v1/user, and stores it in a dedicated credentials file (XDG path, 0600 permissions) separate from the main config. Credential resolution order becomes: environment variable (existing, unchanged for CI/agents) > credentials file. auth status / auth logout are follow-up tasks, not in scope. This is a major change (credential handling + new dependency golang.org/x/term, BSD-3): run the Section 15 pre-implementation check before coding and the independent review after.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 fj auth login verifies the token against the instance before storing it and rejects invalid or unauthorized tokens with actionable errors
- [ ] #2 The stored credentials file is separate from the main config and created with 0600 permissions
- [ ] #3 Environment-variable credentials keep precedence over the stored file, so existing setups and CI behavior are unchanged
- [ ] #4 Token input is hidden on a TTY and available non-interactively via --token-stdin; the token value never appears in output, logs, or errors
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision: Mika approved option C (login only; status/logout later) in chat on 2026-08-19, including the golang.org/x/term dependency for hidden input. License compliance per TASK-15 applies (BSD-3-Clause).
<!-- SECTION:NOTES:END -->
