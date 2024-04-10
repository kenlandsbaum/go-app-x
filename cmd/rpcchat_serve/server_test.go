package main

import (
	"context"
	"go-app-x/internal/chat"
	"io"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	bufSize      = 1024 * 1024
	lis          *bufconn.Listener
	testMessages = []*chat.Chat{
		{Content: "msg 1", Username: "ken"},
		{Content: "msg 2", Username: "ken"},
		{Content: "msg 3", Username: "ken"},
		{Content: "msg 4", Username: "ken"},
		{Content: "msg 5", Username: "ken"},
		{Content: "msg 6", Username: "ken"},
	}
)

func init() {
	lis = bufconn.Listen(bufSize)
	go RunChatServer(lis)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestChatServer(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	testClient := chat.NewChatServiceClient(conn)

	actualSent := handleTestClient(t, testClient)
	if len(actualSent) != len(testMessages) {
		t.Errorf("got %d messages but expected %d\n", len(actualSent), len(testMessages))
	}
}

func handleTestClient(t *testing.T, client chat.ChatServiceClient) []*chat.Chat {
	var messagesSent []*chat.Chat
	stream, err := client.Send(context.Background())
	if err != nil {
		t.Errorf("error sending from test client %s\n", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			chatMsg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("error receiving stream in client %s\n", err)
			}
			messagesSent = append(messagesSent, chatMsg)
		}
		close(waitc)
	}()

	for _, msg := range testMessages {
		if err := stream.Send(msg); err != nil {
			t.Errorf("error sending from client %s\n", err)
		}
	}
	stream.CloseSend()
	<-waitc
	return testMessages
}
