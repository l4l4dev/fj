---
id: TASK-2.10
title: Load configuration from an XDG TOML file
status: Done
assignee:
  - '@codex'
created_date: '2026-07-10 17:08'
updated_date: '2026-07-10 17:16'
labels: []
dependencies:
  - TASK-2.3
  - TASK-2.4
references:
  - ROADMAP.md
modified_files:
  - go.mod
  - go.sum
  - internal/infrastructure/config/loader.go
  - internal/infrastructure/config/loader_test.go
parent_task_id: TASK-2
ordinal: 76000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Load named Forgejo instance profiles from an XDG Base Directory TOML configuration file using BurntSushi/toml. Keep loading separate from the existing configuration model and validation responsibilities.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The configuration file is resolved according to the XDG Base Directory convention.
- [x] #2 TOML profiles are decoded with BurntSushi/toml and passed through existing configuration validation.
- [x] #3 Missing, malformed, and invalid configuration files produce actionable errors without exposing credentials.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Resolve the approved XDG configuration path, rejecting relative XDG_CONFIG_HOME and unset HOME.
2. Add a private TOML storage DTO decoded with BurntSushi/toml, map it to the existing Application configuration model, and run Configuration.Validate().
3. Return safe actionable errors for missing, malformed, and invalid files without exposing credential values; leave permission and symlink behavior unspecified.
4. Add deterministic tests for path precedence, environment errors, missing/malformed/invalid files, and credential redaction.
5. Run required checks and obtain an independent GPT-5 post-implementation review, then finalize TASK-2.10.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Model: GPT-5 — adding BurntSushi/toml and an Infrastructure configuration loader is a Major Change requiring high-capability design and verification

Pre-implementation check (GPT-5): Critical none. Major: external dependency is approved in principle, but exact XDG_CONFIG_HOME path semantics, HOME fallback, required missing-file behavior, and TOML decode-error secret redaction require explicit human approval before implementation. Minor: file permission/symlink and HOME-unset policy should be consciously scoped. Suggestion: private TOML DTO, map into Application model, call existing Validate(), and table-test path/missing/malformed/invalid/redaction cases. No code implementation started.

Human approved the Major Change details: absolute XDG_CONFIG_HOME takes precedence, unset XDG_CONFIG_HOME falls back to $HOME/.config, relative XDG_CONFIG_HOME and unset HOME are errors, the config file is required, TOML/validation errors redact credential values, private DTO mapping calls existing Configuration.Validate(), and file permissions/symlinks remain unspecified for a future Task.

Implemented the approved Infrastructure configuration loader. Absolute XDG_CONFIG_HOME resolves to $XDG_CONFIG_HOME/fj/config.toml; unset XDG_CONFIG_HOME falls back to $HOME/.config/fj/config.toml; relative XDG_CONFIG_HOME and unset HOME return configuration errors. A required TOML file is decoded with BurntSushi/toml into private DTOs, mapped to Application Configuration, and validated with existing Configuration.Validate(). Missing, malformed, and invalid errors are actionable and never include credential values. File permissions and symlink behavior remain out of scope.

Validation passed: gofmt -l . (no output); git diff --check; focused go test ./internal/infrastructure/config; go vet ./...; go test ./.... Initial unsandboxed full verification hit the known Go build-cache permission restriction and passed on approved elevated reruns. go mod tidy normalized BurntSushi/toml as a direct dependency.

Post-implementation review (independent GPT-5): Critical none; Major none; Minor none; Suggestion none after go.mod direct-dependency cleanup. Review confirmed approved scope, acceptance criteria, Infrastructure-to-Application dependency direction, secure redaction, sufficient tests, and accurate Backlog records.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added the approved XDG TOML configuration loader using BurntSushi/toml. It resolves deterministic XDG/HOME paths, decodes through private DTOs, maps to the Application configuration model, runs existing validation, and safely reports missing/malformed/invalid files without credential leakage. Permissions and symlinks remain unspecified for future work. All checks and the independent Major Change review passed.
<!-- SECTION:FINAL_SUMMARY:END -->
