// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xtls_test

import (
	"crypto/tls"
	"testing"

	"go.cipher.host/x/xcrypto/xtls"
)

func TestModernConfig(t *testing.T) {
	t.Parallel()

	config := xtls.ModernConfig()

	if config.MinVersion != tls.VersionTLS13 {
		t.Errorf("ModernConfig(): unexpected MinVersion: got %v want %v", config.MinVersion, tls.VersionTLS13)
	}

	if config.SessionTicketsDisabled {
		t.Errorf("ModernConfig(): unexpected SessionTicketsDisabled: got %v want false", config.SessionTicketsDisabled)
	}

	if config.ClientSessionCache == nil {
		t.Errorf("ModernConfig(): ClientSessionCache was not initialized")
	}
}

func TestIntermediateConfig(t *testing.T) {
	t.Parallel()

	config := xtls.IntermediateConfig()

	if config.MinVersion != tls.VersionTLS12 {
		t.Errorf("IntermediateConfig(): unexpected MinVersion: got %v want %v", config.MinVersion, tls.VersionTLS12)
	}

	if config.SessionTicketsDisabled {
		t.Errorf("IntermediateConfig(): unexpected SessionTicketsDisabled: got %v want false", config.SessionTicketsDisabled)
	}

	if config.ClientSessionCache == nil {
		t.Errorf("IntermediateConfig(): ClientSessionCache was not initialized")
	}

	if len(config.CipherSuites) == 0 {
		t.Errorf("IntermediateConfig(): CipherSuites was not set")
	}
}

func TestDefaultCipherSuites(t *testing.T) {
	t.Parallel()

	cipherSuites := xtls.DefaultCipherSuites()

	if cipherSuites == nil {
		t.Error("cipherSuites is nil")
	}
}
