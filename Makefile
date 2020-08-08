build:
	@GOOS=linux GOARCH=amd64 go build -o osbckp cmd/main.go

clean:
	@rm -rf osbckp.log osbckp