module example.com/chill-contracts-consumer-go

go 1.26.5

require github.com/chill-institute/chill-contracts v0.0.0

require (
	connectrpc.com/connect v1.19.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.8 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/chill-institute/chill-contracts => ../../..
