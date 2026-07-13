---
id: TASK-15
title: Establish Dependency License and Redistribution Compliance
status: Done
assignee:
  - '@codex'
created_date: '2026-07-12 13:17'
updated_date: '2026-07-13 03:31'
labels: []
dependencies:
  - TASK-11
  - TASK-12
modified_files:
  - THIRD_PARTY_NOTICES.md
  - licenses/BurntSushi-toml-COPYING
  - licenses/spf13-cobra-LICENSE.txt
  - licenses/spf13-pflag-LICENSE
  - licenses/Go-LICENSE
  - README.md
  - >-
    backlog/tasks/task-15 -
    Establish-Dependency-License-and-Redistribution-Compliance.md
ordinal: 94000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Establish the dependency-license inventory, redistribution obligations, compliance-artifact policy, and Public Release gate required before fj binaries are publicly distributed. This parent task records the completed investigation and required Human Decisions; it does not authorize implementation, Public Release creation, workflow changes, dependency changes, or release-asset changes.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The dependency inventory distinguishes direct, indirect, target-linked, test/tooling-only, platform-specific, standard-library, copied, generated, embedded, and vendored code as applicable
- [x] #2 Every dependency included in the darwin/arm64 or linux/amd64 release binary has authoritative version, upstream license, copyright, NOTICE, exception, and ambiguity Evidence
- [x] #3 Source, binary, and release-artifact redistribution obligations are documented for every applicable license family, including MIT, Apache License 2.0, and BSD-style notices
- [x] #4 An explicit Human Decision selects the compliance-artifact format, downloadable distribution unit, inventory persistence policy, automation policy, and handling of unresolved attribution ambiguity before implementation begins
- [x] #5 The approved release design provides required project and third-party license texts, copyright notices, disclaimers, and NOTICE attribution to recipients of every distributed binary
- [x] #6 A repeatable Public Release gate verifies target-specific binary dependencies, notice-bundle completeness, inventory freshness, artifact contents, and checksum coverage
- [x] #7 Unknown, incompatible, or unresolved dependency-license and attribution conditions block Public Release pending explicit human or legal resolution
- [x] #8 Independent Review confirms the inventory, obligations, approved artifact design, verification Evidence, scope boundaries, and Public Release gate before Finalization
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Preserve the completed investigation and the approved Human Decisions as the durable basis for TASK-15.
2. Prepare an independently reviewable implementation design for the combined notice bundle, per-platform archives, committed inventory snapshot, standard-Go-tool validation, fail-closed release gate, and archive checksum policy.
3. Present the implementation plan and exact allowed files to the human maintainer and obtain separate explicit implementation authorization before creating or changing any compliance artifact, workflow, archive, checksum behavior, or documentation.
4. After that authorization only, implement and verify the approved scope without adding external license scanners, package-manager integration, source-release publication, signing, notarization, or additional platforms.
5. Obtain Independent Review before checking Acceptance Criteria, Finalization, Public Release, commit, or push.

The eight product and compliance Decisions are resolved and recorded. Implementation remains blocked only on a later explicit implementation authorization with approved file and workflow scope.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Investigation Evidence:
- Current Public Release targets are darwin/arm64 and linux/amd64, built with CGO_ENABLED=0.
- Target-specific production dependency analysis identified github.com/BurntSushi/toml v1.6.0, github.com/spf13/cobra v1.10.2, and github.com/spf13/pflag v1.0.9 as linked external modules for both targets.
- github.com/inconshreveable/mousetrap v1.1.0 is present in the module graph but is Windows-specific and is not linked into the current release targets.
- Other modules present only through transitive go.mod relationships are not reachable from the current production or test package graph and must not be assumed to be distributed.
- BurntSushi/toml is MIT-licensed and includes type_fields.go adapted from Go encoding/json under a Go Authors BSD-style notice; this secondary notice requires explicit treatment.
- Cobra is Apache License 2.0, includes Cobra Authors copyright headers, and has no NOTICE file at v1.10.2.
- pflag uses a BSD 3-Clause license with Alex Ogier and Go Authors copyright notices.
- The current consolidated release-foundation artifact contains only two binaries and the SHA-256 manifest; it does not contain project or third-party license and notice files.
- The root fj MIT LICENSE covers fj-owned code and does not replace third-party redistribution obligations.
- Static linking into a Go executable does not remove applicable license-text, copyright, disclaimer, or attribution obligations.
- This record is engineering compliance planning and is not legal advice.

