package middleware

import (
	"net/http"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(allowedOrigins, allowedMethods, allowedHeaders []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Set allowed methods
			if len(allowedMethods) > 0 {
				methods := ""
				for i, method := range allowedMethods {
					if i > 0 {
						methods += ","
					}
					methods += method
				}
				w.Header().Set("Access-Control-Allow-Methods", methods)
			}

			// Set allowed headers
			if len(allowedHeaders) > 0 {
				headers := ""
				for i, header := range allowedHeaders {
					if i > 0 {
						headers += ","
					}
					headers += header
				}
				w.Header().Set("Access-Control-Allow-Headers", headers)
			}

			// Handle preflight request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
