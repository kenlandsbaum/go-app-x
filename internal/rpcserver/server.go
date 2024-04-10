package rpcserver

import (
	"context"
	"errors"
	"go-app-x/internal/user"
	"log"
	"net"

	"google.golang.org/grpc"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
}

func (u UserServiceServer) GetOne(ctx context.Context, in *user.UserRequest) (*user.User, error) {
	log.Printf("received request for user with id %d\n", in.Id)
	id := in.Id
	for _, u := range someUsers {
		if u.Id == id {
			return u, nil
		}
	}
	return nil, errors.New("no such user")
}

func (u UserServiceServer) Create(ctx context.Context, in *user.NewUserRequest) (*user.NewUserResponse, error) {
	newUser := user.User{FirstName: in.FirstName, LastName: in.LastName, Email: in.Email}
	someUsers = append(someUsers, &newUser)
	return &user.NewUserResponse{Id: int64(len(someUsers))}, nil
}

func RunServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	userServiceServer := UserServiceServer{}
	user.RegisterUserServiceServer(grpcServer, userServiceServer)
	log.Println("starting rpc server...")
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
