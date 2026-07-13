# Third-Party Notices

This document is the reviewed dependency-license inventory and attribution
index for the initial `fj` binary targets. The verbatim license texts are in
[`licenses/`](licenses/). This summary does not replace those license texts.

The inventory applies to binaries built from this repository for
`darwin/arm64` and `linux/amd64` with `CGO_ENABLED=0`. Before a Public Release,
the inventory must be checked against each actual release binary using
target-specific `go list -deps` output and `go version -m`.

## Linked production dependencies

| Module | Version | Classification | License | Copyright or attribution | Upstream NOTICE | Verbatim text |
| --- | --- | --- | --- | --- | --- | --- |
| `github.com/BurntSushi/toml` | `v1.6.0` | Direct / production / linked | MIT | Copyright (c) 2013 TOML authors | None found | [`BurntSushi-toml-COPYING`](licenses/BurntSushi-toml-COPYING) |
| `github.com/spf13/cobra` | `v1.10.2` | Direct / production / linked | Apache License 2.0 | Copyright 2013-2023 The Cobra Authors | None found | [`spf13-cobra-LICENSE.txt`](licenses/spf13-cobra-LICENSE.txt) |
| `github.com/spf13/pflag` | `v1.0.9` | Indirect / production / linked | BSD 3-Clause | Copyright (c) 2012 Alex Ogier; Copyright (c) 2012 The Go Authors | None found | [`spf13-pflag-LICENSE`](licenses/spf13-pflag-LICENSE) |

`github.com/BurntSushi/toml` includes `type_fields.go`, which states that its
struct-field handling is adapted from Go's `encoding/json` package and is
governed by the BSD-style license in the Go distribution. The corresponding
verbatim text is included as [`Go-LICENSE`](licenses/Go-LICENSE).

The Go standard library is classified as a linked toolchain/runtime component.
Its license text is provided in `Go-LICENSE` for the initial Go 1.26.5 release
toolchain.

## Modules not distributed in the initial binaries

The following modules appear in `go.mod`, `go.sum`, or the selected module
graph, but are not reachable from the production package graph for either
initial target. They are not represented as linked release dependencies.

| Module | Version | Classification | Reason not distributed |
| --- | --- | --- | --- |
| `github.com/inconshreveable/mousetrap` | `v1.1.0` | Indirect / Windows-specific / non-distributed | Cobra uses it only for Windows-specific behavior; the initial targets are macOS and Linux. |
| `github.com/cpuguy83/go-md2man/v2` | `v2.0.6` | Module-graph-only / non-distributed | Present through Cobra's module graph but not reachable from the fj production package graph. |
| `github.com/russross/blackfriday/v2` | `v2.1.0` | Module-graph-only / non-distributed | Present through go-md2man's module graph but not reachable from the fj production package graph. |
| `go.yaml.in/yaml/v3` | `v3.0.4` | Module-graph-only / non-distributed | Present through Cobra's module graph but not reachable from the fj production package graph. |
| `gopkg.in/check.v1` | `v0.0.0-20161208181325-20d25e280405` | Module-graph-only / non-distributed | Present through the YAML module graph but not reachable from the fj production package graph. |

## Other dependency and source classifications

- Test-only external dependencies: none.
- Build-tool dependencies: none.
- Vendored dependencies: none.
- Copied third-party source in the fj repository: none.
- Generated third-party source in the fj repository: none.
- Embedded third-party assets: none.

If the build targets, build tags, dependency versions, or reachable packages
change, this exclusion list and the included license texts must be reviewed
again. A future Windows binary would require a new review of `mousetrap` and
its Apache License 2.0 text.

## Redistribution summary

- The MIT notice for `BurntSushi/toml` must accompany copies or substantial
  portions of that software.
- Apache License 2.0 requires recipients of Cobra in source or object form to
  receive a copy of the license. Cobra `v1.10.2` has no upstream `NOTICE` file;
  this document is fj's own third-party notice summary and is not an upstream
  Apache `NOTICE` file.
- The pflag BSD 3-Clause license requires its copyright notice, conditions,
  and disclaimer to accompany binary redistribution and prohibits using named
  contributors to endorse derived products without permission.
- Static linking into a Go executable does not remove these redistribution
  conditions.

Unknown, incompatible, missing, or unresolved license or attribution
conditions block a Public Release. This inventory is an engineering compliance
record and is not legal advice.
