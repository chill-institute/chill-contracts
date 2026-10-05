module example.com/chill-contracts-consumer-go

go 1.26.6

require (
	github.com/chill-institute/chill-contracts/v2 v2.7.6
	google.golang.org/protobuf v1.36.12
)

require (
	connectrpc.com/connect v1.21.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.31.0 // indirect
)

replace github.com/chill-institute/chill-contracts/v2 => ../../..
