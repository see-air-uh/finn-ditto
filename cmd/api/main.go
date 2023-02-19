package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/see-air-uh/finn-ditto/data"
	"github.com/see-air-uh/finn-ditto/token"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// const (
// 	webPort  = "50001"
// 	mongoURL = "mongodb://localhost:27017"
// )

// var counts int64

type Config struct {
	DB           *sql.DB
	Models       data.Models
	M_Model      data.M_Model
	PasetoClient token.GoTokens
	WebPort string
}

var client *mongo.Client

// TODO: Get this key from MONGO DB
var RANDOM_STRING = "Qcv4I4HV9161U6RiaqOggFDmTuQAl6DJ"

func main() {
		// check if there is an environment variable already set to production
	environment := os.Getenv("ENVIRONMENT")

	if environment == "" {
		// TODO: Load .env file
		log.Println("Loading environment variables from local .env file")
		godotenv.Load(".env")
	}
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

	t, err := token.NewPasetoClient(RANDOM_STRING)
	if err != nil {
		panic(err)
	}

	app := setupApp(client, t)

	// start the grpc server
	app.gRPCListen()

}

func connectToMongo() (*mongo.Client, error) {
	dsn := os.Getenv("MONGO_URL")
	user := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	// create connection options
	clientOptions := options.Client().ApplyURI(dsn)
	clientOptions.SetAuth(options.Credential{
		Username: user,
		Password: password,
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

func setupApp( m *mongo.Client,t  token.GoTokens) *Config {
	return &Config{
		M_Model:      data.NewMongo(client),
		PasetoClient: t,
		WebPort: os.Getenv("WEB_PORT"),
	}
}
