package main

import(
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct{
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"` //omitempty en Data significa "Si Data es nulo o vacío, ni siquiera incluyas esa llave al Frontend" (Ayudando a ahorrar ancho de banda).
} // this is a struct that is used to define the response that will be sent to the frontend

// vamos a querer tres funciones: para leer, escribir y generar error JSON  

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error{
	//limitation on the size of the json
	mayBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(mayBytes))

	dec := json.NewDecoder(r.Body) // creates a new decoder that will read from the request body
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{}) // this is used to check if there is any other data in the request body
	if err != io.EOF {
		return errors.New("Body must have only a single JSON value")
	}
	
	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error { // headers is an optional parameter
out, err := json.Marshal(data) // marshal is used to convert a Go object to a JSON object
if err != nil {
	return err
}

if len (headers)> 0 { // if there are any headers, add them to the response
	for key, value := range headers[0]{
		w.Header()[key] = value  // set the header
	}
}

w.Header().Set("Content-Type","application/json") // set the content type
w.WriteHeader(status) // set the status code
_,err = w.Write(out) // write the response
if err != nil {
	return err
}

return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int)error{ // status is an optional parameter
	statusCode := http.StatusBadRequest // set the default status code to 400

	if len (status) > 0 { // if there is a status code provided, use it
		statusCode = status[0]
	}

	var payload jsonResponse // create a json response
	payload.Error = true // set the error to true
	payload.Message = err.Error() // set the message to the error message

	return app.writeJSON(w, statusCode, payload) // write the response
}