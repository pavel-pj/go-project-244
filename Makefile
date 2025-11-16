.PHONY: lint start build

test:
	 go test ./... 
lint:
	golangci-lint run
start:
	go run cmd/gendiff/main.go		
build:
	go build -o bin/gendiff ./cmd/gendiff	
	bin/gendiff