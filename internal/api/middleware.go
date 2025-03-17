package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerProvider func() zerolog.Logger

func recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			instance := recover()

			if instance == nil {
				return
			}

			var err error

			switch val := instance.(type) {
			case string:
				err = fmt.Errorf("%s", val)
			case error:
				err = val
			default:
				err = fmt.Errorf("unknown panic: %v", val)
			}

			reqID := middleware.GetReqID(r.Context())
			stackTrace := string(debug.Stack())

			logger := log.Ctx(r.Context())

			logger.Error().
				Caller().
				Str("request_id", reqID).
				Err(err).
				Str("stack_trace", stackTrace).
				Msg("Recovered from panic")

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "applcation/json")
			w.Write([]byte(`{"error": "Internal Server Error"}`))
		}()

		next.ServeHTTP(w, r)
	})
}

// setLogger injects a request-specific logger with request_id into context.
func loggerMiddleware(next http.Handler) http.Handler {
	return loggerMiddlewareWithLoggerProvider(func() zerolog.Logger {
		return log.Logger
	})(next)
}

func loggerMiddlewareWithLoggerProvider(provider LoggerProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve or generate request ID
			reqID := middleware.GetReqID(r.Context())

			w.Header().Set("x-request-id", reqID)

			// Derive a logger with request_id (inherits global logger's settings)
			requestLogger := provider().With().Str("request_id", reqID).Logger()

			// Inject logger into context
			ctx := requestLogger.WithContext(r.Context())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// zerologRequestLogger is a custom request logger using zerolog.
func zerologRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		// retrieve the request logger set in loggerMiddlewareWithLoggerProvider
		logger := log.Ctx(r.Context())

		logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.Status()).
			Int("bytes", ww.BytesWritten()).
			Dur("duration", time.Since(start)).
			Str("remote", r.RemoteAddr).
			Msg("HTTP request")
	})
}
