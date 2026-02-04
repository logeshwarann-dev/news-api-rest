GO_BIN?=$(shell pwd)/.bin
SHELL:=env PATH=$(GO_BIN):$(PATH) $(SHELL)


fmt::
	${GO_BIN}/golangci-lint.exe run --fix -v ./...
run::
	go run ./cmd/api-server/main.go
test::
	go test -v -cover ./...
tidy::
	go mod tidy -v

tools:
	mkdir -p $(GO_BIN)
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ${GO_BIN} v1.61.0
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % sh -c 'GOBIN=${GO_BIN} go install %'