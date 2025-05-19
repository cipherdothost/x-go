// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xrand

import (
	"crypto/rand"
	"errors"
	"testing"
)

type mockReader struct{}

func (*mockReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("mock read error")
}

func TestUUID_Generate_Panic(t *testing.T) { //nolint:paralleltest // we're modifying rand.Reader
	origReader := rand.Reader

	t.Cleanup(func() {
		rand.Reader = origReader
	})

	rand.Reader = &mockReader{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected GenerateV4 to panic, but it did not")
		}
	}()

	uuid := &UUID{}
	_ = uuid.GenerateV4()
}
