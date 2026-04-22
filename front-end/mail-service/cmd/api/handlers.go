package main

import (
	"net/http"
	
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request){ // handler to send mail ()
	type mailMessage struct {
		From string `json:"from"`
		To string `json:"to"`
		Subject string `json:"subject"`
		Message any `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	msg := mail.Message{
		To:      requestPayload.To,
		From: requestPayload.From,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}
	
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "Mail sent successfully to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}