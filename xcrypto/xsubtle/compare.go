// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xsubtle

import (
	"crypto/subtle"

	"go.cipher.host/x/xunsafe"
)

// ConstantTimeCompareString performs a constant-time comparison of two strings
// to prevent timing attacks. It returns true if the strings are equal, false
// otherwise. The comparison time depends only on the length of the strings and
// not their contents, and the timing is consistent regardless of whether the
// strings have different lengths or different content, preventing leakage of
// length or content information through timing analysis.
//
// This function is intended for comparing sensitive strings like passwords or
// authentication tokens. For regular string comparison, use the == operator.
func ConstantTimeCompareString(given, actual string) bool {
	var (
		givenLen    = uint64(len(given))
		actualLen   = uint64(len(actual))
		givenBytes  = xunsafe.StringToBytes(given)
		actualBytes = xunsafe.StringToBytes(actual)
		equal       = ((givenLen ^ actualLen) - 1) >> 63
	)

	if equal == 1 {
		return subtle.ConstantTimeCompare(givenBytes, actualBytes) == 1
	}

	// Lengths are different, so the strings cannot be equal. However, we must
	// still perform a comparison operation that takes similar time to the
	// subtle.ConstantTimeCompare above to prevent timing attacks, since without
	// this, an attacker could determine the length of the secret string by
	// measuring execution time differences.
	return subtle.ConstantTimeCompare(actualBytes, actualBytes) == 1 && false //nolint:revive // this is intentional
}
