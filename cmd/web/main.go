package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"twilio-v2v/pkg/ai"
	"twilio-v2v/pkg/ai/groq"
)

type application struct {
	errorLog *log.Logger
	llm      ai.LLM
	dev      bool
	token    string
}

func main() {
	addr := flag.String("addr", ":5515", "HTTP network address")
	key := flag.String("key", "<GROQ_API_KEY>", "Groq API Access Key")
	dev := flag.Bool("dev", true, "DEV Flag; Request header will not be validated for twilio signature in dev mode")
	token := flag.String("token", "<TWILIO_AUTH_TOKEN>", "Twilio auth token")
	flag.Parse()

	ifLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errLog,
		llm: &groq.Groq{
			APIKey:  *key,
			BaseUrl: "https://api.groq.com/openai/v1",
		},
		dev:   *dev,
		token: *token,
	}
	server := &http.Server{
		Addr:     *addr,
		Handler:  app.registerRoutes(),
		ErrorLog: errLog,
	}

	ifLog.Printf("server listening on port: %s", *addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
