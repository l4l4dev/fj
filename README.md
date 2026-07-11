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

## Install on macOS

Install `fj` for the current user without `sudo`:

```bash
make install
```

The binary is installed at `$HOME/.local/bin/fj`. The install target does not
modify `PATH` or shell configuration.

Check whether the install directory is already on `PATH`:

```bash
command -v fj
```

If no path is printed, add the directory to your shell configuration. For zsh:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

Confirm the installation:

```bash
fj --help
```

Remove only the installed binary with:

```bash
make uninstall
```

This does not remove configuration files, credentials, or shell settings.

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

## Configuration onboarding

The configuration file is loaded from the XDG config location:

- `$XDG_CONFIG_HOME/fj/config.toml` when `XDG_CONFIG_HOME` is set.
- `$HOME/.config/fj/config.toml` otherwise.

Define instance profiles with the following TOML schema. The endpoint below is
a placeholder; do not replace it with a real credential-bearing URL in
documentation.

```toml
[[instances]]
name = "playground"
endpoint = "https://forgejo-playground.example"
credential = "FORGEJO_PLAYGROUND_TOKEN"
```

The `credential` field contains only the name of an environment variable. Set
the credential value in the environment; never put the token itself in
`config.toml` or documentation:

```bash
export FORGEJO_PLAYGROUND_TOKEN="<token-not-shown>"
```

Use `--instance playground` to select the profile explicitly. With exactly one
configured profile, `--instance` may be omitted. When multiple profiles are
configured, an explicit `--instance` is required.

Credential values, raw tokens, and credentials embedded in URLs must not be
printed in command output, error messages, logs, or examples.

## Quickstart

After installing `fj` and configuring an instance profile, confirm the version
and start with read-only commands confirmed by the current acceptance checks:

```bash
fj --help
fj version
fj --version
fj repo inspect example-owner/example-repository --instance playground
fj issue list example-owner/example-repository --instance playground
fj issue inspect example-owner/example-repository NUMBER --instance playground
```

Replace `NUMBER` with a non-sensitive issue number available in the selected
instance. These examples do not create or modify repositories, issues, or pull
requests.

### Environment-dependent commands

The following commands depend on Forgejo Playground permissions and API
behavior. They are not presented as guaranteed successful examples:

```bash
fj repo list --instance playground
fj pr list example-owner/example-repository --instance playground
```

Treat failures from these commands as environment-dependent until the selected
instance's permissions and API compatibility have been verified. Do not record
credential values, raw tokens, real hostnames, or real repository owners in
command output or support reports.
