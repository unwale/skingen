proto-gen-task:
	@echo ">> Generating gRPC code for task-service..."
	@protoc --go_out=. --go-grpc_out=. \
		--proto_path=./services/task-service/proto \
		./services/task-service/proto/task/v1/task.proto

proto-gen-model-server:
	@echo ">> Generating gRPC code for model-server..."
	@protoc --go_out=. --go-grpc_out=. \
		--proto_path=./services/model-server/proto \
		./services/model-server/proto/model/v1/model.proto