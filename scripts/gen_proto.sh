#!/bin/bash
protoc --go_out=api --go-grpc_out=api --proto_path=api/pricing api/pricing/*.proto
protoc --go_out=api --go-grpc_out=api --proto_path=api/booking api/booking/*.proto
protoc --go_out=api --go-grpc_out=api --proto_path=api/send api/send/*.proto