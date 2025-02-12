package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) registerRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/voice/answer", app.validateTwilioSignature(http.HandlerFunc(app.handleCallStarted)))
	r.Post("/voice/ivr", app.validateTwilioSignature(http.HandlerFunc(app.handleIVR)))

	return r
}
