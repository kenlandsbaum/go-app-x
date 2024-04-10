
test:
	go test ./... -coverprofile cover.out

cov:
	go tool cover -func cover.out

userp:
	protoc --go_out=internal/user --go_opt=paths=source_relative --go-grpc_out=internal/user --go-grpc_opt=paths=source_relative user.proto

rpcc:
	go run ./cmd/rpcclient/...

rpcs:
	go run ./cmd/rpcserve/...