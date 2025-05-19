// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xtls

import (
	"crypto/tls"
)

// ModernConfig returns a [*tls.Config] using the modern profile from the
// [Mozilla SSL Configuration Generator].

//
// [*tls.Config]: https://pkg.go.dev/crypto/tls#Config
// [Mozilla SSL Configuration Generator]: https://ssl-config.mozilla.org/#server=go&version=1.23.3&config=modern&guideline=5.7
func ModernConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.X25519MLKEM768,
		},
		SessionTicketsDisabled: false,
		ClientSessionCache:     tls.NewLRUClientSessionCache(128),
	}
}

// IntermediateConfig returns a [*tls.Config] using the intermediate profile
// from the [Mozilla SSL Configuration Generator].
//
// [*tls.Config]: https://pkg.go.dev/crypto/tls#Config
// [Mozilla SSL Configuration Generator]: https://ssl-config.mozilla.org/#server=go&version=1.23.3&config=intermediate&guideline=5.7
func IntermediateConfig() *tls.Config {
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		CipherSuites: DefaultCipherSuites(),
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.X25519MLKEM768,
		},
		SessionTicketsDisabled: false,
		ClientSessionCache:     tls.NewLRUClientSessionCache(128),
	}
}

// DefaultCipherSuites returns a sensible default list of cipher suites for
// TLS1.2 based on [Mozilla's recommendations].
//
// [Mozilla's recommendations]: https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
func DefaultCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
}
