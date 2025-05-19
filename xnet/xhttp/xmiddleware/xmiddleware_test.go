// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.cipher.host/x/xnet/xhttp/xmiddleware"
)

func TestChain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		giveMiddlewares []func(http.Handler) http.Handler
		giveMethod      string
		wantStatusCode  int
	}{
		{
			name:            "No middlewares",
			giveMiddlewares: make([]func(http.Handler) http.Handler, 0),
			giveMethod:      http.MethodGet,
			wantStatusCode:  http.StatusOK,
		},
		{
			name: "Multiple middlewares",
			giveMiddlewares: []func(http.Handler) http.Handler{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Middleware", "1")

						next.ServeHTTP(w, r)
					})
				},

				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Middleware", "2")

						next.ServeHTTP(w, r)
					})
				},
			},
			giveMethod:     http.MethodGet,
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(tt.giveMethod, "/", http.NoBody)
			)

			var (
				handlerToChain = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				})
				chain = xmiddleware.Chain(handlerToChain, tt.giveMiddlewares...)
			)

			chain.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Chain(): expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}
