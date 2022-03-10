.PHONY: lambda deps build zip clean run mockgen test

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

mockgen: pkg/mock_pkg/mock_pkg.go

pkg/mock_pkg/mock_pkg.go:
	mockgen -destination $@ github.com/takaiyuk/kakeibo-go/pkg InterfaceIFTTT,InterfaceService,InterfaceSlackClient

test:
	go test ./pkg/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

# bench:
# 	go test -bench=. -benchmem ./pkg/...
