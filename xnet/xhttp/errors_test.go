// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xhttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.cipher.host/x/xnet/xhttp"
)

// failingWriter represents a mock http.ResponseWriter whose Write method always
// fails.
type failingWriter struct {
	header http.Header
}

func newFailingWriter() *failingWriter {
	return &failingWriter{
		header: make(http.Header),
	}
}

func (fw *failingWriter) Header() http.Header {
	return fw.header
}

func (*failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("simulated writer error")
}

func (*failingWriter) WriteHeader(_ int) {
}

func TestDetailError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give xhttp.DetailError
		want string
	}{
		{
			name: "all fields set",
			give: xhttp.DetailError{
				Title:         "Resource Not Found",
				Detail:        "The requested resource does not exist.",
				Documentation: "https://example.com/docs/errors#not-found",
				Status:        http.StatusNotFound,
			},
			want: "404: Resource Not Found: The requested resource does not exist.",
		},
		{
			name: "title empty, derived from status",
			give: xhttp.DetailError{
				Title:  "",
				Detail: "Something went wrong internally.",
				Status: http.StatusInternalServerError,
			},
			want: "500: Internal Server Error: Something went wrong internally.",
		},
		{
			name: "title empty, status 0, derived title is default",
			give: xhttp.DetailError{
				Title:  "",
				Detail: "Status was zero.",
				Status: 0,
			},
			want: "0: Unknown Error: Status was zero.",
		},
		{
			name: "title empty, status non-standard http code, derived title is default",
			give: xhttp.DetailError{
				Title:  "",
				Detail: "Non-standard code used.",
				Status: 601,
			},
			want: "601: Unknown Error: Non-standard code used.",
		},
		{
			name: "detail empty",
			give: xhttp.DetailError{
				Title:  "Payment Required",
				Detail: "",
				Status: http.StatusPaymentRequired,
			},
			want: "402: Payment Required",
		},
		{
			name: "title and detail empty",
			give: xhttp.DetailError{
				Title:  "",
				Detail: "",
				Status: http.StatusForbidden,
			},
			want: "403: Forbidden",
		},
		{
			name: "status only, title derived, detail empty",
			give: xhttp.DetailError{
				Status: http.StatusUnauthorized,
			},
			want: "401: Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.give.Error()
			if got != tt.want {
				t.Errorf("DetailError.Error(): expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestDetailError_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give xhttp.DetailError
		want string
	}{
		{
			name: "all fields set",
			give: xhttp.DetailError{
				Title:         "Resource Not Found",
				Detail:        "The requested resource does not exist.",
				Documentation: "https://example.com/docs/errors#not-found",
				Status:        http.StatusNotFound,
			},
			want: `{"title":"Resource Not Found","detail":"The requested resource does not exist.","documentation":"https://example.com/docs/errors#not-found","status":404}`,
		},
		{
			name: "title empty, derived from status",
			give: xhttp.DetailError{
				Title:         "",
				Detail:        "Internal issue.",
				Documentation: "https://example.com/docs/errors#internal",
				Status:        http.StatusInternalServerError,
			},
			want: `{"title":"Internal Server Error","detail":"Internal issue.","documentation":"https://example.com/docs/errors#internal","status":500}`,
		},
		{
			name: "title empty, status 0, derived title is default",
			give: xhttp.DetailError{
				Title:  "",
				Status: 0,
			},
			want: `{"title":"Unknown Error","status":0}`,
		},
		{
			name: "title empty, status non-standard, derived title is default",
			give: xhttp.DetailError{
				Title:  "",
				Status: 789,
			},
			want: `{"title":"Unknown Error","status":789}`,
		},
		{
			name: "detail and documentation empty, omitted",
			give: xhttp.DetailError{
				Title:         "Bad Request",
				Detail:        "",
				Documentation: "",
				Status:        http.StatusBadRequest,
			},
			want: `{"title":"Bad Request","status":400}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotBytes, err := tt.give.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() unexpected error: %v", err)
			}

			if string(gotBytes) != tt.want {
				t.Errorf("MarshalJSON(): expected %s, got %s", tt.want, string(gotBytes))
			}
		})
	}
}

func TestDetailError_WriteJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		give       xhttp.DetailError
		wantBody   string
		wantStatus int
	}{
		{
			name: "all fields set",
			give: xhttp.DetailError{
				Title:         "Resource Not Found",
				Detail:        "The requested resource does not exist.",
				Documentation: "https://example.com/docs/errors#not-found",
				Status:        http.StatusNotFound,
			},
			wantBody:   `{"title":"Resource Not Found","detail":"The requested resource does not exist.","documentation":"https://example.com/docs/errors#not-found","status":404}` + "\n",
			wantStatus: http.StatusNotFound,
		},
		{
			name: "successful write, title empty, derived from status",
			give: xhttp.DetailError{
				Title:  "",
				Detail: "An internal error occurred.",
				Status: http.StatusInternalServerError,
			},
			wantBody:   `{"title":"Internal Server Error","detail":"An internal error occurred.","status":500}` + "\n",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "successful write, title empty, status non-standard, derived title is default",
			give: xhttp.DetailError{
				Title:  "",
				Status: 999,
			},
			wantBody:   `{"title":"Unknown Error","status":999}` + "\n",
			wantStatus: 999,
		},
		{
			name: "successful write, detail and documentation empty, omitted",
			give: xhttp.DetailError{
				Title:         "Payment Required",
				Detail:        "",
				Documentation: "",
				Status:        http.StatusPaymentRequired,
			},
			wantBody:   `{"title":"Payment Required","status":402}` + "\n",
			wantStatus: http.StatusPaymentRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			if err := tt.give.WriteJSON(recorder); err != nil {
				t.Fatalf("WriteJSON() unexpected error: %v", err)
			}

			if recorder.Code != tt.wantStatus {
				t.Errorf("WriteJSON(): expected status %d, got %d", tt.wantStatus, recorder.Code)
			}

			if recorder.Header().Get("Content-Type") != xhttp.HeaderValueApplicationProblemJSON {
				t.Errorf("WriteJSON(): expected Content-Type %s, got %s", xhttp.HeaderValueApplicationProblemJSON, recorder.Header().Get("Content-Type"))
			}

			if recorder.Body.String() != tt.wantBody {
				t.Errorf("WriteJSON(): expected body %s, got %s", tt.wantBody, recorder.Body.String())
			}
		})
	}
}

func TestDetailError_WriteJSON_EncoderError(t *testing.T) {
	t.Parallel()

	detailErr := xhttp.DetailError{
		Status: http.StatusInternalServerError,
		Title:  "Internal Server Error",
		Detail: "This is a test error that should not be fully written.",
	}

	err := detailErr.WriteJSON(newFailingWriter())
	if err == nil {
		t.Fatal("WriteJSON() expected an error, but got nil")
	}

	if !errors.Is(err, xhttp.ErrWriteDetailError) {
		t.Errorf("WriteJSON() error = %v, want error to wrap %v", err, xhttp.ErrWriteDetailError)
	}
}
