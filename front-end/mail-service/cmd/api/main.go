package main


type Config struct {

}

const webPort = "80"

func main() {
	app:= Config{}

	log.Println("Starting mail service on port", webPort)

	// now we define a server that will listen on the port we defined above
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// now we start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}