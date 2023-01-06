build:
	go build -o ./bin/avl-vs-rb .

test:
	go test -v --count=1 ./...

lint:
	golangci-lint run

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
