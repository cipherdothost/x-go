// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.cipher.host/x/xnet/xhttp"
)

func TestNewModernClient_CheckRedirect(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://example.com", http.StatusFound)
	}))
	defer ts.Close()

	client := xhttp.NewModernClient(0)

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, ts.URL, http.NoBody)
	if err != nil {
		t.Fatalf("NewClient(): expected no error, got %v", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		t.Fatalf("NewClient(): expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("NewClient(): expected status code %d, got %d", http.StatusFound, resp.StatusCode)
	}
}

func TestNewIntermediateClient_CheckRedirect(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://example.com", http.StatusFound)
	}))
	defer ts.Close()

	client := xhttp.NewIntermediateClient(0)

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, ts.URL, http.NoBody)
	if err != nil {
		t.Fatalf("NewClient(): expected no error, got %v", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		t.Fatalf("NewClient(): expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("NewClient(): expected status code %d, got %d", http.StatusFound, resp.StatusCode)
	}
}
