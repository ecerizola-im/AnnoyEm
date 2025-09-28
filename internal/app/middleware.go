package app

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

// generateReqID creates a short random hex string.
func generateReqID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		// fallback to timestamp-based ID if randomness fails
		return hex.EncodeToString([]byte(time.Now().Format("150405")))
	}
	return hex.EncodeToString(b)
}

// LoggingMiddleware logs basic request info with a request ID.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := generateReqID()
		start := time.Now()
		log.Printf("[req %s] -> %s %s CL=%d from %s", id, r.Method, r.URL.Path, r.ContentLength, r.RemoteAddr)

		// Wrap ResponseWriter to capture status code
		lrw := &loggingResponseWriter{ResponseWriter: w, status: 200}
		r = r.WithContext(r.Context())
		next.ServeHTTP(lrw, r)

		dur := time.Since(start)
		log.Printf("[req %s] <- %s %s %d %s", id, r.Method, r.URL.Path, lrw.status, dur)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}
