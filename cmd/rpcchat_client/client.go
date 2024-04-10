package main

import (
	"context"
	"go-app-x/internal/chat"
	"go-app-x/internal/env"
	"go-app-x/internal/rpcclient"
	"io"
	"log"
	"os"
)

var (
	chatMessages = []*chat.Chat{
		{Content: "msg 1", Username: "ken"},
		{Content: "msg 2", Username: "ken"},
		{Content: "msg 3", Username: "ken"},
		{Content: "msg 4", Username: "ken"},
		{Content: "msg 5", Username: "ken"},
		{Content: "msg 6", Username: "ken"},
	}
)

func main() {
	env.Load(".env")
	conn := rpcclient.CreateConnection(os.Getenv("RPC_ADDRESS"))
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)

	stream, err := client.Send(context.Background())
	if err != nil {
		log.Fatalf("could not send messages from client %s", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			chatMsg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error receiving stream in client %s\n", err)
			}
			log.Println("msg rec'd", chatMsg)
		}
		close(waitc)
	}()

	for _, msg := range chatMessages {
		if err := stream.Send(msg); err != nil {
			log.Fatalf("error sending from client %s\n", err)
		}
	}
	stream.CloseSend()
	<-waitc
	log.Println("done")
}
