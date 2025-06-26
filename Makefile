dc:
	docker-compose up  --remove-orphans --build

build:
	go build -o app_port cmd/main.go

buildLinux:
	GOOS=linux go build -o app_port cmd/main.go

test:
	go test -v ./...

test-inner:
	go test -v ./inner/...

install-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

lint:
	golangci-lint run ./...