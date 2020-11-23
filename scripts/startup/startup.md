# generate Go grpc stubs and encoders/decoders from protocol buffer definition
# input proto file definitions -> outputs grpc stub and encoders/decoders
protoc ./api/proto/order/order.proto --go_out=plugins=grpc:./

# run specific env
go run main.go -GO_ENV prod
go run cmd/location-service/main.go -GO_ENV dev
go run main.go -GO_ENV test

# test recursively
go test ./...

# show all keys in sorted set
zrangebyscore "sorted set index" -inf +inf

# start redis instance on port
redis-server --port 6385

# update all go modules
go list -m -u all

# run unit tests
richgo test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# generate all interface mocks
go generate ./...
