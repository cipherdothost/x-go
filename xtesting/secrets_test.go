// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xtesting_test

import (
	"os"
	"testing"

	"go.cipher.host/x/xtesting"
)

func TestSetupSecrets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		giveDirectory string
		givePrefix    string
		giveSecrets   map[string]string
		wantError     bool
	}{
		{
			name:          "valid without secrets",
			giveDirectory: t.TempDir(),
			givePrefix:    "batman",
			giveSecrets:   make(map[string]string),
			wantError:     false,
		},
		{
			name:          "valid with secrets",
			giveDirectory: t.TempDir(),
			givePrefix:    "batman",
			giveSecrets: map[string]string{
				"SECRET":         "secret",
				"ANOTHER_SECRET": "another secret",
			},
			wantError: false,
		},
		{
			name:          "invalid directory",
			giveDirectory: "invalid",
			givePrefix:    "batman",
			giveSecrets: map[string]string{
				"SECRET": "secret",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := xtesting.SetupSecrets(t, tt.giveDirectory, tt.givePrefix, tt.giveSecrets)
			if (err != nil) != tt.wantError {
				t.Fatalf("SetupSecrets() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}

func TestSetupEnvironment(t *testing.T) { //nolint:paralleltest // SetupEnvironment uses t.Setenv, which isn't thread safe
	tmpDir := xtesting.SetupEnvironment(t)

	if tmpDir == "" {
		t.Fatal("SetupEnvironment() tmpDir is empty")
	}

	if os.Getenv("CREDENTIALS_DIRECTORY") != tmpDir {
		t.Fatalf("SetupEnvironment() CREDENTIALS_DIRECTORY = %q, want %q", os.Getenv("CREDENTIALS_DIRECTORY"), tmpDir)
	}

	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Fatalf("SetupEnvironment() tmpDir %q does not exist", tmpDir)
	}
}
