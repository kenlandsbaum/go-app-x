package main

import (
	"context"
	"go-app-x/internal/env"
	"go-app-x/internal/pb/user"
	"go-app-x/internal/router"
	"go-app-x/internal/rpcclient"
	"os"
)

func main() {
	env.Load(".env")
	ctx := context.Background()
	conn := rpcclient.CreateConnection(os.Getenv("RPC_ADDRESS"))
	defer conn.Close()

	client := user.NewUserServiceClient(conn)
	r := router.New()

	wserver := createServer(ctx, client, r)
	wserver.routes()
	wserver.start()
}
