.PHONY: build protogen run clean migrate test cover

dev: 
	go run main.go

test:
	go test -v ./...

cover:
	go test ./... -coverprofile=testprofile.out
	go tool cover -html=testprofile.out
