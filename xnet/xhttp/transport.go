// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xhttp

import (
	"crypto/tls"
	"net/http"
	"runtime"
	"time"

	"go.cipher.host/x/xcrypto/xtls"
	"go.cipher.host/x/xnet"
)

const (
	// DefaultTransportMaxIddleConns is the default maximum number of idle
	// connections in the pool.
	DefaultTransportMaxIddleConns int = 100

	// DefaultTransportIdleConnTimeout is the default maximum amount of time a
	// connection can remain idle before being closed.
	DefaultTransportIdleConnTimeout = 90 * time.Second

	// DefaultTransportTLSHandshakeTimeout is the default maximum amount of time
	// allowed to complete a TLS handshake.
	DefaultTransportTLSHandshakeTimeout = 10 * time.Second

	// DefaultTransportExpectContinueTimeout is the default maximum amount of
	// time to wait for a server's first response headers after fully writing
	// the request headers if the request has an "Expect: 100-continue" header.
	DefaultTransportExpectContinueTimeout = 1 * time.Second
)

// NewTransport returns a new http.Transport with similar default values to
// http.DefaultTransport, but is not shared, thus cannot be altered by other
// packages.
//
// If tlsConfig is nil, then [xtls.ModernConfig()] is used.
//
// [xtls.ModernConfig()]: https://pkg.go.dev/go.cipher.host/x/xcrypto/xtls#ModernConfig
func NewTransport(tlsConfig *tls.Config) *http.Transport {
	if tlsConfig == nil {
		tlsConfig = xtls.ModernConfig()
	}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           xnet.NewDialer().DialContext,
		TLSClientConfig:       tlsConfig,
		MaxIdleConns:          DefaultTransportMaxIddleConns,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		IdleConnTimeout:       DefaultTransportIdleConnTimeout,
		TLSHandshakeTimeout:   DefaultTransportTLSHandshakeTimeout,
		ExpectContinueTimeout: DefaultTransportExpectContinueTimeout,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
	}
}
