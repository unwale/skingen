proto-gen-task:
	@echo ">> Generating gRPC code for task-service..."
	@protoc --go_out=. --go-grpc_out=. \
		--proto_path=./services/task-service/proto \
		./services/task-service/proto/task/v1/task.proto

test-unit:
	@echo ">> Running unit tests..."
	@go test -v ./... -coverprofile=coverage.txt

test-integration:
	@echo ">> Running integration tests..."
	@touch .test.env; \
	set -a; \
	source .test.env; \
	set +a; \
	docker compose -f docker-compose.test.yaml up -d --build; \
	go test --tags=integration -p=1 ./... -coverprofile=coverage.txt; \
	EXIT_CODE=$$?; \
	docker compose -f docker-compose.test.yaml down --remove-orphans; \
	exit $$EXIT_CODE