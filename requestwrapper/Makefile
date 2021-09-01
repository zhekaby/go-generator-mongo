all: generate test build

build:
	CGO_ENABLED=0 go build -a -ldflags="-w -s -X main.VERSION=`git rev-parse HEAD | cut -c1-8`" -o ~/go/bin/requestwraper

generate:
	rm -f tests/requestwrapper_validator_middleware.go
	easyjson tests/models.go
	go run . tests/

test:
	go test -v -coverprofile cover.out ./tests && \
		go tool cover -func cover.out