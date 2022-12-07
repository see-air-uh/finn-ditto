package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/see-air-uh/finn-ditto/data"
	"github.com/see-air-uh/finn-ditto/token"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort  = "50001"
	mongoURL = "mongodb://localhost:27017"
)

var counts int64

type Config struct {
	DB           *sql.DB
	Models       data.Models
	M_Model      data.M_Model
	PasetoClient token.GoTokens
}

var client *mongo.Client

// TODO: Get this key from MONGO DB
var RANDOMSTRING = "Qcv4I4HV9161U6RiaqOggFDmTuQAl6DJ"

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	// log.Panic("Error reading .env file...")

	// }
	log.Println("Attempting to start authentication service...")

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// conn := connectToDB()

	// // check for failed connection
	// if conn == nil {
	// 	log.Panic("Couldn't connect to database!")
	// }

	t, err := token.NewPasetoClient(RANDOMSTRING)
	if err != nil {
		panic(err)
	}

	// set up config var
	app := Config{
		// DB:           conn,
		// Models:       data.New(conn),
		M_Model:      data.NewMongo(client),
		PasetoClient: t,
	}

	// start the grpc server
	app.gRPCListen()

}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "admin12345",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}
	log.Println("Connected to mongo.")

	return c, nil
}

// func connectToDB() *sql.DB {

// 	// dsn := os.Getenv("DSN")
// 	// log.Println("LOGGING OUT DSN >>" + dsn + "<<<")
// 	// //loop until connection to DB is made
// 	// for {
// 	// 	connection, err := openDB(dsn)
// 	// 	if err != nil {
// 	// 		log.Println("Postgres not yet ready to connect...")
// 	// 		log.Println(err)
// 	// 		counts++
// 	// 	} else {
// 	// 		log.Println("Connected to postgres")
// 	// 		return connection
// 	// 	}

// 	// 	// handle case can't connect to db
// 	// 	if counts > 10 {
// 	// 		log.Println(err)
// 	// 		return nil
// 	// 	}
// 	// 	log.Println("Backing off to 2 seconds...")
// 	// 	time.Sleep(2 * time.Second)
// 	// 	continue
// 	// }
// }

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
