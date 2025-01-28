.PHONY: dc run test lint

dc:
	docker-compose up --remove-orphans --build

build:
	go build -race -o app cmd/main.go

run:
	go build -race -o app cmd/main.go && \
	HTTP_ADDR=:8080 \
	./app

test:
	go test -race ./...

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run ./...



