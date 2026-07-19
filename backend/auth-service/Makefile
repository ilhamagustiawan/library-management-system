.PHONY: help run build swagger test test-race vet fmt precommit migrate-up migrate-down create-client docker-up docker-down

help:
	@echo "Available commands"
	@echo "  make run            Start the auth service"
	@echo "  make build          Build the service"
	@echo "  make swagger        Regenerate Swagger documentation"
	@echo "  make test           Run tests"
	@echo "  make test-race      Run tests with race detection"
	@echo "  make precommit      Format, vet, test, and build"
	@echo "  make migrate-up     Apply database migrations"
	@echo "  make migrate-down   Roll back one migration"
	@echo "  make create-client  Create an additional OAuth client"

run:
	go run . serve

build:
	go build ./...

swagger:
	go run . swagger

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

precommit: fmt vet test-race build

migrate-up:
	go run . migrate --action up

migrate-down:
	go run . migrate --action down

create-client:
	go run . create-client --name "$(NAME)" --id "$(ID)" --kind "$(or $(KIND),authorization_code)" --redirect-uri "$(REDIRECT_URI)" --scopes "$(or $(SCOPES),loans:borrow:self loans:return:self transactions:read:self books:read)"

docker-up:
	docker compose -f ../../docker-compose.yaml up --build auth-service

docker-down:
	docker compose -f ../../docker-compose.yaml down
