// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xos_test

import (
	"testing"
	"time"

	"go.cipher.host/x/xos"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		giveKey      string
		giveFallback string
		giveEnvValue string
		want         string
	}{
		{
			name:         "non-existent variable with a fallback",
			giveKey:      "SOMETHING_THAT_DOES_NOT_EXIST",
			giveFallback: "fallback",
			giveEnvValue: "",
			want:         "fallback",
		},
		{
			name:         "existent variable with a value",
			giveKey:      "SOME_EXISTENT_VARIABLE",
			giveFallback: "fallback",
			giveEnvValue: "test-value",
			want:         "test-value",
		},
		{
			name:         "existent variable with an empty value",
			giveKey:      "SOME_EXISTENT_VARIABLE",
			giveFallback: "fallback",
			giveEnvValue: "",
			want:         "fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.giveKey, tt.giveEnvValue)

			got := xos.GetEnv(tt.giveKey, tt.giveFallback)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntEnv(t *testing.T) {
	tests := []struct {
		name         string
		giveKey      string
		giveFallback int
		giveEnvValue string
		want         int
	}{
		{
			name:         "non-existent variable with a fallback",
			giveKey:      "SOMETHING_THAT_DOES_NOT_EXIST",
			giveFallback: 123,
			giveEnvValue: "",
			want:         123,
		},
		{
			name:         "existent variable with a valid integer value",
			giveKey:      "SOME_EXISTENT_VARIABLE",
			giveFallback: 123,
			giveEnvValue: "456",
			want:         456,
		},
		{
			name:         "existent variable with an invalid integer value",
			giveKey:      "SOME_EXISTENT_VARIABLE",
			giveFallback: 123,
			giveEnvValue: "not-an-integer",
			want:         123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.giveKey, tt.giveEnvValue)

			got := xos.GetIntEnv(tt.giveKey, tt.giveFallback)
			if got != tt.want {
				t.Errorf("GetIntEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		giveKey      string
		giveFallback bool
		giveEnvValue string
		want         bool
	}{
		{
			name:         "non-existent variable with a fallback",
			giveKey:      "SOMETHING_THAT_DOES_NOT_EXIST",
			giveFallback: true,
			giveEnvValue: "",
			want:         true,
		},
		{
			name:         "existent variable with a valid boolean value",
			giveKey:      "SOME_EXISTENT_VARIABLE",
			giveFallback: true,
			giveEnvValue: "false",
			want:         false,
		},
		{
			name:         "existent variable with an invalid boolean value",
			giveKey:      "SOME_EXISTENT_VARIABLE",
			giveFallback: true,
			giveEnvValue: "not-a-boolean",
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.giveKey, tt.giveEnvValue)

			got := xos.GetBoolEnv(tt.giveKey, tt.giveFallback)
			if got != tt.want {
				t.Errorf("GetBoolEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDurationEnv(t *testing.T) {
	tests := []struct {
		name         string
		giveKey      string
		giveFallback time.Duration
		giveEnvValue string
		want         time.Duration
	}{
		{
			name:         "non-existent variable with a fallback",
			giveKey:      "SOMETHING_THAT_DOES_NOT_EXIST",
			giveFallback: 5 * time.Minute,
			giveEnvValue: "",
			want:         5 * time.Minute,
		},
		{
			name:         "existent variable with a valid duration value",
			giveKey:      "SOME_EXISTENT_VARIABLE",
			giveFallback: 5 * time.Minute,
			giveEnvValue: "1h30m",
			want:         1*time.Hour + 30*time.Minute,
		},
		{
			name:         "existent variable with an invalid duration value",
			giveKey:      "SOME_EXISTENT_VARIABLE",
			giveFallback: 5 * time.Minute,
			giveEnvValue: "not-a-duration",
			want:         5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.giveKey, tt.giveEnvValue)

			got := xos.GetDurationEnv(tt.giveKey, tt.giveFallback)
			if got != tt.want {
				t.Errorf("GetDurationEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
