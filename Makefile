dc:
	docker-compose -f docker/docker-compose.yml up  --remove-orphans --build

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

test-coverage:
	go clean -testcache
	go test -v ./inner/... -coverprofile=coverage.tmp.out
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

swag-generate:
	swag init -d cmd,inner --parseDependency --parseInternal
