// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xsubtle_test

import (
	"strings"
	"testing"

	"go.cipher.host/x/xcrypto/xsubtle"
)

func TestConstantTimeCompareString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		giveStringA string
		giveStringB string
		want        bool
	}{
		{
			name:        "equal strings",
			giveStringA: "secret",
			giveStringB: "secret",
			want:        true,
		},
		{
			name:        "different strings same length",
			giveStringA: "secret",
			giveStringB: "necret",
			want:        false,
		},
		{
			name:        "different lengths",
			giveStringA: "secret",
			giveStringB: "secrets",
			want:        false,
		},
		{
			name:        "very large strings equal",
			giveStringA: strings.Repeat("a", 1<<30),
			giveStringB: strings.Repeat("a", 1<<30),
			want:        true,
		},
		{
			name:        "very large strings different",
			giveStringA: strings.Repeat("a", 1<<30),
			giveStringB: strings.Repeat("b", 1<<30),
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xsubtle.ConstantTimeCompareString(tt.giveStringA, tt.giveStringB)
			if got != tt.want {
				t.Errorf("ConstantTimeStringEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
