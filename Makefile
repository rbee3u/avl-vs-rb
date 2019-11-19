build:
	go build -o ./bin/avl-vs-rb .

lint:
	golangci-lint run

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
