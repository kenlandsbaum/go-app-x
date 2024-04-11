package main

import (
	"go-app-x/internal/pb/chat"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type chatserver struct {
	chat.UnimplementedChatServiceServer
}

func (ch chatserver) Send(stream chat.ChatService_SendServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error receiving stream %s\n", err)
			return err
		}
		log.Println("msg rec'd", msg)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func RunChatServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	chatServiceServer := chatserver{}
	chat.RegisterChatServiceServer(grpcServer, chatServiceServer)
	log.Println("starting streaming chat rpc server...")
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
