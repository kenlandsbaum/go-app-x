package main

import (
	"context"
	"go-app-x/internal/rpcserver"
	"go-app-x/internal/user"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	go rpcserver.RunServer(lis)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestServerStream(t *testing.T) {
	expected := []*user.User{
		{FirstName: "ken"},
		{FirstName: "jen"},
		{FirstName: "sam"},
		{FirstName: "jill"},
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	testClient := user.NewUserServiceClient(conn)

	actual := doTheNeedful(ctx, testClient)
	for i, ex := range expected {
		if actual[i].FirstName != ex.FirstName {
			t.Errorf("got %s but expected %s\n", actual[i].FirstName, ex.FirstName)
		}
	}
}
