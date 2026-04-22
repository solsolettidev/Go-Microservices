package main


type Config struct { // our application configuration (where we store the config from .env file)
	Mailer Mail

}

const webPort = "80"

func main() {
	app:= Config{
		Mailer: createMail(), // create the mailer instance
	}

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

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT")) // convert the port to an integer
	m := Mail{
		Domain: os.Getenv("MAIL_DOMAIN"),
		Host: os.Getenv("MAIL_HOST"),
		Port: port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		FromName: os.Getenv("FROM_NAME"),
	}
	
	return m
}