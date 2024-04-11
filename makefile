
test:
	go test ./... -coverprofile cover.out | grep -v ".pb.go"

cov:
	go tool cover -func cover.out | grep -v ".pb.go"

userp:
	protoc --proto_path=protos --go_out=internal/user --go_opt=paths=source_relative --go-grpc_out=internal/user --go-grpc_opt=paths=source_relative protos/user.proto

rpcc:
	go run ./cmd/rpcclient/...

rpccs:
	go run ./cmd/rpcssclient/...

rpcs:
	go run ./cmd/rpcserve/...

chatp:
	protoc --proto_path=protos --go_out=internal/chat --go_opt=paths=source_relative --go-grpc_out=internal/chat --go-grpc_opt=paths=source_relative protos/chat.proto

chats:
	go run ./cmd/rpcchat_serve/...

chatc:
	go run ./cmd/rpcchat_client/...