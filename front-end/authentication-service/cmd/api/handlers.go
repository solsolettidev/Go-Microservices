package main

import(
	"errors"
	"net/http"
	"fmt"
	"bytes"
	"encoding/json"
)

func(app *Config) Authenticate(w http.ResponseWriter, r *http.Request){
	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate the user against the database
	user, err:= app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	
	valid, err:= user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("User %s logged in", user.Email)) // this is how we log requests to the logger service
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) logRequest(name, data string) error { // this function is used to log requests to the logger service
	var entry struct{
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData)) // this is how we send data to the logger service, bytes.NewBuffer is used to convert the json data to a byte buffer
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		return errors.New("error calling logger service")
	}

	return nil
}