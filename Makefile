start:
	go run cmd/shortener/main.go

build:
	go build -o shortenerBuild cmd/shortener/main.go

lint:
	go vet ./...

test:
	go test ./...

test-coverage:
	go test ./... -cover
