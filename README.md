<!--
SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>

SPDX-License-Identifier: CC0-1.0
-->

# x

[![Go Documentation](https://pkg.go.dev/badge/v3.svg)](https://pkg.go.dev/go.cipher.host/x)
[![Go Report Card](https://goreportcard.com/badge/go.cipher.host/x)](https://goreportcard.com/report/go.cipher.host/x)
[![Tests](https://github.com/cipherdothost/x-go/actions/workflows/ci.yml/badge.svg)](https://github.com/cipherdothost/x-go/actions/workflows/ci.yml)

Package **x** is a collection of extensions to Go's standard library
that are shared across multiple packages, services, and applications at
Cipher Host. The naming scheme follows the `x[pkgname]` pattern.

If we're using the same code across two or more projects, it'll either
end up here or as its own package.

> [!WARNING]
> These packages are considered internal to Cipher Host and are provided
> without any stability guarantees for external usage. We recommend
> copying the code instead of importing the packages.

## Contributing

Anyone can help make these packages better. Check out [the contribution
guidelines](CONTRIBUTING.md) for more information.

---

The work in this repository complies with [the REUSE
specification](https://reuse.software/spec-3.3/). While [the default
license is MIT](LICENSE.md), individual files may be licensed
differently.

Please see the individual files for details and [the LICENSES
directory](LICENSES/) for a full list of licenses used in this
repository.
