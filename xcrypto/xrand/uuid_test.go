// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xrand_test

import (
	"regexp"
	"testing"

	"go.cipher.host/x/xcrypto/xrand"
)

func TestUUID_Generate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "Valid UUIDv4",
			want: "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				u   = &xrand.UUID{}
				got = u.GenerateV4()
			)

			matched, err := regexp.MatchString(tt.want, got)
			if err != nil {
				t.Fatalf("failed to match UUID regex: %v", err)
			}

			if !matched {
				t.Errorf("Generate() = %q, want match %q", got, tt.want)
			}
		})
	}
}
