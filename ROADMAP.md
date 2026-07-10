# Roadmap

This roadmap describes the long-term evolution of `fj` through independently reviewable milestones. It contains no schedule or task assignments. All milestones must remain consistent with [PROJECT_CONSTITUTION.md](PROJECT_CONSTITUTION.md), follow the workflow in [AGENTS.md](AGENTS.md), and preserve the boundaries defined in [ARCHITECTURE.md](ARCHITECTURE.md).

## M0 - Foundation

### Goal

Establish the project principles, architectural direction, development workflow, and repository standards required for sustainable collaboration.

### Scope

- Project constitution and governance principles
- AI-agent working rules and approval gates
- High-level architecture and dependency boundaries
- Milestone-based roadmap
- Repository-wide documentation, review, testing, and commit standards

### Success Criteria

- The project's purpose, values, and human decision authority are documented.
- Contributors and AI agents have clear, consistent working rules.
- The high-level architecture and dependency direction are documented without premature implementation commitments.
- Repository standards support small, reviewable, and traceable changes.
- Project documents agree on their authority and relationships.

## M1 - Core CLI

### Goal

Provide a dependable CLI foundation that can connect users safely to the intended Forgejo instance.

### Scope

- CLI framework and global command behavior
- Configuration lifecycle and validation
- Authentication and credential safety
- Multiple Forgejo instance profiles and selection
- Repository context selection
- Common human-readable errors and output behavior

### Success Criteria

- Users can identify and select a configured Forgejo instance predictably.
- Authentication succeeds without exposing credentials in output or diagnostics.
- Repository context is explicit and observable.
- Common command, configuration, and failure behavior is consistent.
- The foundation respects the architectural boundaries in `ARCHITECTURE.md`.

## M2 - Repository

### Goal

Support the essential Forgejo repository workflows through a consistent CLI experience.

### Scope

- Repository discovery and inspection
- Repository creation and lifecycle operations
- Repository metadata and settings
- Collaborator and access-related views appropriate to repository workflows

### Success Criteria

- Users can discover and inspect repositories across configured instances.
- Supported repository changes make their target and effect clear before execution.
- Repository operations use consistent inputs, outputs, and error semantics.
- Behavior is documented and covered by relevant compatibility expectations.

## M3 - Issues

### Goal

Enable complete, safe, and predictable issue-management workflows.

### Scope

- Issue discovery and inspection
- Issue creation and updates
- Issue state transitions
- Comments, labels, milestones, and assignments
- Issue filtering and pagination behavior

### Success Criteria

- Users can complete common issue workflows without leaving the CLI.
- State-changing operations clearly identify the issue, repository, and instance.
- Issue metadata and lifecycle behavior remain consistent across commands.
- Large issue collections can be navigated predictably.

## M4 - Pull Requests

### Goal

Support pull request collaboration from creation through review and completion while preserving human control.

### Scope

- Pull request discovery and inspection
- Pull request creation and updates
- Review, comment, and approval workflows
- Merge readiness and status visibility
- Merge and close operations with explicit safeguards

### Success Criteria

- Users can follow a pull request from creation through completion in the CLI.
- Review state, checks, and merge readiness are presented clearly.
- Consequential operations require an explicit target and intent.
- Pull request behavior remains predictable for both interactive and automated use.

## M5 - Releases

### Goal

Provide reliable release-management workflows for Forgejo repositories.

### Scope

- Release discovery and inspection
- Draft and published release lifecycle
- Tags, release notes, and release metadata
- Release asset management
- Pre-release and latest-release semantics

### Success Criteria

- Users can inspect and manage the supported release lifecycle through the CLI.
- Release, tag, and asset relationships are presented unambiguously.
- Publishing and deletion operations make their consequences explicit.
- Release behavior is documented and compatibility-sensitive interfaces remain stable.

## M6 - Actions

### Goal

Make Forgejo Actions activity observable and controllable through safe CLI workflows.

### Scope

- Workflow and run discovery
- Run status, job, and log visibility
- Supported workflow dispatch and run-control operations
- Artifact discovery and retrieval
- Clear representation of asynchronous execution states

### Success Criteria

- Users can determine the current and final state of supported workflow runs.
- Jobs, logs, and artifacts are associated with the correct instance, repository, and run.
- Run-control operations communicate their target and expected effect.
- Asynchronous states and remote failures are represented consistently.

## M7 - AI Experience

### Goal

Make `fj` a stable, predictable interface for AI agents and automation without reducing human readability or oversight.

### Scope

- Stable JSON output across supported command families
- Machine-readable success and error contracts
- Deterministic non-interactive behavior
- Explicit operation context, effects, and next-step information
- Compatibility policy for machine-readable interfaces
- AI-friendly command discovery and usage guidance

### Success Criteria

- Supported commands provide valid, documented, and stable JSON output.
- Automation can distinguish successful results, user errors, remote failures, and internal failures.
- Machine-readable output contains sufficient context without exposing sensitive data.
- Non-interactive execution avoids ambiguous prompts and hidden state.
- AI-focused behavior remains understandable and verifiable by people.

## M8 - Advanced Features

### Goal

Extend protocol support, integration options, and performance while preserving the established architecture and compatibility guarantees.

### Scope

- GraphQL support alongside REST
- Reviewed extension boundaries for additional integrations
- Performance and resource-efficiency improvements
- Advanced instance and repository discovery
- Broader interoperability with Forgejo-compatible environments

### Success Criteria

- GraphQL capabilities coexist with REST without changing core use-case behavior.
- Extensions use explicit, documented boundaries with clear ownership and security expectations.
- Performance improvements are measurable and do not weaken correctness or observability.
- Existing commands, configuration, and JSON contracts remain compatible or provide documented migration paths.
- Advanced capabilities preserve human oversight and the principles in `PROJECT_CONSTITUTION.md`.
