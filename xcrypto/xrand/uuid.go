// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xrand

import (
	"crypto/rand"
	"io"

	"go.cipher.host/x/xunsafe"
)

const (
	// _hexDigits is a lookup table for converting a byte into a hex string.
	hexDigits string = "0123456789abcdef"

	// _separator is the separator used in the UUID string.
	separator byte = '-'
)

// UUID represents a 128-bit Universal Unique Identifier (UUID) as defined in
// [RFC 4122].
//
// [RFC 4122]: https://tools.ietf.org/html/rfc4122
type UUID struct {
	// bytes is a byte slice of length 16 that is used as a buffer for the UUID.
	bytes [16]byte
}

// Generate generates a UUIDv4 and returns it as a string.
func (u *UUID) GenerateV4() string {
	if _, err := io.ReadFull(rand.Reader, u.bytes[:]); err != nil {
		panic(err)
	}

	u.bytes[6] = (u.bytes[6] & 0x0f) | 0x40
	u.bytes[8] = (u.bytes[8] & 0x3f) | 0x80

	buf := make([]byte, 0, 36)

	for i := range 16 {
		if i == 4 || i == 6 || i == 8 || i == 10 {
			buf = append(buf, separator)
		}

		buf = append(buf, hexDigits[u.bytes[i]>>4], hexDigits[u.bytes[i]&0x0f])
	}

	return xunsafe.BytesToString(buf)
}
