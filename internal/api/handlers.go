package api

import (
	"fmt"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello, Go-ETL-Pipeline")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func recoveryCheckHandler(w http.ResponseWriter, r *http.Request) {
	// This endpoint is only to test internal/api/recovererMiddlewarerecovererMiddleware
	panic("Panic to test internal/api/recovererMiddleware")
}
