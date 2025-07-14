proto-gen-task:
	@echo ">> Generating gRPC code for task-service..."
	@protoc --go_out=. --go-grpc_out=. \
		--proto_path=./services/task-service/proto \
		./services/task-service/proto/task/v1/task.proto