Human Decisions Required Before Implementation:
1. Compliance artifact format: THIRD_PARTY_NOTICES.md, a licenses directory, or a combined summary plus verbatim license directory.
2. Downloadable distribution unit: loose binaries, per-platform archives, or another explicitly approved structure that reliably accompanies each binary with required notices.
3. Inventory persistence: whether the verified inventory and notice mapping are committed or generated only for each release.
4. Automation policy: standard Go tooling and local deterministic validation only, or adoption of a separately reviewed license-scanning or SBOM tool.
5. Ambiguity gate: whether any unknown, incompatible, or unresolved attribution condition blocks Public Release; the investigation recommends blocking.
6. Checksum policy: whether checksums cover archives or each contained compliance file.
7. Package-manager scope: whether Homebrew and other package-manager requirements are included now or deferred.
8. Source distribution treatment: how GitHub-generated or explicit source archives provide required notices.

No option above is approved by this task record. Implementation must not begin until the maintainer records the selected Decisions.

Scope Boundaries:
- Allowed after future approval: verified inventory records, approved third-party notice artifacts, approved release packaging and checksum changes, deterministic dependency/license validation, and minimal notice-location documentation.
- Prohibited now: implementation, source changes, workflow changes, README or LICENSE changes, release-asset creation, dependency changes, new tools, subtasks, Public Release, tag, commit, and push.
- TASK-11 and TASK-12 remain unchanged Sources of Truth for the release foundation and fj project-license selection.

Deferred Items:
- SPDX or other SBOM generation may be considered separately after the first compliant Public Release path is established.
- GitHub dependency review or license scanning may be added later only through a separate dependency/tooling Decision.
- Homebrew, other package managers, signing, notarization, Windows support, and additional platform-specific compliance remain outside the initial task unless explicitly adopted by Human Decision.
- mousetrap notice inclusion becomes release-relevant if a Windows artifact is approved later.
- Legal advice or jurisdiction-specific interpretation is outside engineering scope; unresolved legal ambiguity blocks release until resolved by the maintainer with appropriate counsel.
- Subtasks may be proposed after the Human Decisions are recorded but are not created by this task.

Model: GPT-5 (Codex) — dependency-license inventory and redistribution boundaries require cross-cutting release, legal-policy, artifact, and verification reasoning.

Approved Human Decisions — Dependency License and Redistribution Compliance:

1. Compliance artifact format:
- Adopt the combined format THIRD_PARTY_NOTICES.md plus licenses/.
- THIRD_PARTY_NOTICES.md is the human-readable dependency inventory and attribution index.
- licenses/ contains verbatim upstream license or notice files required for distribution.
- Summary text does not replace required verbatim license text.
- The documentation must distinguish an upstream NOTICE file from fj’s own third-party notice summary.

2. Downloadable distribution unit:
- Use one per-platform archive as the primary downloadable binary distribution unit.
- Initial archive names are fj-<version>-darwin-arm64.tar.gz and fj-<version>-linux-amd64.tar.gz.
- Each archive contains fj, LICENSE, THIRD_PARTY_NOTICES.md, and licenses/.
- Loose standalone binaries are not treated as the primary compliant Public Release unit.
- This Decision does not remove or change the current release-foundation artifact behavior.

3. Inventory storage:
- Commit a reviewed dependency-license inventory snapshot to the repository.
- The snapshot identifies module and version, distinguishes linked production dependencies from non-distributed modules, identifies license and notice sources, remains reviewable in Git, and is verified against actual release binaries before Public Release.
- The exact filename and format remain an implementation-planning detail.

4. Initial automation approach:
- For the initial Public Release, use standard Go tooling plus explicit reviewed mappings.
- Expected Evidence may include go list -m -json all, target-specific go list -deps, and go version -m <release-binary>.
- TASK-15 does not adopt go-licenses, an SPDX or SBOM generator, or another external license scanner unless a later explicit Human Decision approves it.
- Automated classification does not replace manual legal interpretation.

5. Ambiguity gate:
- Unknown, incompatible, missing, or unresolved license or attribution conditions block Public Release.
- Warnings alone are insufficient for unresolved distributed dependencies.
- Resolution may require a corrected inventory, additional notice text, a dependency change, a scope change, or explicit legal review.
- TASK-15 does not provide legal advice.

