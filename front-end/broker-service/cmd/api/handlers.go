package main

import (
	"net/http"
	"encoding/json"
	"errors"
	"bytes"
)

type RequestPayload struct{
	Action string `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request){
//w (ResponseWriter): Es tu "lapicera" para escribirle en la cara la respuesta final al Frontend.
//r (Request): Es una radiografía de todo lo que envió el Frontend (incluye cabeceras, qué navegador es, qué IP, y el body/cuerpo json si mandó algún formulario).
	payload := jsonResponse {
		Error: false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK,payload)
}

/*
handlers.go - La lógica de negocio (Controladores)
Una vez el enrutador dejó pasar al Frontend y lo mandó para acá, este código define qué hacer.

type jsonResponse struct {...}: Con este simple paquete defines el "Contrato" de comunicación con el Frontend. Pase lo que pase, el Frontend siempre sabrá que va a recibir un JSON con esa misma forma.
Los anotadores entre comillas invertidas como `json:"message"` le avisan a la librería de Go llamada encoding/json cómo debe de llamar esa llave cuando se pase a formato universal JSON para la web. omitempty en Data significa "Si Data es nulo o vacío, ni siquiera incluyas esa llave al Frontend" (Ayudando a ahorrar ancho de banda).
func (app *Config) Broker(w http.ResponseWriter, r *http.Request):
w (ResponseWriter): Es tu "lapicera" para escribirle en la cara la respuesta final al Frontend.
r (Request): Es una radiografía de todo lo que envió el Frontend (incluye cabeceras, qué navegador es, qué IP, y el body/cuerpo json si mandó algún formulario).
La ejecución:
Instancias el payload ("la caja que enviarás").
json.MarshalIndent(...): Toma tu objeto puro de Go (el payload), y lo serializa transformándolo en puros bytes que lucen como un JSON bonito y con sangrías \t.
w.Header().Set("Content-Type","application/json"): Es importantísimo. Le indica al navegador explícitamente: "Cuidado Frontend, lo que te mando a continuación no es texto plano, ni es una foto, es un JSON. Trátalo como tal".
w.WriteHeader(http.StatusAccepted): Manda un código de estado número 202. Es perfecto para un Broker, porque el 202 significa "He aceptado tu petición, pero aún la voy a procesar a futuro" en lugar de un clásico 200 de "Todo hecho".
Finaliza escribiendo el resultado serializado out dentro de w.
*/


func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request){
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.Authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("invalid action"))
		return
	}
}

func (app *Config) Authenticate(w http.ResponseWriter, a AuthPayload){
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "","\t")
	// call the service

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get bsck the correct status code
	if response.StatusCode != http.StatusUnauthorized{
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}else if response.StatusCode != http.StatusAccepted{
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}
	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message))
		return
	}

	// everything is good, send it back to the frontend
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)

}