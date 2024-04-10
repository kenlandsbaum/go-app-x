package main

import (
	"context"
	"fmt"
	"go-app-x/internal/chat"
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	bufSize              = 1024 * 1024
	lis                  *bufconn.Listener
	testMessagesReceived []*chat.Chat
)

func init() {
	lis = bufconn.Listen(bufSize)
	go RunTestChatServer(lis)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type testChatServer struct {
	chat.UnimplementedChatServiceServer
}

func (ch testChatServer) Send(stream chat.ChatService_SendServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error receiving stream %s\n", err)
			return err
		}
		testMessagesReceived = append(testMessagesReceived, msg)
	}
	return nil
}

func RunTestChatServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	chatServiceServer := testChatServer{}
	chat.RegisterChatServiceServer(grpcServer, chatServiceServer)
	log.Println("starting streaming chat rpc server...")
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func TestChatClient(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)
	HandleChatClient(client)

	if len(testMessagesReceived) != len(chatMessages) {
		t.Errorf("got %d messages but expected %d\n", len(testMessagesReceived), len(chatMessages))
	}
	for i, msg := range testMessagesReceived {
		actual := msg.Content
		expected := fmt.Sprintf("msg %d", i+1)
		if actual != expected {
			t.Errorf("got %s but expected %s\n", actual, expected)
		}
	}
}
