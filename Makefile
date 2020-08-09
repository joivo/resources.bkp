.PHONY: build
build:
	@GOOS=linux GOARCH=amd64 go build -o osbckp cmd/main.go

.PHONY: clean
clean:
	@rm -rf osbckp.log osbckp