6. Checksum policy:
- Checksums cover the per-platform archives as the primary distribution units.
- Compliance files inside each archive are indirectly protected by the archive checksum.
- An additional checksum for standalone compliance files is an implementation detail and must not weaken archive verification.

7. Package-manager scope:
- Homebrew and other package-manager-specific compliance requirements are deferred.
- TASK-15 should produce a distribution structure reasonably reusable later, but package-manager formulae, metadata, and channel-specific notice placement are outside TASK-15 completion scope.

8. Source archive treatment:
- The final policy for GitHub-generated versus explicitly created source archives is deferred to the future Public Release publication task.
- TASK-15 documents identified notice implications but does not redesign source distribution or publish a Release.

Decision Rationale:
- A human-readable summary plus verbatim license files clearly separates inventory and attribution guidance from authoritative license text.
- Per-platform archives ensure each primary binary distribution unit contains project and third-party licensing materials.
- A committed inventory improves Git reviewability and detection of unnoticed dependency changes.
- Standard Go tooling is sufficient for the small initial dependency graph and avoids premature tooling dependencies.
- Fail-closed handling of unresolved attribution reduces Public Release compliance risk.
- Archive checksums align integrity verification with the actual downloadable distribution unit.
- Homebrew, SBOM, signing, notarization, and additional platforms must not expand the first compliance milestone.
- Source archive publication policy belongs to the later Public Release Decision.

Implementation Authorization and Scope Boundary:
- These Decisions authorize implementation planning within TASK-15.
- They establish the design constraints for later implementation, but do not themselves authorize repository implementation changes.
- Separate explicit human authorization is required before implementation begins and must identify the allowed files and workflow scope.
- Not yet authorized: creating THIRD_PARTY_NOTICES.md; creating licenses/; changing workflows; generating archives; changing checksum implementation; changing dependencies; adopting external tooling; creating or publishing a Public Release; creating tags; publishing source archives; or integrating Homebrew.
- Subtasks are not created by this Decision-recording step.
- TASK-15 remains In Progress with assignee @codex, Acceptance Criteria #1-#8 unchecked, and no Final Summary.

Decision Record Supersession:
- The earlier Human Decisions Required Before Implementation list records the pre-Decision investigation state.
- It is superseded in full by the Approved Human Decisions recorded above and is no longer an unresolved Decision list.
- All eight listed Decisions are resolved. Only separate explicit implementation authorization and approved file/workflow scope remain pending.

Implementation Progress:
- Added THIRD_PARTY_NOTICES.md as the committed, human-readable inventory and attribution index for the initial darwin/arm64 and linux/amd64 binaries.
- Recorded the three target-linked production modules and versions: BurntSushi/toml v1.6.0, spf13/cobra v1.10.2, and spf13/pflag v1.0.9.
- Recorded mousetrap and the other module-graph-only modules as not distributed in the initial binaries, with the reason for exclusion.
- Added verbatim upstream license files for BurntSushi/toml, spf13/cobra, spf13/pflag, and the Go 1.26.5 distribution.
- Distinguished fj’s THIRD_PARTY_NOTICES.md summary from an upstream Apache NOTICE file and recorded that Cobra v1.10.2 has no upstream NOTICE file.
- Added the minimal README link to THIRD_PARTY_NOTICES.md.
- No source code, tests, workflow, dependency, go.mod, go.sum, release artifact, archive, checksum behavior, tag, or Public Release was changed.

Implementation Verification:
- Byte-for-byte comparison of all four files in licenses/ against the installed upstream module or Go distribution source: Pass.
- Target-specific production dependency listing for darwin/arm64 with CGO_ENABLED=0: Pass; toml v1.6.0, cobra v1.10.2, and pflag v1.0.9.
- Target-specific production dependency listing for linux/amd64 with CGO_ENABLED=0: Pass; toml v1.6.0, cobra v1.10.2, and pflag v1.0.9.
- Required notice files and README/notice links: Pass.
- Cobra v1.10.2 upstream NOTICE absence check: Pass.
- git diff --check: Pass.
- Go tooling emitted sandbox stat-cache write warnings, but both dependency-list commands completed successfully.
- Acceptance Criteria remain unchecked, Status remains In Progress, and no Final Summary or Independent Review is recorded.

