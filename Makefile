.PHONY: build protogen run clean migrate test cover

dev: 
	go run main.go

test:
	go test ./...

cover:
	go test ./... -coverprofile=testprofile.out
	go tool cover -html=testprofile.out
