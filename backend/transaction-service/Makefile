.PHONY: run build swagger test test-race vet fmt precommit migrate-up migrate-down docker-up docker-down

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

docker-up:
	docker compose -f ../../docker-compose.yaml up --build transaction-service

docker-down:
	docker compose -f ../../docker-compose.yaml down
