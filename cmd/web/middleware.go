package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/twilio/twilio-go/client"
)

func (app *application) validateTwilioSignature(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* turn off middleware validation in dev mode  */
		if app.dev {
			next.ServeHTTP(w, r)
			return
		}

		requestValidator := client.NewRequestValidator(app.token)
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Restore the request body so the next handler can read it
		r.Body = io.NopCloser(bytes.NewReader(body))
		url := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
		signature := r.Header.Get("X-Twilio-Signature")
		ok := requestValidator.ValidateBody(url, []byte(body), signature)
		if !ok {
			app.clientError(w, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
