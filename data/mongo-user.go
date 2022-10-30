package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func NewMongo(mongo *mongo.Client) M_Model {
	client = mongo

	return M_Model{
		M_User: M_User{},
	}
}

type M_Model struct {
	M_User M_User
}

type M_User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string    `bson:"email" json:"email"`
	Username  string    `bson:"username" json:"username"`
	FirstName string    `bson:"first_name,omitempty" json:"first_name"`
	LastName  string    `bson:"last_name,omitempty" json:"last_name"`
	Password  string    `bson:"-" json:"-"`
	Active    bool      `bson:"active" json:"active"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

func (c *M_User) CreateUser(arg_user M_User) error {
	collection := client.Database("users").Collection("users")

	_, err := collection.InsertOne(context.TODO(), M_User{
		Email:     arg_user.Email,
		Username:  arg_user.Username,
		FirstName: arg_user.FirstName,
		LastName:  arg_user.LastName,
		Password:  arg_user.Password,
		Active:    arg_user.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting user to the database: ", err)
		return err
	}

	return nil
}

func (c *M_User) GetUserByEmail(arg_email string) (*M_User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("users").Collection("users")

	var user M_User

	err := collection.FindOne(ctx, bson.M{"email": arg_email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}

func (c *M_User) GetUserByUsername(arg_username string) (*M_User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("users").Collection("users")

	var user M_User
	log.Println("USERNAME-->", arg_username)
	err := collection.FindOne(ctx, bson.M{"username": arg_username}).Decode(&user)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
