package main

import (
	"database/sql"
	"authentication/data"
	"log"
	"net/http"
	"time"
	"os"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"  
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgconn"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}
func main() {
	log.Println("Starting authentication service")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to Postgres")
	}

	// set up config
	app := Config{
		DB: conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string)(*sql.DB, error){
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready . . .")
			counts++

		}else{
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {  // if we have tried to connect 10 times and failed, we will exit
			log.Println(err)
			return nil
		}
		
		log.Println("Backing off for 2 seconds . . .") // wait 2 seconds before trying again
		time.Sleep(2 * time.Second)
		continue
	}
}