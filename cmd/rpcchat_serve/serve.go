package main

import (
	"go-app-x/internal/env"
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

	if err := RunChatServer(lis); err != nil {
		log.Fatalf("error starting server %s\n", err)
	}
}
