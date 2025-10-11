package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		ip := r.RemoteAddr
		timeReq := time.Now().Format("02/Jan/2006:15:04:05 -0700")
		method := r.Method
		url := r.URL.Path
		proto := r.Proto
		ua := r.Header.Get("User-Agent")
		if ua == "" {
			ua = "-"
		}

		msg := fmt.Sprintf(
			"%s - [%s] \"%s %s %s\" %d %s \"%s\"",
			ip,
			timeReq,
			method,
			url,
			proto,
			lrw.statusCode,
			duration,
			ua,
		)

		h.logger.Info(msg)
	})
}
