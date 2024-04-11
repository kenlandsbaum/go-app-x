package rpcserver

import (
	"context"
	"errors"
	"go-app-x/internal/pb/user"
	"log"
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
	go RunServer(lis)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func assertOnResults(t *testing.T, res *user.User, expected *user.User) {
	if res.FirstName != expected.FirstName {
		t.Errorf("got %s but expected %s\n", res.FirstName, expected.FirstName)
	}
	if res.Email != expected.Email {
		t.Errorf("got %s but expected %s\n", res.Email, expected.Email)
	}
}

func getTestConnection(ctx context.Context) *grpc.ClientConn {
	conn, err := grpc.DialContext(
		ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}
func TestUserServiceServerGetOne(t *testing.T) {
	ctx := context.Background()
	conn := getTestConnection(ctx)
	defer conn.Close()
	testClient := user.NewUserServiceClient(conn)
	type getCases struct {
		in          *user.UserRequest
		expected    *user.User
		expectedErr error
	}
	tcs := map[string]getCases{
		"ken":  {in: &user.UserRequest{Id: 1}, expected: someUsers[0], expectedErr: nil},
		"jen":  {in: &user.UserRequest{Id: 2}, expected: someUsers[1], expectedErr: nil},
		"bill": {in: &user.UserRequest{Id: 5}, expected: nil, expectedErr: errors.New("no such user")},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			res, err := testClient.GetOne(ctx, tc.in)
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("got error %s but expected nil\n", err)
				}
			}
			if res != nil {
				if tc.expected == nil {
					t.Errorf("got a response %v but expected nil\n", res)
				} else {
					assertOnResults(t, res, tc.expected)
				}
			}
		})
	}
}

func TestUserServiceServerCreate(t *testing.T) {
	ctx := context.Background()
	conn := getTestConnection(ctx)
	defer conn.Close()
	testClient := user.NewUserServiceClient(conn)

	res, err := testClient.Create(
		ctx, &user.NewUserRequest{FirstName: "testfirst", LastName: "testlast", Email: "t@t.com"})
	if err != nil {
		t.Errorf("expected nil error but got %s\n", err)
	}

	if res.Id != 5 {
		t.Errorf("got %d but expected 5\n", res.Id)
	}
}
