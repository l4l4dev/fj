---
id: TASK-12
title: Select and add project license
status: Done
assignee:
  - '@codex'
created_date: '2026-07-12 00:36'
updated_date: '2026-07-12 07:29'
labels: []
dependencies: []
references:
  - doc-1
ordinal: 91000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Evaluate suitable open-source license options for fj, obtain an explicit human decision on the exact license, and add the approved license text in a separate focused change. This task must not select a license automatically and must not change licensing files before the human decision is recorded.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Relevant license options and their project implications are presented without selecting one automatically.
- [x] #2 The exact license is explicitly approved by a human before repository files are changed.
- [x] #3 The approved license text is added in a dedicated LICENSE file without unrelated changes.
- [x] #4 User-facing and package metadata are checked for consistency with the approved license where applicable.
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Record the explicit Human Decision selecting MIT License and copyright holder l4l4dev.
2. Add the standard OSI MIT License text in the repository root with Copyright (c) 2026 l4l4dev.
3. Add a minimal License section to README.md without restructuring unrelated content.
4. Assess CONTRIBUTING.md and third-party notice needs without changing or adding those files.
5. Verify license text, copyright, README link, task Decision record, diff scope, and whitespace.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Human Decision — License Selection:
- Adopt the MIT License.
- Copyright Holder: `l4l4dev`.
- Rationale: lower adoption and contribution barriers; permit commercial use, forks, modification, and redistribution; minimize maintenance burden for an individually maintained Go CLI; and align with common practice around GitHub CLI and Gitea-related tooling.
- The project prioritizes the simplicity and broad adoption of MIT over an explicit patent grant at this stage.
- Proprietary forks and non-public modifications are explicitly accepted consequences of this choice.
- This Decision is a project-governance choice and does not constitute legal advice.

Model: GPT-5 (Codex) — license implementation is policy-sensitive repository work requiring exact text, scope control, and verification.

Pre-Implementation Check:
- Constitution, architecture, roadmap, scope, compatibility, security, and task granularity: Pass.
- Additional license-selection approval: not required; the Human Decision above is explicit.
- Approved implementation scope: root LICENSE and minimal README License section only; CONTRIBUTING.md and third-party notice changes remain assessment-only.

License Comparison Evidence:

MIT License:
- Category: permissive.
- Commercial use: allowed.
- Modification and redistribution: allowed.
- Source disclosure: not required.
- Patent grant: no explicit patent grant.
- Contributor and enterprise impact: easy to understand with low adoption and contribution barriers; permits proprietary forks and non-public modifications.
- fj fit: low maintenance burden for an individually maintained Go CLI and the strongest fit for the current goals of user adoption and interoperability.

Apache License 2.0:
- Category: permissive.
- Commercial use: allowed.
- Modification and redistribution: allowed.
- Source disclosure: not required.
- Patent grant: explicit contributor patent grant with patent termination provisions.
- Redistribution considerations: preserve the license, identify modified files, and process NOTICE attribution when applicable.
- Contributor and enterprise impact: strong patent clarity and enterprise contribution fit, with somewhat greater compliance explanation than MIT.
- fj fit: a strong candidate, but the project currently prioritizes MIT simplicity.

BSD 3-Clause:
- Category: permissive.
- Commercial use: allowed.
- Modification and redistribution: allowed.
- Source disclosure: not required.
- Patent grant: no explicit patent grant.
- Redistribution considerations: preserve copyright, license conditions, and disclaimer; contributor names may not be used for endorsement without permission.
- Contributor and enterprise impact: generally easy to adopt, with a clearer non-endorsement condition than MIT.
- fj fit: highly compatible, but offers no decisive advantage over MIT or Apache 2.0 for the current project.

Mozilla Public License 2.0:
- Category: file-level weak copyleft.
- Commercial use: allowed.
- Modification and redistribution: allowed.
- Source disclosure: modified MPL-covered files distributed to others must be made available under MPL terms.
- Patent grant: explicit patent grant.
- Contributor and enterprise impact: balances contribution return with coexistence alongside proprietary code, but requires file-level compliance explanation and management.
- fj fit: operationally complex for a small single-binary Go CLI.

