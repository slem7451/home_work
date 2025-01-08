package internalhttp

import (
	"log"
	"net/http"
	"os"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	if err := os.MkdirAll("logs/", 0o666); err != nil {
		panic(err)
	}
	
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(err)
	}

	logger := log.New(logFile, "", 0)
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		logger.Printf("%s [%s] %s %s %s %d %d \"%s\"\n", r.RemoteAddr, start.String(), r.Method, r.URL.String(), r.Proto, lrw.statusCode, time.Since(start), r.UserAgent())
	})
}
