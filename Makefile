BINARY_NAME=live-server

dev:
	air
run:
	go run main.go
build:
	go build -o ${BINARY_NAME}
test:
	go test -cover ./...
fmt:
	go fmt ./...
lint:
	golangci-lint run ./...
clean:
	go clean
	rm ${BINARY_NAME}
