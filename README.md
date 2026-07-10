# fj

AI-first CLI for Forgejo

## Status

🚧 Under development

## Command behavior

- `fj` and `fj --help` print root help to standard output and exit successfully.
- `-h` and `--help` are the supported global help flags.
- Successful command output is written to standard output.
- Invalid input and other failures are written to standard error and return a non-zero exit status without repeating usage text.

## Instance selection

Instance selection uses this precedence:

1. An explicitly requested profile.
2. The only configured profile when no profile is explicitly requested.

Selection fails when an explicitly requested profile does not exist or when multiple profiles are configured without an explicit selection.
