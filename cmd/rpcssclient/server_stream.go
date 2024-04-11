package main

import (
	"context"
	"fmt"
	"go-app-x/internal/env"
	"go-app-x/internal/pb/user"
	"go-app-x/internal/rpcclient"
	"io"
	"os"
)

func main() {
	env.Load(".env")
	ctx := context.Background()
	conn := rpcclient.CreateConnection(os.Getenv("RPC_ADDRESS"))
	defer conn.Close()
	client := user.NewUserServiceClient(conn)

	doTheNeedful(ctx, client)
}

func doTheNeedful(ctx context.Context, client user.UserServiceClient) []*user.User {
	var receivedMsgs []*user.User
	stream, err := client.GetAll(ctx, &user.AllUsersRequest{})
	if err != nil {
		panic(err)
	}
	for {
		userMsg, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("done")
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("user rec'd", userMsg)
		receivedMsgs = append(receivedMsgs, userMsg)
	}
	return receivedMsgs
}
