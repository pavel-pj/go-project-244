.PHONY: lint start build

test:
	go test ./... 
#go test -coverprofile=coverage.out ./...  \
#go tool cover -func=coverage.out

lint:
	golangci-lint run ./...
	
start:
	go run cmd/gendiff/main.go		
build:
	go build -o bin/gendiff ./cmd/gendiff/main.go	
	
#--bin/gendiff

#	go build -o bin/gendiff ./cmd/main.go