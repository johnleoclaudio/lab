package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery middleware recovers from panics
func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					reqID := GetRequestID(r.Context())
					logger.ErrorContext(r.Context(), "panic recovered",
						slog.String("request_id", reqID),
						slog.Any("error", err),
						slog.String("stack", string(debug.Stack())),
					)

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error"}]}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
