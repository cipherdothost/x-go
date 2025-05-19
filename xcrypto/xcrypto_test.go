// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xcrypto_test

import (
	"testing"

	"go.cipher.host/x/xcrypto"
)

func TestZeroMemory(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		xcrypto.ZeroMemory(nil)
	})

	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()

		b := []byte{1, 2, 3, 4}
		xcrypto.ZeroMemory(b)

		if b[0] != 0 || b[1] != 0 || b[2] != 0 || b[3] != 0 {
			t.Fatalf("ZeroMemory() = %v, want empty slice", b)
		}
	})
}
