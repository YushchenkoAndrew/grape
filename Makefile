.PHONY: dev migrate test cover

dev: 
	go run main.go

csv: 
	docker cp ./tmp/Colors.csv psql:/home
	docker cp ./tmp/Patterns.csv psql:/home
	docker cp ./tmp/GeoLite2-Country-Blocks.csv psql:/home
	docker cp ./tmp/GeoLite2-Country-Locations-en.csv psql:/home

test:
	go test -v ./...

migrate:
	read -p "Migration name: " DESC; \
	goose -dir ./migrations create $$DESC go

cover:
	go test ./... -coverprofile=testprofile.out
	go tool cover -html=testprofile.out
