module example.com/chill-contracts-consumer-go

go 1.26.6

require github.com/chill-institute/chill-contracts/v2 v2.3.7

require (
	connectrpc.com/connect v1.20.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.30.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/chill-institute/chill-contracts/v2 => ../../..
