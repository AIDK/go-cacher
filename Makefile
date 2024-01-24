build:
	@go build -o bin/gocacher

run: build
	@go run main.go

test:
	@go test -v ./...