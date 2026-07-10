# Architecture

This document defines the high-level architecture of `fj`. It is subordinate to [PROJECT_CONSTITUTION.md](PROJECT_CONSTITUTION.md), and work governed by it must follow [AGENTS.md](AGENTS.md).

## 1. Goals

The architecture must support a reliable, maintainable Go CLI for interacting with Forgejo. It should:

- keep business behavior independent of command-line and remote API technologies;
- provide predictable interfaces for both people and AI agents;
- support multiple Forgejo instances without coupling operations to a single server;
- accommodate REST APIs today and GraphQL in the future;
- make features independently testable and replaceable;
- preserve backward compatibility and allow incremental evolution;
- remain understandable to maintainers and contributors over the long term.

## 2. Design Principles

- **Clean Architecture:** Separate policy from delivery mechanisms, external services, and storage.
- **Explicit boundaries:** Define clear responsibilities and contracts between architectural layers.
- **Dependency inversion:** Core policies depend on abstractions, not infrastructure implementations.
- **Small, cohesive components:** Organize behavior by responsibility and avoid broad shared utilities.
- **Stable interfaces:** Treat commands, JSON output, configuration, and exit behavior as user-facing contracts.
- **Technology isolation:** Contain Cobra, HTTP, REST, GraphQL, and persistence concerns at the edges.
- **Incremental extensibility:** Add capabilities through established boundaries rather than speculative frameworks.
- **Observability and safety:** Make targets, effects, and failures clear, especially for state-changing operations.

These principles refine the project values in `PROJECT_CONSTITUTION.md`; they do not replace them.

## 3. Layered Architecture

`fj` follows four conceptual layers:

1. **Domain:** Forgejo-related concepts, rules, and stable domain errors. It has no knowledge of Cobra, transport protocols, configuration formats, or output rendering.
2. **Application:** Use cases that coordinate domain behavior through explicit input and output boundaries. It defines the ports required from external systems.
3. **Interface:** CLI commands, input validation, presentation models, human-readable output, and AI-friendly JSON output. Cobra belongs to this layer.
4. **Infrastructure:** Forgejo API clients, authentication, configuration sources, and other external integrations. REST adapters belong here, with GraphQL adapters able to coexist later.

The layers are logical boundaries. Their value comes from dependency discipline and focused responsibilities, not from maximizing the number of packages.

## 4. Directory Structure

The intended high-level structure is:

```text
cmd/                 Executable entry points
internal/
  domain/            Domain concepts and rules
  application/       Use cases and external-system ports
  interface/         Cobra commands and output presentation
  infrastructure/    API, authentication, and configuration adapters
docs/                 Supporting design and user documentation
```

This structure is a direction for future implementation, not a requirement to create empty directories. New packages should be introduced only when an approved task needs them. Package boundaries should reflect responsibilities rather than mirror every conceptual type.

## 5. Core Components

- **Command interface:** Maps Cobra commands and flags to application use cases without containing domain policy.
- **Use-case layer:** Coordinates each user operation and remains independent of transport and presentation choices.
- **Forgejo service ports:** Express the capabilities that use cases require from a Forgejo instance.
- **Transport adapters:** Fulfill service ports through REST, with room for future GraphQL implementations.
- **Instance registry:** Resolves a selected Forgejo instance and its connection context for each operation.
- **Configuration boundary:** Supplies validated settings without exposing storage details to core layers.
- **Authentication boundary:** Provides credentials to transports while preventing their exposure in output or errors.
- **Output presenters:** Produce human-readable and stable, AI-friendly JSON representations from the same application results.
- **Error presentation:** Converts domain, application, and infrastructure failures into consistent messages, JSON errors, and process outcomes.

These components describe responsibilities only; their concrete APIs and implementation are defined by approved feature work.

## 6. Dependency Direction

Dependencies point inward:

```text
Interface ───────► Application ───────► Domain
Infrastructure ─► Application ports ─► Domain
```

The Domain layer depends on no outer layer. The Application layer may depend on Domain concepts and defines ports for required external capabilities. Interface and Infrastructure layers depend on those inner contracts and are composed at the application boundary.

Core behavior must not depend directly on Cobra, HTTP clients, REST or GraphQL schemas, configuration files, credential stores, or rendering formats. Data crossing a boundary is translated into the model owned by the receiving layer.

## 7. Error Handling

Errors are explicit outcomes and retain enough context to support diagnosis without exposing sensitive data.

- Domain errors describe violated rules or unavailable domain outcomes.
- Application errors describe use-case failures and preserve meaningful causes.
- Infrastructure errors translate transport, authentication, rate-limit, and remote-service failures at the adapter boundary.
- Interface presenters map failures to concise human messages, structured JSON errors, and stable exit behavior.

Error categories and machine-readable identifiers are compatibility-sensitive contracts. Messages should identify the failed operation and target while redacting credentials, tokens, and other sensitive values. Recoverable conditions should be distinguishable from invalid input and internal failures.

## 8. Configuration

Configuration is modeled independently of how it is stored. It supports named Forgejo instances, each with its own endpoint, authentication reference, and instance-specific preferences.

Instance selection must be explicit or derived from a documented, deterministic precedence order. Commands must make the effective target observable before consequential changes. Configuration validation occurs before a use case contacts a remote service.

Secrets are referenced through the configuration boundary but managed by an appropriate credential mechanism. They must not be embedded in ordinary configuration output, logs, diagnostics, or JSON responses. Changes to configuration shape require compatibility and migration consideration.

## 9. Testing Strategy

Testing follows architectural boundaries:

- **Domain tests** verify rules and invariants without external systems.
- **Application tests** verify use cases using controlled implementations of external ports.
- **Interface tests** verify command contracts, validation, human output, JSON schemas, and exit behavior.
- **Infrastructure tests** verify REST adapters, authentication, configuration, and protocol translation against controlled boundaries.
- **Integration tests** verify composition across layers and selected Forgejo-compatible API behavior.
- **Compatibility tests** protect stable commands, configuration, and machine-readable output from unintended breaking changes.

Tests should be deterministic and should not require live services unless explicitly designated. Each behavior change must be verified at the narrowest effective layer, with broader tests used where boundaries interact. Testing work must follow the execution and reporting rules in `AGENTS.md`.

## 10. Future Extensions

The architecture should allow the following additions without changing core policies:

- GraphQL transport adapters alongside REST adapters;
- support for new Forgejo resources and workflows through focused use cases;
- additional authentication and secure credential mechanisms;
- alternative configuration sources and instance discovery mechanisms;
- new presentation formats while preserving the JSON contract;
- non-interactive integrations and AI-agent tooling built on stable command and output boundaries;
- optional extension points whose ownership, compatibility, and security model are explicitly approved.

Future capabilities must enter through existing architectural boundaries or justify a reviewed architectural change. Extensibility does not permit unbounded plugin mechanisms, premature abstraction, or weakening human oversight. Architecture changes remain subject to the approval gates in `AGENTS.md`.
