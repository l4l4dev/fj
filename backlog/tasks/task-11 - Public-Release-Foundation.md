---
id: TASK-11
title: Public Release Foundation
status: Done
assignee: []
created_date: '2026-07-11 17:33'
updated_date: '2026-07-12 07:02'
labels: []
milestone: m-10
dependencies: []
priority: medium
ordinal: 11000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Establish the foundation for reproducible fj release builds, cross-platform artifacts, SHA-256 checksums, installation guidance, and release acceptance. This task does not create or publish public releases. TASK-10.5 is not a required dependency; its results are inputs to release acceptance.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Release workflow foundation is defined without public release creation or publication.
- [x] #2 Version policy uses vMAJOR.MINOR.PATCH, strips the v prefix, and preserves dev fallback.
- [x] #3 Initial build targets include darwin/arm64 and linux/amd64.
- [x] #4 SHA-256 checksum and artifact management policies are defined.
- [x] #5 TASK-10.1 version metadata and User-Agent integration are preserved.
- [x] #6 Release acceptance accounts for TASK-10.5 results without treating failed commands as successful.
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Decision Summary:
- TASK-11 establishes a reproducible Public Release Foundation without creating or publishing a public GitHub Release.
- Version inputs use `vMAJOR.MINOR.PATCH`, are normalized by removing the leading `v`, and are injected through the TASK-10.1 ldflags path.
- Local or uninjected builds preserve the `dev` fallback.
- Initial artifact targets are limited to `darwin/arm64` and `linux/amd64`.
- The foundation generates a deterministic SHA-256 manifest and a consolidated workflow artifact containing exactly both binaries and the checksum manifest.
- Public Release publication, tag-only public release policy, signing, notarization, Homebrew, system-wide installation, and permanent distribution remain outside TASK-11.

Evidence Summary:
- TASK-11.1 through TASK-11.5 are Done.
- TASK-11.1 verified version validation, leading-`v` normalization, ldflags injection, read-only workflow permissions, and absence of public Release creation.
- TASK-11.2 verified `darwin/arm64` and `linux/amd64` builds, artifact naming, binary formats, and injected version metadata.
- TASK-11.3 verified deterministic SHA-256 generation, checksum verification, tamper detection, missing-input detection, and consolidated artifact contents.
- TASK-11.4 verified retrieval, checksum, current-user installation, version/help, uninstall documentation, and distribution/security boundaries.
- TASK-11.5 verified macOS arm64 and Linux amd64 runtime, installed version/help, read-only Forgejo smoke operations, secret redaction, and successful live `workflow_dispatch`.
- TASK-10.1 metadata continuity is preserved through `internal/version.Value`, CLI version output, and `fj/<version>` Forgejo User-Agent behavior.
- TASK-10.5 `repo list` and `pr list` remain unsuccessful and are not included in TASK-11 successful evidence.

Independent Review (post-implementation, parent TASK-11):
- Scope: TASK-11, TASK-11.1 through TASK-11.5, `release-foundation.yml`, installation documentation, release acceptance evidence, and TASK-10.1/TASK-10.5 boundaries.
- Critical: none.
- Major: none.
- Minor: none.
- Suggestions: none.
- Acceptance Criteria #1-#6 are supported by implementation and recorded evidence.
- No public Release is created.
- Version and User-Agent metadata continuity is preserved.
- Both platform artifacts and checksum policy are verified.
- Secret-redaction boundaries are maintained.
- TASK-10.5 unsuccessful commands are not treated as successful.
- Result: Ready for Finalization.

Deferred Items:
- License selection remains TASK-12 and requires a separate explicit human decision.
- Public GitHub Release creation and publication remain outside the Public Release Foundation.
- A future public Release workflow must define its tag-only release policy separately.
- Signing, notarization, Homebrew, system-wide installation, and permanent distribution channels remain future work.
- TASK-10.5 `repo list` and `pr list` remain unsuccessful. Future re-acceptance belongs to TASK-10.5 and is not a TASK-11 completion requirement.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Completed the Public Release Foundation across TASK-11.1 through TASK-11.5. The foundation validates and normalizes `vMAJOR.MINOR.PATCH` versions, preserves the TASK-10.1 ldflags metadata and User-Agent path, builds versioned `darwin/arm64` and `linux/amd64` artifacts, generates and verifies deterministic SHA-256 checksums, consolidates the expected release artifacts, and documents safe retrieval, verification, current-user installation, version/help confirmation, and uninstall behavior.

Release acceptance verified both platform artifacts, the live workflow, read-only Forgejo smoke operations, and secret redaction. TASK-10.5 unresolved `repo list` and `pr list` commands remain explicitly unsuccessful and were not counted as release successes.

Public Release publication, License selection, tag-only public release policy, signing, notarization, Homebrew, system-wide installation, and permanent distribution remain outside this task. Parent-level Independent Review found no Critical, Major, or Minor findings and judged TASK-11 Ready for Finalization.
<!-- SECTION:FINAL_SUMMARY:END -->
