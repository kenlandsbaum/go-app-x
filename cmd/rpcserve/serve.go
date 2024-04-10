package main

import (
	"go-app-x/internal/env"
	"go-app-x/internal/rpcserver"
	"log"
	"net"
	"os"
)

func main() {
	env.Load(".env")

	lis, err := net.Listen("tcp", os.Getenv("RPC_ADDRESS"))
	if err != nil {
		log.Fatalf("error creating listener %s\n", err.Error())
	}
	if err := rpcserver.RunServer(lis); err != nil {
		log.Fatalf("error serving %s\n", err.Error())
	}
}
