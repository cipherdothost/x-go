// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xrand_test

import (
	"errors"
	"testing"

	"go.cipher.host/x/xcrypto/xrand"
)

func TestBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give int
		want int
	}{
		{
			name: "Default give (32 bytes)",
			give: 0,
			want: 32,
		},
		{
			name: "8 random bytes",
			give: 8,
			want: 8,
		},
		{
			name: "16 random bytes",
			give: 16,
			want: 16,
		},
		{
			name: "64 random bytes",
			give: 64,
			want: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xrand.Bytes(tt.give)
			if len(got) != tt.want {
				t.Errorf("Bytes(): got %d bytes, want %d bytes", len(got), tt.want)
			}
		})
	}
}

func TestBytesWithReader_Panic(t *testing.T) {
	t.Parallel()

	reader := errorReader{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected BytesWithReader() to panic, but it didn't")
		}
	}()

	_ = xrand.BytesWithReader(16, reader)
}

type errorReader struct{}

func (errorReader) Read(_ []byte) (int, error) {
	return 0, errors.New("forced error")
}
