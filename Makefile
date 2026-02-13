GO_BIN?=$(shell pwd)/.bin
SHELL:=env PATH=$(GO_BIN):$(PATH) $(SHELL)

# format the code
fmt::
	${GO_BIN}/golangci-lint.exe run --fix -v ./...

# run generate command
generate::
	go generate ./...

# run the server
run::
	go run ./cmd/api-server/main.go

# run tests
test::
	go clean -testcache && go test -v -cover ./...

test-integration::
	go clean -testcache && go test -v -tags integration ./test/...

# run tidy
tidy::
	go mod tidy -v

# setup tools
tools::
	mkdir -p $(GO_BIN)
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ${GO_BIN} v2.0.0
# 	GOBIN=${GO_BIN} go install tool
	go install github.com/golang/mock/mockgen@v1.6.0