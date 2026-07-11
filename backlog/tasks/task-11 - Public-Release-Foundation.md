---
id: TASK-11
title: Public Release Foundation
status: To Do
assignee: []
created_date: '2026-07-11 17:33'
updated_date: '2026-07-11 19:35'
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
- [ ] #1 Release workflow foundation is defined without public release creation or publication.
- [ ] #2 Version policy uses vMAJOR.MINOR.PATCH, strips the v prefix, and preserves dev fallback.
- [ ] #3 Initial build targets include darwin/arm64 and linux/amd64.
- [ ] #4 SHA-256 checksum and artifact management policies are defined.
- [ ] #5 TASK-10.1 version metadata and User-Agent integration are preserved.
- [ ] #6 Release acceptance accounts for TASK-10.5 results without treating failed commands as successful.
<!-- AC:END -->
