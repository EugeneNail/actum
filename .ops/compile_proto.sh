#!/bin/bash

find grpc/proto -name "*.proto" | while read -r proto_file; do
    dir_name=$(basename "${proto_file}" .proto)
    mkdir -p "grpc/gen/${dir_name}"
    protoc --proto_path=grpc/proto \
           --go_out=grpc/gen/"${dir_name}" \
           --go_opt=paths=source_relative \
           --go-grpc_out=grpc/gen/"${dir_name}" \
           --go-grpc_opt=paths=source_relative "${proto_file}"
done
