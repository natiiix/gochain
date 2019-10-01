.PHONY: all
all: proto

.PHONY: proto
proto:
	mkdir proto_types
	protoc --proto_path=./proto/ --go_out=./proto_types/ ./proto/block.proto ./proto/blockchain.proto
