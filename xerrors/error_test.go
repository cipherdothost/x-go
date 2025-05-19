// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xerrors_test

import (
	"errors"
	"testing"

	"go.cipher.host/x/xerrors"
)

func TestError(t *testing.T) {
	t.Parallel()

	const (
		errGeneric xerrors.Error = "generic error"
	)

	tests := []struct {
		name       string
		give       error
		wantString string
		wantError  error
	}{
		{
			name:       "generic error",
			give:       errGeneric,
			wantString: "generic error",
			wantError:  errGeneric,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.give.Error(); got != tt.wantString {
				t.Errorf("Error(): got %v, want %v", got, tt.wantString)
			}

			if !errors.Is(tt.give, tt.wantError) {
				t.Errorf("Error(): got %v, want %v", tt.give, tt.wantError)
			}
		})
	}
}
