// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xtesting

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go.cipher.host/x/xerrors"
)

// ErrWriteSecret is returned when writing a secret fails.
const ErrWriteSecret xerrors.Error = "failed to write secret"

// SetupSecrets writes the provided test secrets to files with the given prefix
// in the given directory.
func SetupSecrets(t *testing.T, dir, prefix string, secrets map[string]string) error {
	t.Helper()

	for k, v := range secrets {
		if err := os.WriteFile(filepath.Join(dir, prefix+"-"+k), []byte(v), 0o600); err != nil {
			return fmt.Errorf("%w %q: %w", ErrWriteSecret, k, err)
		}
	}

	return nil
}

// SetupEnvironment sets up the environment variables necessary for the
// [credential-go] package to work, and return the directory where the secrets
// are written.
//
// [credential-go]: https://sr.ht/~jamesponddotco/credential-go/
func SetupEnvironment(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	t.Setenv("CREDENTIALS_DIRECTORY", dir)

	return dir
}
