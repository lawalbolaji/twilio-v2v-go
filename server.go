package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"twilio-v2v/groq"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go/twiml"
)

type IVRPayload struct {
	SpeechText string `form:"SpeechText" json:"SpeechText" binding:"required"`
}

func (payload *IVRPayload) Validate() error {
	/* custom validation logic here */
	return nil
}

func main() {
	router := gin.Default()
	groq := groq.NewGroqClient(os.Getenv("GROQ_API_KEY"))

	router.POST("/voice/answer", func(ctx *gin.Context) {
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
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Header("Content-Type", "text/xml")
		ctx.String(http.StatusOK, twimlXml)
	})

	router.POST("/voice/ivr", func(ctx *gin.Context) {
		var form IVRPayload
		if err := ctx.ShouldBind(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, errors.New(err.Error()))
			return
		}

		/* business logic */
		llmResponse, err := groq.GetCompletion(form.SpeechText)
		if err != nil {
			fmt.Println(err)
		}

		/* twiml actions */
		answerUserQuestion := &twiml.VoiceSay{Message: llmResponse}
		collectNewUserQuestion := &twiml.VoiceGather{Action: "/voice/ivr", Input: "speech", SpeechTimeout: "2"}

		twimlXml, err := twiml.Voice([]twiml.Element{answerUserQuestion, collectNewUserQuestion})
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Header("Content-Type", "text/xml")
		ctx.String(http.StatusOK, twimlXml)
	})

	PORT := os.Getenv("TWILIO_V2V_PORT")
	if PORT == "" {
		PORT = "5515"
	}

	router.Run(fmt.Sprintf(":%s", PORT))
}
