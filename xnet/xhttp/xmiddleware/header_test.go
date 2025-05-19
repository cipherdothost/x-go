// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"go.cipher.host/x/xcontext"
	"go.cipher.host/x/xnet/xhttp"
	"go.cipher.host/x/xnet/xhttp/xmiddleware"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		giveUserAgent string
		wantStatus    int
		wantBody      string
	}{
		{
			name:          "missing User-Agent",
			giveUserAgent: "",
			wantStatus:    http.StatusBadRequest,
			wantBody: `{"title":"Bad Request","detail":"User agent is missing. Please provide a valid user agent.","status":400}
`,
		},
		{
			name:          "valid User-Agent",
			giveUserAgent: "TestAgent",
			wantStatus:    http.StatusOK,
			wantBody:      "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			r.Header.Set("User-Agent", tt.giveUserAgent)

			middleware := xmiddleware.UserAgent(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write([]byte("OK")) //nolint:errcheck // safe to ignore, we're testing
			}))
			middleware.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				t.Errorf("UserAgent(): expected status %d, got %d", tt.wantStatus, w.Code)
			}

			body := w.Body.String()
			if body != tt.wantBody {
				t.Errorf("UserAgent(): expected body %s, got %s", tt.wantBody, body)
			}
		})
	}
}

func TestPrivacyPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		giveURI string
		wantURI string
	}{
		{
			name:    "valid URI",
			giveURI: "https://example.com/privacy",
			wantURI: "https://example.com/privacy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			middleware := xmiddleware.PrivacyPolicy(tt.giveURI, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
			middleware.ServeHTTP(w, r)

			if w.Header().Get("Privacy-Policy") != tt.wantURI {
				t.Errorf("PrivacyPolicy(): expected header %s, got %s", tt.wantURI, w.Header().Get("Privacy-Policy"))
			}
		})
	}
}

func TestTermsOfService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		giveURI string
		wantURI string
	}{
		{
			name:    "valid URI",
			giveURI: "https://example.com/terms",
			wantURI: "https://example.com/terms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			middleware := xmiddleware.TermsOfService(tt.giveURI, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
			middleware.ServeHTTP(w, r)

			if w.Header().Get("Terms-Of-Service") != tt.wantURI {
				t.Errorf("TermsOfService(): expected header %s, got %s", tt.wantURI, w.Header().Get("Terms-Of-Service"))
			}
		})
	}
}

func TestRequestID(t *testing.T) {
	t.Parallel()

	const contextKeyRequestID xcontext.ContextKey = "request_id"

	tests := []struct {
		name           string
		giveExistingID string
		validateOutput func(*testing.T, string)
	}{
		{
			name:           "with existing request ID",
			giveExistingID: "test-request-id",
			validateOutput: func(t *testing.T, id string) {
				t.Helper()

				if id != "test-request-id" {
					t.Errorf("RequestID(): expected request ID %q, got %q", "test-request-id", id)
				}
			},
		},
		{
			name:           "without existing request ID",
			giveExistingID: "",
			validateOutput: func(t *testing.T, id string) {
				t.Helper()

				if len(id) != 36 {
					t.Errorf("RequestID(): expected UUID length of 36, got %d", len(id))
				}

				matched := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(id)
				if !matched {
					t.Errorf("RequestID(): invalid UUID format: %s", id)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			if tt.giveExistingID != "" {
				r.Header.Set(xhttp.HeaderKeyRequestID, tt.giveExistingID)
			}

			var contextID string

			middleware := xmiddleware.RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				id, ok := r.Context().Value(contextKeyRequestID).(string)
				if !ok {
					t.Error("RequestID(): request ID not found in context")
				}

				contextID = id
			}))

			middleware.ServeHTTP(w, r)

			headerID := w.Header().Get(xhttp.HeaderKeyRequestID)

			tt.validateOutput(t, headerID)

			if headerID != contextID {
				t.Errorf("RequestID(): header ID %q does not match context ID %q", headerID, contextID)
			}
		})
	}
}