Implementation Authorization Update:
- The human’s TASK-15 Implementation request supersedes the earlier pending-authorization statement only for the compliance artifacts implemented above and the minimum README consistency link.
- Authorized files are THIRD_PARTY_NOTICES.md, licenses/, README.md, and this TASK-15 implementation record.
- Workflow, archive generation, checksum behavior, dependency, source, test, release artifact, Public Release, tag, commit, push, Acceptance Criteria, Status, Final Summary, and Independent Review changes remain unauthorized or outside this implementation phase.

Public Release Gate Definition:

Purpose and boundary:
- TASK-15 defines the release-compliance conditions that must be true before a Public Release.
- TASK-15 does not implement or execute release packaging, workflow automation, archive generation, checksum generation, artifact upload, or live release acceptance.

Required Distribution Design:
- The primary binary distribution units are fj-<version>-darwin-arm64.tar.gz and fj-<version>-linux-amd64.tar.gz.
- Each archive contains exactly the required distribution payload: fj, LICENSE, THIRD_PARTY_NOTICES.md, and licenses/.
- Loose standalone binaries are not treated as the primary compliant Public Release unit.
- This is a mandatory future distribution design and is not packaging implementation completed by TASK-15.

Required Dependency Evidence:
Before Public Release, the release process provides Evidence that:
- Both darwin/arm64 and linux/amd64 release binaries were built.
- go version -m was executed against both actual release binaries.
- Binary module metadata matches the approved committed dependency inventory.
- Target-specific production dependencies were checked for both targets.
- Every linked third-party module has authoritative license, copyright, disclaimer, exception, and NOTICE-presence Evidence as applicable.
- Unknown or changed linked modules fail the gate.
- Inventory entries not linked into the target binaries are explicitly classified as non-distributed with their exclusion reason.
- All required copyright, disclaimer, license, and NOTICE obligations are resolved.

Required Compliance Files:
The future release gate verifies that:
- The root LICENSE is included.
- THIRD_PARTY_NOTICES.md is included.
- Every required verbatim file under licenses/ is included.
- No required license file is missing.
- Verbatim license files match the reviewed upstream sources for the exact dependency versions.
- Upstream NOTICE presence or absence is confirmed for each exact dependency version.
- fj’s THIRD_PARTY_NOTICES.md is represented as fj’s own inventory and attribution summary and is not represented as an upstream Apache NOTICE file.

Required Integrity Evidence:
The future release implementation verifies that:
- Checksums cover each per-platform archive as the primary distribution unit.
- Archive checksum verification succeeds.
- Archive contents match the approved layout.
- Missing compliance files fail verification.
- Dependency inventory drift fails verification.
- Changed or unknown dependency versions fail verification.
- Tampered archives fail checksum verification.
- TASK-15 defines these conditions but does not implement or execute them.

Pass Conditions:
Public Release may proceed only when:
- All linked dependencies are identified.
- All applicable license and attribution obligations are resolved.
- No unknown, incompatible, missing, or unresolved license condition remains.
- The approved notice bundle is included in each primary distribution unit.
- Actual binary metadata and the committed inventory agree for both targets.
- Archive structure and checksum Evidence pass.
- Independent Review of the release-compliance Evidence passes.

Fail Conditions:
Public Release is blocked when any of the following occurs:
- An unknown linked dependency is present.
- A changed dependency or version is not reflected in the inventory.
- A required license file is missing or stale.
- A copyright, disclaimer, exception, or NOTICE requirement is unresolved.
- Actual binary metadata does not match the committed inventory.
- Required archive content is missing.
- Archive checksum verification fails.
- A legal or attribution ambiguity remains unresolved.
- Warnings alone are insufficient for any of these conditions.

Deferred Packaging Implementation:
The following are outside TASK-15 implementation scope:
- Archive-generation code.
- Release workflow modifications.
- Archive checksum generation.
- Archive upload or release-artifact generation.
- Automated archive-content verification.
- Tamper-test and missing-file-test implementation.
- Live release acceptance.
- Source archive publication policy.
- GitHub Release creation or publication.
- Tag creation.
- These items require a later Public Release Preparation, Release Packaging, or Release Workflow Enhancement task with its own Human Decision, implementation authorization, verification, and Independent Review.

