package rpcclient

import (
	"encoding/json"
	"go-app-x/internal/pb/user"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateConnection(address string) *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}
	return conn
}

func MarshalUser(response *user.User) ([]byte, error) {
	bts, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return bts, nil
}

func UnmarshalUser(userBts []byte) (*user.User, error) {
	var newUser user.User
	if err := json.Unmarshal(userBts, &newUser); err != nil {
		return nil, err
	}
	return &newUser, nil
}
