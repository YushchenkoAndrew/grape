.PHONY: build protogen run clean migrate test cover

dev: 
	go run main.go

csv: 
	docker cp ./tmp/default_patterns.csv psql:/home
	docker cp ./tmp/GeoLite2-Country-Blocks.csv psql:/home
	docker cp ./tmp/GeoLite2-Country-Locations-en.csv psql:/home

test:
	go test -v ./...

cover:
	go test ./... -coverprofile=testprofile.out
	go tool cover -html=testprofile.out
