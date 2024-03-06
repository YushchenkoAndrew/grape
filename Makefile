.PHONY: build protogen run clean migrate test cover

dev: csv
	go run main.go

csv: 
	docker cp ./tmp/Colors.csv psql:/home
	docker cp ./tmp/Patterns.csv psql:/home
	docker cp ./tmp/GeoLite2-Country-Blocks.csv psql:/home
	docker cp ./tmp/GeoLite2-Country-Locations-en.csv psql:/home

test:
	go test -v ./...

cover:
	go test ./... -coverprofile=testprofile.out
	go tool cover -html=testprofile.out