TASK-15 Independent Review Entry Conditions:
TASK-15 may proceed to Independent Review when it contains:
- A verified target-specific dependency inventory.
- Authoritative verbatim license files and NOTICE-presence Evidence.
- THIRD_PARTY_NOTICES.md as the committed inventory and attribution index.
- Source, binary, and release-artifact redistribution obligation analysis.
- The approved per-platform archive-layout policy.
- The approved archive-checksum policy.
- Explicit pass and fail conditions for the Public Release gate.
- Clear Deferred Items separating compliance definition from packaging implementation.
- These entry conditions do not claim that future archives, workflow changes, checksums, or live release acceptance have been implemented.

Remediation — MAJ-2 Redistribution Obligation Matrix:

MIT — BurntSushi/toml:
- Source distribution obligations: retain the dependency copyright notice, MIT permission notice, and disclaimer in copies or substantial portions of the software.
- Binary distribution obligations: provide the same copyright, permission, and disclaimer text with the binary distribution; static linking does not remove this condition.
- Per-platform archive contents: include fj’s root LICENSE for fj-owned code, THIRD_PARTY_NOTICES.md as the inventory and attribution index, and licenses/BurntSushi-toml-COPYING as the authoritative dependency text.
- NOTICE handling: the exact toml v1.6.0 distribution has no upstream NOTICE file; no Apache-style NOTICE propagation applies.
- Modification notice requirements: MIT has no separate prominent modified-file marking requirement, but the original copyright and permission notice must be retained.
- fj response: retain the verbatim upstream COPYING file, record the module/version and attribution in THIRD_PARTY_NOTICES.md, and require both files in each future per-platform archive.

Apache License 2.0 — spf13/cobra:
- Source distribution obligations: provide a copy of Apache License 2.0; under Section 4(b), cause modified upstream files to carry prominent notices stating that they were changed; under Section 4(c), retain applicable copyright, patent, trademark, and attribution notices from the upstream source form.
- Binary distribution obligations: provide recipients of the Object form with a copy of Apache License 2.0. Static linking does not remove this requirement.
- Per-platform archive contents: include fj’s root LICENSE, THIRD_PARTY_NOTICES.md, and licenses/spf13-cobra-LICENSE.txt.
- NOTICE handling: if an exact upstream version includes a NOTICE file, propagate its applicable attribution notices as required by Section 4(d). Cobra v1.10.2 has no upstream NOTICE file. THIRD_PARTY_NOTICES.md is fj’s own summary and must not be represented as an upstream Apache NOTICE file.
- Modification notice requirements: Section 4(b) applies when fj distributes modified Cobra source files; fj currently vendors or modifies no Cobra source.
- fj response: retain the verbatim v1.10.2 license, record the exact version, copyright, and NOTICE absence in THIRD_PARTY_NOTICES.md, and re-check NOTICE and modification status whenever the dependency version or source treatment changes.

BSD-style — spf13/pflag and Go:
- Source distribution obligations: retain the applicable copyright notice, redistribution conditions, and disclaimer in redistributed source.
- Binary distribution obligations: reproduce the applicable copyright notice, conditions, and disclaimer in documentation or other materials provided with the binary distribution.
- Per-platform archive contents: include fj’s root LICENSE, THIRD_PARTY_NOTICES.md, licenses/spf13-pflag-LICENSE for pflag, and licenses/Go-LICENSE for the Go standard library and the Go-derived code identified in BurntSushi/toml.
- NOTICE handling: these BSD-style texts do not impose a separate Apache NOTICE-file mechanism; exact-version NOTICE presence remains part of dependency Evidence.
- Modification notice requirements: the applicable BSD-style texts do not impose a separate prominent modified-file notice requirement; their copyright, conditions, disclaimer, and non-endorsement terms remain applicable.
- fj response: retain both verbatim license files, record pflag and Go-derived attribution in THIRD_PARTY_NOTICES.md, avoid endorsement claims, and include the notice bundle in each future per-platform archive.

Cross-license file mapping:
- LICENSE contains the MIT terms for fj-owned code and does not replace third-party terms.
- THIRD_PARTY_NOTICES.md is fj’s human-readable dependency inventory and attribution index; it does not replace verbatim license text and is not an upstream Apache NOTICE file.
- licenses/ contains the reviewed verbatim dependency and Go distribution license texts required by the approved inventory.
- Each future primary per-platform archive must contain LICENSE, THIRD_PARTY_NOTICES.md, and the complete approved licenses/ set alongside the fj binary.
- Source archive publication policy remains deferred, but the known source redistribution obligations above are not deferred.

