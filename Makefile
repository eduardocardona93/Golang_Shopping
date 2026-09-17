.PHONY: run build test vet fmt tidy docker-up docker-down docker-logs

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l .

tidy:
	go mod tidy

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f
