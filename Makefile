.PHONY: lambda deps build clean zip run test

lambda: deps build zip

deps:
	go mod tidy
  
build:
	GOOS=linux GOARCH=amd64 go build -o kakeibo ./src/main.go

zip:
	zip handler.zip ./kakeibo ./.env

clean:
	rm ./kakeibo ./handler.zip

run:
	go run ./src/main.go

test:
	go test ./src/...
