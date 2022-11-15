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

	u, err := a.M_Model.M_User.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
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
	// input := req.GetArgUser()

	arg_user := req.GetArgUser()

	if arg_user.GetUsername() == "" && arg_user.GetEmail() == "" {
		return nil, fmt.Errorf("error. no email or username supplied")
	}

	var user *data.M_User

	var err error
	// determine which auth strategy should be used
	if arg_user.GetUsername() != "" {
		user, err = a.M_Model.M_User.GetUserByUsername(arg_user.GetUsername())
	} else {
		user, err = a.M_Model.M_User.GetUserByEmail(arg_user.GetEmail())
	}
	if err != nil {
		return nil, err
	}
	result, err := user.PasswordMatches(arg_user.GetPassword())

	if err != nil {
		return nil, err
	}

	res := &auth.AuthResponse{
		Authed: result,
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
