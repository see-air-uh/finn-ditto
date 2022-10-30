package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/see-air-uh/finn-ditto/auth"
	"github.com/see-air-uh/finn-ditto/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Models  data.Models
	M_Model data.M_Model
}

type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AuthServer) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	input := req.GetArgUser()

	//check if username or email exists

	_, err := a.M_Model.M_User.GetUserByEmail(input.Email)
	if err == nil {
		return nil, status.Errorf(401, "error. email in use")
	}

	_, err = a.M_Model.M_User.GetUserByUsername(input.Username)
	if err == nil {
		return nil, status.Errorf(402, "error. username in use")
	}

	u := data.M_User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Password:  req.GetPassword(),
		Active:    true,
	}

	err = a.M_Model.M_User.CreateUser(u)
	if err != nil {
		return nil, err
	}

	res := &auth.CreateUserResponse{
		Created:  true,
		Username: u.Username,
	}
	return res, nil

}

func (a *AuthServer) GetUserByUsername(ctx context.Context, req *auth.GetUserByUsernameRequest) (*auth.GetUserByUsernameResponse, error) {
	username := req.GetUsername()

	log.Println("USERNAME", username)

	u, err := a.M_Model.M_User.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	log.Println(u)
	res := &auth.GetUserByUsernameResponse{
		Found: true,
		User: &auth.M_User{
			Username:  u.Username,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
		},
	}

	return res, nil
}

func (a *AuthServer) AuthUser(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	input := req.GetArgUser()

	// attempt to grab user by passed in email
	user, err := a.Models.User.GetByUsername(input.Username)
	// if the user does not exist
	if err != nil {
		res := &auth.AuthResponse{
			Authed: false,
		}
		return res, err
	}

	// check to see if the passwords match
	valid, err := user.PasswordMatches(input.Password)
	if err != nil || !valid {
		res := &auth.AuthResponse{
			Authed: false,
		}
		return res, err
	}
	res := &auth.AuthResponse{
		Authed: true,
	}
	return res, nil

}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", webPort))
	if err != nil {
		log.Fatalf("failed to listen for grpc %v", err)
	}

	s := grpc.NewServer()

	auth.RegisterAuthServiceServer(s, &AuthServer{Models: app.Models})

	log.Printf("GRPC server started on port %s", webPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for grpc %v", err)
	}
}
