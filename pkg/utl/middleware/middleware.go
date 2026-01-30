package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/urfave/negroni"
)

const (
	// HeaderXRequestID is the header name for the request ID
	HeaderXRequestID = "X-Request-ID"
)

// RequestID middleware generates a unique request ID for each incoming request.
// It sets the X-Request-ID header in the response. If the client sends an
// X-Request-ID header, that value is used instead of generating a new UUID.
func RequestID() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		requestID := r.Header.Get(HeaderXRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		rw.Header().Set(HeaderXRequestID, requestID)

		// Optionally, you might want to store the requestID in the request context
		// for access by other parts of the application (e.g., loggers).
		// ctx := context.WithValue(r.Context(), contextKeyRequestID, requestID)
		// r = r.WithContext(ctx)

		next(rw, r)
	}
}
