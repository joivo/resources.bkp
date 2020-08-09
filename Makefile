.PHONY: build
build:
	@go fmt ./... && GOOS=linux GOARCH=amd64 go build -o sssaver cmd/main.go

.PHONY: clean
clean:
	@rm -rf sssaver.log sssaver
