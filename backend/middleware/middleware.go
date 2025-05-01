package middleware

import (
	"log"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}
func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		wrapped := &wrappedWriter {
			ResponseWriter: w,
			statusCode: http.StatusOK,
		}
		next.ServeHTTP(wrapped,r)
		log.Println(
			"-----------------\n",
			"request:",r.Method, r.URL.Path, "\n",
			"response:",wrapped.statusCode, time.Since(start),
			"\n--------------------------------------",
		)

	})
}