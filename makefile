
test:
	go test ./... -coverprofile cover.out

cov:
	go tool cover -func cover.out

userp:
	protoc --proto_path=protos --go_out=internal/user --go_opt=paths=source_relative --go-grpc_out=internal/user --go-grpc_opt=paths=source_relative protos/user.proto

rpcc:
	go run ./cmd/rpcclient/...

rpccs:
	go run ./cmd/rpcssclient/...

rpcs:
	go run ./cmd/rpcserve/...