GNU General Public License v3:
- Category: strong copyleft.
- Commercial use: allowed.
- Modification and redistribution: allowed.
- Source disclosure: distribution of a covered derivative work requires corresponding source and GPLv3 terms.
- Patent grant: explicit patent provisions.
- Contributor and enterprise impact: strongly preserves downstream software freedom, while potentially increasing legal and compliance burden for enterprise use, embedding, and redistribution.
- fj fit: appropriate when copyleft is the primary goal, but MIT better matches the current goal of reducing adoption barriers.

Comparison Conclusion and Decision Sequence:
1. The research presented relevant license candidates, trade-offs, and project implications; the AI did not make the license decision automatically.
2. The final selection was left to a human.
3. The Human Decision approved the MIT License and copyright holder `l4l4dev`.
4. That Decision was recorded in TASK-12 before repository licensing files were changed.
5. Only after the Decision record, the root LICENSE and README License section were implemented.
6. The comparison does not reopen or alter the approved MIT selection or copyright holder.

Deferred Items:

Dependency license inventory:
- Before a Public Release, perform a formal license inventory for direct and indirect dependencies.
- Confirm source and binary redistribution conditions.
- Determine how required license texts, copyright notices, and NOTICE attribution will accompany distributed artifacts.
- Decide whether to use THIRD_PARTY_NOTICES.md based on the inventory result.
- This work is not a TASK-12 completion condition; managing it in an independent Backlog task is recommended. No follow-up task is created by this remediation.

Contribution policy:
- Before accepting external contributions at scale, confirm the inbound licensing policy.
- Decide whether a Developer Certificate of Origin is needed.
- Decide whether a Contributor License Agreement is needed.
- If dual licensing or copyright assignment becomes necessary, obtain a separate Human Decision.
- This work is not a TASK-12 completion condition.

Evidence Summary:
- MIT License, Apache License 2.0, BSD 3-Clause, MPL 2.0, and GPL v3 were compared across commercial use, modification and redistribution, source-disclosure obligations, patent provisions, contributor and enterprise impact, and suitability for fj.
- The comparison presented candidates and trade-offs without selecting a license automatically. The final selection was explicitly left to a human.
- The Human Decision selected the MIT License with copyright holder `l4l4dev`.
- The Decision was recorded before repository licensing files were changed.
- The root LICENSE contains the standard MIT License text without omission, additional conditions, or modification, using `Copyright (c) 2026 l4l4dev`.
- README.md identifies MIT as the project license and links to the root LICENSE file.
- Existing Go package, CLI metadata, artifact naming, and workflow metadata were checked and require no additional license fields or changes.
- Verification completed successfully: MIT text comparison, copyright check, README link check, Decision record check, `git diff --check`, `go vet ./...`, `go test ./...`, and `make pre-commit`.

Independent Review (post-implementation re-review):
- Scope: TASK-12, LICENSE, README.md, package and distribution metadata, current diff scope, license comparison evidence, Decision sequence, and Deferred Items.
- Previous finding MAJ-1: Resolved. The five-license comparison and project implications are now recorded in TASK-12 as persistent Evidence.
- Acceptance Criteria #1-#4: Pass.
- Critical: none.
- Major: none.
- Minor: none.
- The review confirmed the approved MIT text and copyright holder, Decision-gate compliance, README and metadata consistency, and absence of unrelated changes.
- Dependency license inventory and contribution-policy decisions remain appropriately deferred and do not block TASK-12 completion.
- Result: Ready for Finalization.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Selected and added the MIT License following an explicit Human Decision. TASK-12 records the comparison of MIT, Apache License 2.0, BSD 3-Clause, MPL 2.0, and GPL v3 without automatic license selection, and records approval of MIT with copyright holder `l4l4dev` before implementation.

The repository now contains the standard MIT License text in the root LICENSE file, and README.md includes a concise MIT License section linking to it. Existing package and distribution metadata were checked and require no additional license changes.

Verification and independent re-review passed with no Critical, Major, or Minor findings. Dependency license inventory and external-contribution policy remain separately deferred and do not block TASK-12 completion.
<!-- SECTION:FINAL_SUMMARY:END -->
