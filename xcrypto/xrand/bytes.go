// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xrand

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
)

// Bytes returns an amount of cryptographically-secure random bytes equal to the
// the specified size.
//
// The bytes are generated using the [crypto/rand] package and will panic if
// crypto/rand fails to generate the requested amount of bytes.
//
// Defaults to 32 if no size is specified.
//
// [crypto/rand]: https://pkg.go.dev/crypto/rand
func Bytes(size int) []byte {
	return BytesWithReader(size, rand.Reader)
}

// BytesWithReader is like Bytes, but allows you to specify your own reader.
func BytesWithReader(size int, reader io.Reader) []byte {
	if size < 1 {
		size = 32
	}

	var (
		pool   = make([]byte, 0, size)
		buffer = bytes.NewBuffer(pool)
	)

	_, err := io.CopyN(buffer, reader, int64(size))
	if err != nil {
		panic(
			fmt.Errorf("%w", err),
		)
	}

	return buffer.Bytes()
}
