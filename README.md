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

## Install release artifacts

Public releases are not available yet. Successful `Release foundation` workflow
runs provide a temporary `fj-<version>-artifacts` workflow artifact. Artifact
access and retention depend on the repository and its Actions settings.

The initial supported targets are:

| Platform | `uname -s` / `uname -m` | Binary |
| --- | --- | --- |
| macOS Apple Silicon | `Darwin` / `arm64` | `fj-<version>-darwin-arm64` |
| Linux x86-64 | `Linux` / `x86_64` | `fj-<version>-linux-amd64` |

Intel Macs, Linux arm64, Windows, and other targets are not currently
supported. On Linux, `x86_64` corresponds to the `amd64` artifact name.

### Obtain the workflow artifact

Open the repository's Actions page, select a successful `Release foundation`
run for the required version, and download `fj-<version>-artifacts` from its
Artifacts section. Extract the download into a directory of your choice. It
must contain the two binaries and checksum manifest:

```text
fj-<version>-checksums.txt
fj-<version>-darwin-arm64
fj-<version>-linux-amd64
```

If the GitHub CLI is already configured, the same artifact can be retrieved
without putting credentials in the command:

```bash
VERSION="1.2.3"
RUN_ID="<run-id>"

gh run download "$RUN_ID" \
  --name "fj-${VERSION}-artifacts" \
  --dir "fj-${VERSION}-artifacts"
```

`VERSION` is the normalized version without the tag's `v` prefix. The example
tag `v1.2.3` therefore produces filenames containing `1.2.3`.

### Verify the checksums

Run the verification command inside the directory containing all three files.
On macOS:

```bash
VERSION="1.2.3"
shasum -a 256 --check "fj-${VERSION}-checksums.txt"
```

On Linux:

```bash
VERSION="1.2.3"
sha256sum --check "fj-${VERSION}-checksums.txt"
```

Both commands must report `OK` for both binaries. Do not install or run the
artifacts if a file is missing or checksum verification fails. A matching
checksum verifies file integrity; it does not provide signing, notarization,
provenance, or proof of publisher identity.

### Install fj

After successful verification, install only the binary for the current
platform. On macOS Apple Silicon:

```bash
VERSION="1.2.3"
install -d "$HOME/.local/bin"
install -m 0755 \
  "fj-${VERSION}-darwin-arm64" \
  "$HOME/.local/bin/fj"
```

On Linux x86-64:

```bash
VERSION="1.2.3"
install -d "$HOME/.local/bin"
install -m 0755 \
  "fj-${VERSION}-linux-amd64" \
  "$HOME/.local/bin/fj"
```

These commands do not require `sudo`. Confirm the injected version directly,
then verify that the installed directory is on `PATH`:

```bash
"$HOME/.local/bin/fj" version
command -v fj
fj version
fj --help
```

The reported version must match `VERSION`. If `command -v fj` does not resolve
to `$HOME/.local/bin/fj`, add `$HOME/.local/bin` to the appropriate shell
configuration before using the remaining commands. Continue with
[Configuration onboarding](#configuration-onboarding) and the
[Quickstart](#quickstart); credential handling requirements are documented
there rather than repeated here.

### Uninstall

Confirm that `command -v fj` resolves to `$HOME/.local/bin/fj`, then remove only
the installed binary:

```bash
rm "$HOME/.local/bin/fj"
command -v fj
```

This does not remove `$HOME/.local/bin`, downloaded artifacts, configuration,
credentials, or shell settings. A different `fj` installation may still be
reported after removal.

These workflow artifacts are not permanent public distribution assets. The
initial binaries are not documented as signed or notarized, and checksum
success is not a reason to bypass operating-system security warnings. Homebrew,
package-manager, and system-wide installation are outside the current scope.

## Install from source on macOS

From a repository checkout, install `fj` for the current user without `sudo`:

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

## License

fj is available under the MIT License.
See [LICENSE](LICENSE) for details.
