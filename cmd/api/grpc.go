package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/see-air-uh/asxce-toga/auth"
	"github.com/see-air-uh/asxce-toga/data"
	"google.golang.org/grpc"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Models data.Models
}

type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
