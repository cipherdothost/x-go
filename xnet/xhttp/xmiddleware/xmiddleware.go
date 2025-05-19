// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

// Package xmiddleware contains simple middleware functions.
package xmiddleware

import "net/http"

// Chain wraps a given http.Handler with multiple middleware functions.
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
