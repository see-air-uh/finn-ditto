package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/see-air-uh/asxce-toga/data"
)

const webPort = "50001"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	log.Println("Attempting to start authentication service...")

	conn := connectToDB()

	// check for failed connection
	if conn == nil {
		log.Panic("Couldn't connect to Postgres!")
	}

	// set up config var
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// start the grpc server
	app.gRPCListen()

}

func connectToDB() *sql.DB {

	dsn := os.Getenv("DSN")

	//loop until connection to DB is made
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready to connect...")
			counts++
		} else {
			log.Println("Connected to postgres")
			return connection
		}

		// handle case can't connect to db
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off to 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {

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
