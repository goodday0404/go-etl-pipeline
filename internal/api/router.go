package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer) // recover panic from chi RequestID middleware
	r.Use(middleware.RequestID)
	r.Use(recovererMiddleware) // recover panic from all user defined middlewares
	r.Use(loggerMiddleware)
	r.Use(zerologRequestLogger)

	// register handlers here
	r.Get("/", defaultHandler)
	r.Get("/healthz", healthCheckHandler)
	r.Get("/panic", recoveryCheckHandler) // only for test from panic

	return r
}
