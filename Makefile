.PHONY: build protogen run clean migrate test cover

dev: 
	go run main.go

test:
	go test ./... -coverprofile=testprofile.out

cover:
	go tool cover -html=testprofile.out
