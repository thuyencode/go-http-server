BINARY_NAME=http-server
ENV_FILE=.env

ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
endif

dev:
	air
run:
	go run .
build:
	GOAMD64=$(GOAMD64) CGO_ENABLED=${CGO_ENABLED} go build -o $(BINARY_NAME)
test:
	go test -cover ./...
fmt:
	go fmt ./...
lint:
	golangci-lint run ./...
clean:
	go clean
	rm $(BINARY_NAME)
	rm -rf tmp
