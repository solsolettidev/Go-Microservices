package main
import(
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)
func (app *Config) routes() http.Handler { // func that returns an http.handler

	mux := chi.NewRouter() // variable that is the comun name when using routing in go, its the abreviation of Multiplexor (it takes
	// multiple entry signals and redirect them where their belongs
	

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{ // cors is used to allow other origins to connect to our api
		AllowedOrigins: []string{"http://*","http://*"},//the allowed ones r everyone for dev purposes
	    AllowedMethods: []string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders: []string{"Accept","Authorization","Content-Type","X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	mux.Use(middleware.Heartbeat("/ping")) 
	// this is a healthcheck endpoint, it is used to check if the server is alive, if it is it will return a 200 status code (enterying the url /ping)

	mux.Post("/send", app.SendMail) // path to send mail that goes to the handlers.
	return mux
}
