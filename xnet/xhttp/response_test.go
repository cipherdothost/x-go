// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xhttp_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"go.cipher.host/x/xnet/xhttp"
)

type errReader struct{}

func (*errReader) Read(_ []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

type customReadCloser struct {
	data *bytes.Buffer
}

func (c *customReadCloser) Read(p []byte) (n int, err error) {
	return c.data.Read(p)
}

func (*customReadCloser) Close() error {
	return errors.New("mock close error")
}

func TestDrainResponseBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveResponse *http.Response
		wantErr      bool
	}{
		{
			name: "valid response body",
			giveResponse: &http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte("valid response body"))),
			},
			wantErr: false,
		},
		{
			name: "empty response body",
			giveResponse: &http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte(""))),
			},
			wantErr: false,
		},
		{
			name: "error response body",
			giveResponse: &http.Response{
				Body: io.NopCloser(&errReader{}),
			},
			wantErr: true,
		},
		{
			name:         "nil response",
			giveResponse: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := xhttp.DrainResponseBody(tt.giveResponse)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestDrainResponseBody_ErrorClose(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: &customReadCloser{
			data: bytes.NewBufferString("test"),
		},
	}

	err := xhttp.DrainResponseBody(resp)
	if err == nil {
		t.Error("expected error, got nil")
	}

	want := fmt.Errorf("%w: %w", xhttp.ErrCannotCloseResponse, errors.New("mock close error"))
	if err.Error() != want.Error() {
		t.Errorf("got: %v, want: %v", err, want)
	}
}
