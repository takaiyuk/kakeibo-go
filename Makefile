.PHONY: lambda deps build zip clean run test bench

lambda: deps build zip clean

deps:
	go mod tidy
  
build:
	GOOS=linux GOARCH=amd64 go build -o kakeibo ./cmd/lambda/main.go

zip:
	zip -r handler.zip .env kakeibo pkg

clean:
	rm ./kakeibo

run:
	go run ./cmd/kakeibo/main.go

test:
	go test ./pkg/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

bench:
	go test -bench=. -benchmem ./pkg/...
