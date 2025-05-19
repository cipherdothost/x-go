// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xmiddleware

import (
	"context"
	"log/slog"
	"net/http"

	"go.cipher.host/x/xnet/xhttp"
)

// PanicRecovery tries to recover from panics and returns a 500 error if there
// was one.
func PanicRecovery(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		defer func(ctx context.Context) {
			if err := recover(); err != nil {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"panic recovered",
					slog.Any("error", err),
				)

				response := xhttp.DetailError{
					Status: http.StatusInternalServerError,
					Detail: "Something went wrong. Please try again later.",
				}

				if err = response.WriteJSON(w); err != nil {
					logger.LogAttrs(
						ctx,
						slog.LevelError,
						"failed to write response",
						slog.Any("error", err),
					)
				}
			}
		}(ctx)

		next.ServeHTTP(w, r)
	})
}
