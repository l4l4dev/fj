# fj

AI-first CLI for Forgejo

## Status

🚧 Under development

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for the repository's contribution
standards and required review workflow.

## Development

The shared development commands are provided by the repository `Makefile`:

- `make fmt` formats Go sources.
- `make check-fmt` checks that Go sources are formatted.
- `make vet` runs static analysis.
- `make test` runs the Go test suite.
- `make build` builds all Go packages.
- `make verify` runs formatting, whitespace, vet, and test checks.
- `make pre-commit` runs the complete pre-commit verification.

## Command behavior

- `fj` and `fj --help` print root help to standard output and exit successfully.
- `-h` and `--help` are the supported global help flags.
- Successful command output is written to standard output.
- Invalid input and other failures are written to standard error without repeating
  usage text.
- Process outcomes distinguish validation, authentication, remote-service, and
  internal failures. Numeric exit-code values remain an internal implementation
  detail rather than a published compatibility contract.
- Error messages identify the failed operation and use category-safe text rather
  than exposing authentication, remote-service, or internal causes.

## Instance selection

Instance selection uses this precedence:

1. An explicitly requested profile.
2. The only configured profile when no profile is explicitly requested.

Selection fails when an explicitly requested profile does not exist or when multiple profiles are configured without an explicit selection.
