// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xmiddleware

import (
	"context"
	"net/http"

	"go.cipher.host/x/xcontext"
	"go.cipher.host/x/xcrypto/xrand"
	"go.cipher.host/x/xnet/xhttp"
)

// contextKeyRequestID is the key used for the request ID in the context.
const contextKeyRequestID xcontext.ContextKey = "request_id"

// UserAgent ensures that the request has the User-Agent header set.
func UserAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() == "" {
			response := xhttp.DetailError{
				Status: http.StatusBadRequest,
				Detail: "User agent is missing. Please provide a valid user agent.",
			}

			response.WriteJSON(w) //nolint:errcheck // if this fails, there isn't much we can do

			return
		}

		next.ServeHTTP(w, r)
	})
}

// PrivacyPolicy adds a privacy policy header to the response.
func PrivacyPolicy(uri string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Privacy-Policy", uri)

		next.ServeHTTP(w, r)
	})
}

// TermsOfService adds a terms of service header to the response.
func TermsOfService(uri string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Terms-Of-Service", uri)

		next.ServeHTTP(w, r)
	})
}

// RequestID adds a request ID header to the response.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(xhttp.HeaderKeyRequestID)
		if requestID == "" {
			var uuid xrand.UUID

			requestID = uuid.GenerateV4()
		}

		ctx := context.WithValue(r.Context(), contextKeyRequestID, requestID)

		r = r.WithContext(ctx)

		w.Header().Set(xhttp.HeaderKeyRequestID, requestID)

		next.ServeHTTP(w, r)
	})
}
