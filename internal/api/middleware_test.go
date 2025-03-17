package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func newTestRouter(middlewares ...func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	r.Get("/healthz", healthCheckHandler)

	return r
}

func TestRecoverMiddleware_PanicRecovery(t *testing.T) {
	/* Arrange */
	router := newRouter()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)

	/* Act */
	router.ServeHTTP(recorder, request)

	/* Assert */
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.JSONEq(t, `{"error":"Internal Server Error"}`, recorder.Body.String())
}

func TestMiddlewareRequestID(t *testing.T) {
	/* Arrange */
	r := newRouter()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	/* Act */
	r.ServeHTTP(recorder, request)

	/* Assert */
	reqID := recorder.Header().Get("X-Request-id")
	assert.NotEmpty(t, reqID, "Request ID should be present in response headers")
}

func TestLoggerMiddlewareWithLoggerProvider(t *testing.T) {
	/* Arrange */
	var logBuffer bytes.Buffer

	// set a test log
	logProvider := func() zerolog.Logger {
		return zerolog.New(&logBuffer).With().Timestamp().Logger()
	}

	// set middlewares to test a target middleware
	lastMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			logger := log.Ctx(r.Context())
			logger.Info().Msg("test request_id")
			assert.NotNil(t, logger)
			w.WriteHeader(http.StatusOK)
		})
	}

	targetMiddleware := loggerMiddlewareWithLoggerProvider(logProvider)

	middlewares := []func(http.Handler) http.Handler{
		middleware.RequestID, targetMiddleware, lastMiddleware,
	}

	// set test responseWriter and request
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	// set a test router
	testRouter := newTestRouter(middlewares...)

	// set a map to read the test log content
	var logEntry map[string]any

	/* Act */
	testRouter.ServeHTTP(recorder, request)

	err := json.Unmarshal(logBuffer.Bytes(), &logEntry)

	/* Assert */
	t.Logf("log captured: %s", logBuffer.String())

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NoError(t, err)

	reqID, ok := logEntry["request_id"].(string)
	assert.Equal(t, recorder.Header().Get("x-request-id"), reqID)
	assert.True(t, ok, "request_id should be string representation of unique uuid")
	assert.Greater(t, len(reqID), 0)
}

func TestZerologRequestLoggerWithLoggerProvider(t *testing.T) {
	/* Arrange */
	var logBuffer bytes.Buffer

	// set a test log
	logProvider := func() zerolog.Logger {
		return zerolog.New(&logBuffer).With().Timestamp().Logger()
	}

	// set middlewares to test a target middleware
	targetMiddleware := zerologRequestLogger

	// set responseWriter and request
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	// set a test router
	testRouter := newTestRouter(
		middleware.RequestID,
		loggerMiddlewareWithLoggerProvider(logProvider),
		targetMiddleware,
	)

	// set a map to read the test log content
	var logEntry map[string]any

	/* Act */
	testRouter.ServeHTTP(recorder, request)
	err := json.Unmarshal(logBuffer.Bytes(), &logEntry)

	/* Assert */
	t.Logf("Canptured log: %s\n", logBuffer.String())

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, logEntry["method"])
	assert.Equal(t, "/healthz", logEntry["path"])
	assert.Equal(t, float64(http.StatusOK), logEntry["status"])

	requestID, ok := logEntry["request_id"].(string)
	assert.True(t, ok, "request_id should be string representation")
	assert.Greater(t, len(requestID), 0, "request_id is non-empty unique uuid")

	bytesWritten, ok := logEntry["bytes"].(float64)
	assert.True(t, ok, "bytes should be float64")
	assert.Greater(t, bytesWritten, float64(0), "bytesWritten should be greater than 0")

	duration, ok := logEntry["duration"].(float64)
	assert.True(t, ok, "duration should be a float64")
	assert.Greater(t, duration, float64(0), "duration should be greater than 0")
}
