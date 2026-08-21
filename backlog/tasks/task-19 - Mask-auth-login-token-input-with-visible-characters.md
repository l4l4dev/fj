---
id: TASK-19
title: Mask auth login token input with visible characters
status: Done
assignee:
  - '@l4l4dev'
created_date: '2026-08-21 12:38'
updated_date: '2026-08-21 12:42'
labels: []
dependencies: []
ordinal: 98000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
fj auth login reads the token via term.ReadPassword, which shows nothing while typing/pasting. Mika could not tell whether a paste succeeded. Show a mask character (e.g. ●) per byte typed/pasted so input is visible without revealing the token.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Interactive (non --token-stdin) token entry echoes one mask character per input byte instead of nothing
- [x] #2 Backspace/Delete removes the last mask character and the last buffered byte
- [x] #3 Enter finalizes input and the token is trimmed and passed through exactly as before
- [x] #4 Ctrl+C aborts without storing a partial token
- [x] #5 The actual token is never printed to the terminal
- [x] #6 Unit tests cover the masked-read byte-processing logic (append, backspace, enter, interrupt) without requiring a real TTY
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add readMaskedToken(reader *bufio.Reader, writer io.Writer) in auth.go: byte-loop, echo mask char per printable byte, backspace removes last byte+erases mask, CR/LF finalizes, Ctrl+C returns an interrupt error.
2. Wire it into readLoginToken: keep term.IsTerminal guard, call term.MakeRaw/Restore around the byte loop instead of term.ReadPassword, echo the prompt beforehand via command.PrintErrf as today.
3. Unit test readMaskedToken directly (no TTY needed) for: normal input, backspace, empty+enter, interrupt.
4. make pre-commit.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Replaced term.ReadPassword with a raw-mode byte loop (readMaskedToken) that echoes ● per printable byte; backspace erases, Enter finalizes, Ctrl+C aborts with errTokenEntryInterrupted. Trimming/behavior otherwise unchanged. make pre-commit green.

Manual verification by Mika: pasted a Playground token, saw 42 mask characters (matching a single clean paste, no duplication from the earlier plaintext attempts), pressed Enter, got 'Logged in to https://forgejo.yamaneco.io as Mika' and credentials.toml was written. make pre-commit green (gofmt/vet/test, including the 4 new readMaskedToken unit tests).
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Replaced the silent term.ReadPassword call in interactive auth login with readMaskedToken: a raw-mode byte loop that echoes ● per printable byte, erases on backspace/delete, finalizes on Enter, and aborts on Ctrl+C without touching the verifier/store. Token value, trimming, and downstream behavior are unchanged. Verified with 4 new unit tests (append/mask, backspace, empty-buffer backspace, Ctrl+C interrupt) plus a real login against Forgejo Playground by Mika (visible mask count matched a clean paste, login succeeded, credentials.toml written). make pre-commit green.
<!-- SECTION:FINAL_SUMMARY:END -->
