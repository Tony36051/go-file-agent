mkdir -p generated/go/file_transfer
protoc --go_out=paths=source_relative:generated/go/file_transfer --go-grpc_out=paths=source_relative:generated/go/file_transfer api/v1/file_transfer.proto