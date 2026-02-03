fmt::
	go fmt ./...
run::
	go run ./cmd/api-server/main.go
test::
	go test -v -cover ./...
tidy::
	go mod tidy -v