Evidence Summary:
- THIRD_PARTY_NOTICES.md records the reviewed inventory for the initial darwin/arm64 and linux/amd64 targets, including direct, indirect, production-linked, Windows-specific, module-graph-only, non-distributed, standard-library, test/tooling, vendored, copied, generated, and embedded classifications.
- Target-specific dependency checks identify BurntSushi/toml v1.6.0, spf13/cobra v1.10.2, and spf13/pflag v1.0.9 as the linked external production modules for both initial targets.
- The repository includes byte-for-byte verified upstream license texts for toml, Cobra, pflag, and the Go 1.26.5 distribution under licenses/.
- THIRD_PARTY_NOTICES.md records exact versions, copyright attribution, upstream NOTICE absence, the Go-derived toml source attribution, non-distributed module reasons, and the distinction between fj’s summary and an upstream Apache NOTICE file.
- README.md links users to THIRD_PARTY_NOTICES.md while the root LICENSE remains the license for fj-owned code.
- TASK-15 documents MIT, Apache License 2.0, and BSD-style source, binary, and future per-platform archive obligations, including Apache Sections 4(b), 4(c), and 4(d), modification notices, and file mapping.
- The approved Public Release gate defines target-binary metadata checks, inventory-drift failure, notice-bundle completeness, archive layout, checksum coverage, pass conditions, and fail-closed handling.
- Verification passed for verbatim license comparison, both target-specific production dependency lists, notice-file presence checks, required links and files, git diff --check, go vet ./..., go test ./..., and make pre-commit.

Independent Re-Review:
- Scope: TASK-15, THIRD_PARTY_NOTICES.md, licenses/, README.md, LICENSE, go.mod, go.sum, governing documents, TASK-11/TASK-12 boundaries, current Git scope, and remediation for MAJ-1 and MAJ-2.
- MAJ-1 Inventory classification is incomplete: Resolved.
- MAJ-2 Redistribution obligations are not fully documented: Resolved.
- Acceptance Criteria #1-#8: Pass.
- Critical findings: None.
- Major findings: None.
- Minor findings: None.
- Inventory classification, dependency reachability, license fidelity, NOTICE handling, redistribution matrix, compliance-file mapping, Public Release gate, and Deferred Item boundaries: Pass.
- Result: Ready for Finalization.

Deferred Items:
- Archive-generation code, release workflow changes, archive checksum generation, archive upload, automated archive-content verification, and tamper/missing-file test implementation belong to a later Public Release Preparation, Release Packaging, or Release Workflow Enhancement task.
- Live release acceptance, GitHub Release creation or publication, and tag creation remain future explicitly authorized work.
- The final policy for GitHub-generated versus explicit source archives belongs to the future Public Release publication task.
- Homebrew and other package-manager integration, SBOM generation, external license scanners, signing, notarization, Windows artifacts, and additional platforms remain outside TASK-15.
- A future Windows artifact requires renewed mousetrap license and NOTICE review.
- Unknown or unresolved legal or attribution ambiguity remains a Public Release blocker and may require explicit legal review.
- These items are outside TASK-15 completion scope and were not implemented by this task.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Established fj’s dependency-license inventory and redistribution-compliance foundation for the initial darwin/arm64 and linux/amd64 binaries. The repository now includes THIRD_PARTY_NOTICES.md as the reviewed inventory and attribution index, verbatim dependency and Go license texts under licenses/, and a minimal README link.

TASK-15 records target-linked and non-distributed dependency classifications, authoritative license and NOTICE Evidence, MIT, Apache License 2.0, and BSD-style redistribution obligations, and the approved relationship among LICENSE, THIRD_PARTY_NOTICES.md, and licenses/. It also defines a fail-closed Public Release gate and the required future per-platform archive layout without implementing release packaging or workflow changes.

Independent Re-Review confirmed MAJ-1 and MAJ-2 are resolved, Acceptance Criteria #1-#8 pass, and there are no Critical, Major, or Minor findings. Future archive, workflow, checksum, release-publication, package-manager, signing, notarization, SBOM, and additional-platform work remains explicitly deferred.
<!-- SECTION:FINAL_SUMMARY:END -->
