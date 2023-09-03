package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func IP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buffer     []byte
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		slog.Info("",
			"url", r.RequestURI,
			"status-code", lw.statusCode,
			"method", r.Method,
			"duration", duration.String(),
			"client-ip", IP(r),
			"user-agent", r.UserAgent(),
		)
	}
	return http.HandlerFunc(handler)
}
