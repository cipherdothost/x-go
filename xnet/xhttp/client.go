// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xhttp

import (
	"net/http"
	"time"

	"go.cipher.host/x/xcrypto/xtls"
)

// Doer represents the ability to execute HTTP requests. It abstracts HTTP
// client implementations to make testing easier and enable flexible client injection.
//
// The standard http.Client implements this interface.
type Doer interface {
	// Do executes the given HTTP request and returns the response. The caller
	// is responsible for closing the response body when done.
	Do(req *http.Request) (*http.Response, error)
}

// DefaultClientTimeout is the default timeout for the http.Client.
const DefaultClientTimeout = 15 * time.Second

// NewModernClient returns a new http.Client given the provided timeout.
//
// Unlike Go's http.DefaultClient:
// - It is not shared, thus cannot be altered by other modules.
// - It uses [xtls.ModernConfig()] for TLS configuration.
// - It doesn't follow redirects.
// - It doesn't accept cookies.
//
// A zero timeout means the default timeout is used.
//
// [xtls.ModernConfig()]: https://pkg.go.dev/go.cipher.host/x/xcrypto/xtls#ModernConfig
func NewModernClient(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = DefaultClientTimeout
	}

	return &http.Client{
		Transport: NewTransport(nil),
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     nil,
		Timeout: timeout,
	}
}

// NewIntermediateClient is like NewModernClient, but uses
// [xtls.IntermediateConfig()] for TLS configuration.
//
// [xtls.IntermediateConfig()]: https://pkg.go.dev/go.cipher.host/x/xcrypto/xtls#IntermediateConfig
func NewIntermediateClient(timeout time.Duration) *http.Client {
	client := NewModernClient(timeout)
	client.Transport = NewTransport(xtls.IntermediateConfig())

	return client
}
