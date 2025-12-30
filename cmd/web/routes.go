package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fileserver := http.FileServer(http.Dir("./ui/static/"))

	r.Handle("/static/*", http.StripPrefix("/static", fileserver))
	r.Get("/v1/health", app.healthCheckHandler)
	r.Get("/", app.homeHandler)
	r.Get("/config/{name}", app.configFormPageHandler)
	r.Post("/v1/generate", app.configGeneratorHandler)
	r.Post("/download", app.downloadHandler)

	return r
}
