start:
	go run cmd/shortener/main.go

lint:
	go vet ./...

test:
	go test ./...

test-coverage:
	go test ./... -cover
