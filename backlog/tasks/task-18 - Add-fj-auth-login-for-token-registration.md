---
id: TASK-18
title: Add fj auth login for token registration
status: Done
assignee:
  - '@claude-opus'
created_date: '2026-08-19 08:37'
updated_date: '2026-08-19 10:33'
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
- [x] #1 fj auth login verifies the token against the instance before storing it and rejects invalid or unauthorized tokens with actionable errors
- [x] #2 The stored credentials file is separate from the main config and created with 0600 permissions
- [x] #3 Environment-variable credentials keep precedence over the stored file, so existing setups and CI behavior are unchanged
- [x] #4 Token input is hidden on a TTY and available non-interactively via --token-stdin; the token value never appears in output, logs, or errors
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision: Mika approved option C (login only; status/logout later) in chat on 2026-08-19, including the golang.org/x/term dependency for hidden input. License compliance per TASK-15 applies (BSD-3-Clause).

Pre-implementation check (Fable subagent, 2026-08-19): PROCEED. Mandatory correction: relax config.Validate to allow an empty credential reference (currently blocks the feature; backward compatible). Minors incorporated: atomic 0600 temp+rename write, missing credentials file maps to ErrCredentialUnavailable, trim and reject blank stdin tokens. Suggestions adopted: resolve-credential failure now hints at fj auth login; success output prints the storage path. Constitution/architecture/roadmap consistent; x/term+x/sys BSD-3 notices required per TASK-15. Implementation delegated to an Opus subagent.

Independent review (Fable subagent, independent of the Opus implementer, 2026-08-19): Review Ready. Critical/Major: none. Minor (fixed by orchestrator): credentials-file write/corruption errors now surface path-only actionable messages via apperror instead of generic 'internal error'. Suggestions deferred: os.CreateTemp for the temp filename; newAuthCommand accepting deps. Adversarial token-leak and permission checks all clean (0600 temp+rename, 0700 dir, no token in argv/output/errors). Config validation relaxation (empty credential reference now allowed) was a mandatory pre-check correction, backward compatible.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added fj auth login: hidden TTY or --token-stdin input, token verified via GET /api/v1/user before storing in a 0600 credentials.toml separate from config, env-var precedence preserved, resolve failures now hint at fj auth login. x/term v0.45.0 and x/sys v0.47.0 added with TASK-15 notices. Verified with go test ./..., make pre-commit, smoke test, and an independent Fable review (Review Ready).
<!-- SECTION:FINAL_SUMMARY:END -->
