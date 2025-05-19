// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xargon2_test

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"go.cipher.host/x/xcrypto/xargon2"
)

func TestCompareHashAndPassword(t *testing.T) {
	t.Parallel()

	var (
		testParams   = xargon2.RecommendedParameters()
		testPepper   = []byte("somePepper")
		testPassword = "password123"
		testHash     = xargon2.GenerateFromPassword(testPassword, testParams, nil, testPepper)
	)

	tests := []struct {
		name         string
		giveHash     string
		givePassword string
		givePepper   []byte
		wantErr      error
	}{
		{
			name:         "Valid",
			giveHash:     testHash,
			givePassword: testPassword,
			givePepper:   testPepper,
			wantErr:      nil,
		},
		{
			name:         "InvalidPassword",
			giveHash:     testHash,
			givePassword: "wrongPassword",
			givePepper:   testPepper,
			wantErr:      xargon2.ErrInvalidPassword,
		},
		{
			name:         "InvalidHash",
			giveHash:     "invalid$hash",
			givePassword: testPassword,
			givePepper:   testPepper,
			wantErr:      xargon2.ErrInvalidHash,
		},
		{
			name:         "MalformedParameters",
			giveHash:     "$argon2id$v=19$m=bad,t=param,p=string$someSalt$someHash",
			givePassword: testPassword,
			givePepper:   testPepper,
			wantErr:      xargon2.ErrParseHashParameters,
		},
		{
			name:         "DecodeErrorSalt",
			giveHash:     "$argon2id$v=19$m=65536,t=3,p=4$!nvalidSalt$someHash",
			givePassword: testPassword,
			givePepper:   testPepper,
			wantErr:      xargon2.ErrDecode,
		},
		{
			name:         "DecodeErrorHash",
			giveHash:     "$argon2id$v=19$m=65536,t=3,p=4$someSalt$!nvalidHash",
			givePassword: testPassword,
			givePepper:   testPepper,
			wantErr:      xargon2.ErrDecode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := xargon2.CompareHashAndPassword(tt.giveHash, tt.givePassword, tt.givePepper)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateFromPassword(t *testing.T) {
	t.Parallel()

	testPepper := []byte("somePepper")

	tests := []struct {
		name           string
		givePassword   string
		giveParameters xargon2.Parameters
		givePepper     []byte
	}{
		{
			name:           "DefaultParameters",
			givePassword:   "password123",
			giveParameters: xargon2.MemoryConstrainedParameters(),
			givePepper:     testPepper,
		},
		{
			name:           "EmptyPassword",
			givePassword:   "",
			giveParameters: xargon2.RecommendedParameters(),
			givePepper:     testPepper,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash := xargon2.GenerateFromPassword(tt.givePassword, tt.giveParameters, nil, tt.givePepper)

			parts := strings.Split(hash, "$")
			if len(parts) != 6 {
				t.Errorf("invalid hash format: %s", hash)

				return
			}

			saltBase64, hashBase64 := parts[4], parts[5]

			_, err := base64.RawStdEncoding.DecodeString(saltBase64)
			if err != nil {
				t.Errorf("failed to decode salt: %v", err)
			}

			_, err = base64.RawStdEncoding.DecodeString(hashBase64)
			if err != nil {
				t.Errorf("failed to decode hash: %v", err)
			}
		})
	}
}
