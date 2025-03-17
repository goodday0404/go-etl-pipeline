package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHandler(t *testing.T) {
	/* Arrange */
	router := newRouter()
	record := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	/* Act */
	router.ServeHTTP(record, request)

	/* Assert */
	assert.Equal(t, http.StatusOK, record.Code)
	assert.Equal(t, "Hello, Go-ETL-Pipeline", record.Body.String())
}

func TestHealthCheckHandler(t *testing.T) {
	/* Arrange */
	router := newRouter()
	record := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	/* Act */
	router.ServeHTTP(record, request)

	/* Assert */
	assert.Equal(t, http.StatusOK, record.Code)
	assert.JSONEq(t, `{"status": "ok"}`, record.Body.String())
}

func TestRecoveryCheckHandler(t *testing.T) {
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

func TestUnknownHandler(t *testing.T) {
	/* Arrange */
	record := httptest.NewRecorder()
	reqeust := httptest.NewRequest(http.MethodGet, "/unknown", nil)

	r := newRouter()
	r.ServeHTTP(record, reqeust)

	assert.Equal(t, http.StatusNotFound, record.Code)
}
