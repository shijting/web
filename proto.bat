protoc -I protos -I third_party/grpc-gateway/ -I third_party/googleapis --go_out=plugins=grpc,paths=source_relative:protos --grpc-gateway_out=paths=source_relative:./protos --openapiv2_out=third_party/OpenAPI/ protos/users.proto
statik -m -f -src third_party/OpenAPI/
