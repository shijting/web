install:
	go get \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		github.com/envoyproxy/protoc-gen-validate \
		github.com/rakyll/statik

gen-proto:
	protoc \
		-I protos \
		-I third_party/grpc-gateway/ \
		-I third_party/googleapis \
		--go_out=plugins=grpc,paths=source_relative:protos \
		--grpc-gateway_out=paths=source_relative:./protos \
		--openapiv2_out=third_party/OpenAPI/ \
		--validate_out=lang=go:protos \
		protos/users.proto

statik:
	statik -m -f -src third_party/OpenAPI/