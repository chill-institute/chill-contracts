module example.com/chill-contracts-consumer-go

go 1.26.6

require (
	connectrpc.com/connect/v2 v2.0.0
	github.com/chill-institute/chill-contracts/v3 v3.0.0
	google.golang.org/protobuf v1.36.12
)

require github.com/grpc-ecosystem/grpc-gateway/v2 v2.31.0 // indirect

replace github.com/chill-institute/chill-contracts/v3 => ../../..
