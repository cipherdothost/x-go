// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

// Package xcrypto provides functions and utilities that extend [Go's standard
// crypto package].
//
// [Go's standard crypto package]: https://pkg.go.dev/crypto
package xcrypto

// ZeroMemory sets all bytes in a byte slice (up to its capacity) to zero,
// erasing its contents. It does not provide guarantees against memory analysis
// or raw memory dumping by operators, but it does minimize the exposure.
func ZeroMemory(b []byte) {
	if b == nil {
		return
	}

	b = b[:cap(b):cap(b)]

	for i := range b {
		b[i] = 0
	}
}
