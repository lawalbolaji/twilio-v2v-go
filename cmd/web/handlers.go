package main

import (
	"net/http"
	"twilio-v2v/pkg/forms"

	"github.com/twilio/twilio-go/twiml"
)

func (app *application) handleCallStarted(w http.ResponseWriter, r *http.Request) {
	/* TODO: the details of what to say to the user should be in a service not in this handler */
	intro := &twiml.VoiceSay{
		Message: "hello... you have reached the twilio voice assistant... how can I help you today?",
	}

	collectUserInput := &twiml.VoiceGather{
		Input:         "speech",
		Action:        "/voice/ivr",
		SpeechTimeout: "2",
	}

	endCall := &twiml.VoiceSay{
		Message: "You can give us a call back any time",
	}

	twimlXml, err := twiml.Voice([]twiml.Element{intro, collectUserInput, endCall})
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(twimlXml))
}

func (app *application) handleIVR(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	ok := form.Required("SpeechText") /* we can chain the validation steps */
	if !ok {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	llmResponse, err := app.llm.GetCompletion(form.Get("SpeechText"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	answerUserQuestion := &twiml.VoiceSay{Message: llmResponse}
	collectNewUserQuestion := &twiml.VoiceGather{Action: "/voice/ivr", Input: "speech", SpeechTimeout: "2"}
	twimlXml, err := twiml.Voice([]twiml.Element{answerUserQuestion, collectNewUserQuestion})

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(twimlXml))
}
