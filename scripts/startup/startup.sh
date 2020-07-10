# generate grpc stubs and encoders/decoders from protocol buffer definition
# (codegen dest directory to create, proto file definition)
protoc --go_out=plugins=grpc:./ driver.proto

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
