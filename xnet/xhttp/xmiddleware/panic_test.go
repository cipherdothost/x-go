// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xmiddleware_test

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.cipher.host/x/xnet/xhttp/xmiddleware"
)

// failingResponseWriter represents an http.ResponseWriter whose Write method
// always fails.
type failingResponseWriter struct {
	header http.Header
	code   int
}

func newFailingResponseWriter() *failingResponseWriter {
	return &failingResponseWriter{
		header: make(http.Header),
	}
}

func (fw *failingResponseWriter) Header() http.Header {
	return fw.header
}

func (*failingResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("simulated writer error")
}

func (fw *failingResponseWriter) WriteHeader(statusCode int) {
	if fw.code == 0 {
		fw.code = statusCode
	}
}

func TestPanicRecovery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		handler    http.Handler
		wantStatus int
		wantBody   string
	}{
		{
			name: "with panic",
			handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				panic("forced panic")
			}),
			wantStatus: http.StatusInternalServerError,
			wantBody: `{"title":"Internal Server Error","detail":"Something went wrong. Please try again later.","status":500}
`,
		},
		{
			name: "without panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK")) //nolint:errcheck // safe to ignore, we're testing
			}),
			wantStatus: http.StatusOK,
			wantBody:   "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				buf     bytes.Buffer
				handler = slog.NewJSONHandler(&buf, nil)
				logger  = slog.New(handler)
			)

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			middleware := xmiddleware.PanicRecovery(logger, tt.handler)
			middleware.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				t.Errorf("PanicRecovery(): expected status %d, got %d", tt.wantStatus, w.Code)
			}

			body := w.Body.String()
			if body != tt.wantBody {
				t.Errorf("PanicRecovery(): expected body %s, got %s", tt.wantBody, body)
			}
		})
	}
}

func TestPanicRecovery_WriteJSONError(t *testing.T) {
	t.Parallel()

	var (
		buf     bytes.Buffer
		handler = slog.NewJSONHandler(&buf, nil)
		logger  = slog.New(handler)
	)

	panickingHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("forced panic to trigger WriteJSON error path")
	})

	var (
		writer     = newFailingResponseWriter()
		recorder   = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		middleware = xmiddleware.PanicRecovery(logger, panickingHandler)
	)

	middleware.ServeHTTP(writer, recorder)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "panic recovered") {
		t.Errorf("PanicRecovery(): expected log output to contain 'panic recovered', got %s", logOutput)
	}

	if !strings.Contains(logOutput, "failed to write response") {
		t.Errorf("PanicRecovery(): expected log output to contain 'failed to write response', got %s", logOutput)
	}

	if !strings.Contains(logOutput, "simulated writer error") {
		t.Errorf("PanicRecovery(): expected log output to contain 'simulated writer error', got %s", logOutput)
	}
}
