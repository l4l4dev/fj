# Project Constitution

## 1. Purpose of This Project

`fj` is a project that provides an AI-first CLI for safe and transparent collaboration between people who use Forgejo and AI.

Rather than merely invoking the Forgejo API from a terminal, it aims to understand development context—including repositories, Issues, Pull Requests, and reviews—and support day-to-day software development. This support must not replace human decision-making; it must enable people to make better decisions with less friction.

## 2. The CLI We Aim to Build

`fj` aims to be a CLI with the following characteristics:

- It provides a consistent command-line interface for development activities on Forgejo.
- It provides input and output that are easy for both people and AI to understand and compose.
- It clearly communicates state and scope of impact before and after execution, allowing users to operate it with confidence.
- It supports both interactive use and automation, with operations that have the same meaning behaving predictably.
- It respects Forgejo and its users, without creating unnecessary dependence on a particular AI service or closed environment.

AI-first does not mean granting AI unlimited authority. It means designing explicit, structured interfaces so that AI can accurately interpret the meaning, results, and failures of commands and provide assistance safely under human supervision.

## 3. Design Principles

This project is designed according to the following principles:

1. **Prefer explicitness.** Clearly identify inputs, targets, changes, and results without relying on implicit behavior or assumptions.
2. **Choose safe defaults.** Prevent destructive operations and external changes from being performed unintentionally, and make them subject to confirmation when appropriate.
3. **Keep components small and composable.** Give each feature a focused responsibility and make it reusable from interactive sessions, scripts, and AI agents.
4. **Balance machine readability with human readability.** Provide structured output without losing the information people need for diagnosis and auditing.
5. **Be deterministic and predictable.** Under the same conditions, operations should have the same meaning and outcome whenever possible, and failures must be expressed clearly.
6. **Choose dependencies carefully.** Do not casually add dependencies that compromise maintainability, security, or portability.
7. **Keep the Public API minimal.** Preserve room for future change and expose only what can be maintained over time as a contract with users.
8. **Ensure observability.** Design the system so that users can determine what was done, why it was done, and which target it affected.

We prioritize understandability, verifiability, and ease of change over short-term implementation speed. At the same time, we avoid speculative over-abstraction and evolve incrementally based on actual requirements.

## 4. Rules for Collaborative Development with AI

AI is treated as a collaborator that assists with research, design proposals, implementation, testing, review, and documentation. However, AI output is always a proposal or work product and does not, by itself, constitute evidence of correctness.

- AI must review existing documentation and approved designs before starting work.
- AI must strictly respect the requested scope, prohibitions, and approval gates.
- Design and implementation must remain separate, and AI must wait for explicit human approval at stages where approval is required.
- AI must not present uncertain information as fact and must clearly identify assumptions, inferences, and unverified information.
- Changes must remain small, reviewable, and traceable.
- AI must explain not only the code, but also the design rationale, scope of impact, verification results, and remaining issues.
- AI-generated results must not be adopted without testing or review.
- AI must not read, record, or transmit credentials, personal information, or non-public information when doing so is unnecessary.
- AI-driven automation must be introduced with minimal privileges and in a form that people can stop, reject, or correct.

## 5. Humans Make the Final Decisions

Humans hold final decision-making authority and responsibility for the project.

AI may present options, rationale, risks, and trade-offs, but it must not unilaterally finalize specifications, approve designs, perform destructive operations, issue releases, make security judgments, or make decisions about community governance. When the information required for a decision is insufficient, AI must ask a human for clarification.

Humans must verify AI proposals and reject or revise them when necessary. The use of AI does not remove responsibility for accountability or review.

## 6. Documentation as the Source of Truth

Approved specifications, designs, operational policies, and decisions are recorded in repository documentation, which serves as the Source of Truth. Chat history, AI memory, verbal agreements, and inferences from implementation alone are not authoritative sources.

When code and documentation conflict, the discrepancy must be made explicit and synchronized after a human determines which is correct. Important decisions must be recorded with their background and rationale so that future developers can understand not only the conclusion, but also the context in which it was reached.

This constitution defines the project's highest-level policies. Any proposed change to it must be handled separately from ordinary implementation changes and undergo explicit human review and approval.

## 7. Designed for Long-Term Operation

`fj` is to be developed as software intended for continued use and maintenance, not as a temporary prototype.

- Consider backward compatibility and migration paths, and avoid breaking users unnecessarily.
- Treat tests, documentation, and change history as maintained project assets.
- Design boundaries and responsibilities that can adapt to changes in Forgejo and its surrounding environment.
- Continuously address security issues and dependency updates.
- Reduce knowledge and procedures that only a particular individual can understand.
- For deprecations and incompatible changes, communicate the rationale, impact, and migration path.

We assess project health not only by the number of new features, but also by the reliability of existing functionality, maintenance burden, and impact on users.

## 8. Values as an OSS Project

`fj` respects Forgejo's philosophy and community and aims to be open source software that creates public value.

- **Transparency:** Make designs, constraints, reasons for changes, and known issues public wherever possible.
- **Accessibility for contributors:** Provide documentation that first-time contributors can understand and foster respectful communication.
- **User autonomy:** Respect users' control over their data, credentials, execution environment, and use of AI.
- **Interoperability:** Respect open standards and the Forgejo API, and avoid unnecessary lock-in.
- **Sustainability:** Treat maintainers' time and capacity as finite resources and choose development processes that can continue without unreasonable burden.
- **Inclusivity and respect:** Respect users and contributors from different backgrounds and with different levels of experience, and encourage constructive, safe collaboration.
- **Integrity:** Communicate capabilities and limitations accurately, and do not overstate quality or safety.

The project's success is measured not only by the breadth of its features, but also by whether users trust it, understand it, can participate in improving it, and can continue using it over the long term.
