// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xtime_test

import (
	"encoding/json"
	"testing"
	"time"

	"go.cipher.host/x/xtime"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    string
		want    xtime.Duration
		wantErr bool
	}{
		{
			name:    "valid duration",
			give:    `"1h"`,
			want:    xtime.Duration(time.Hour),
			wantErr: false,
		},
		{
			name:    "valid zero duration",
			give:    `"0"`,
			want:    xtime.Duration(0),
			wantErr: false,
		},
		{
			name:    "invalid duration",
			give:    `"invalid"`,
			wantErr: true,
		},
		{
			name:    "non-string duration",
			give:    `1`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got xtime.Duration

			err := json.Unmarshal([]byte(tt.give), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && got != tt.want {
				t.Fatalf("UnmarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
