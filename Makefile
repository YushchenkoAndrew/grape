.PHONY: dev build migrate test cover swagger

dev: 
	go run main.go

build: swagger
	go build -o ./grape main.go

csv: 
	docker cp ./tmp/Colors.csv psql:/home
	docker cp ./tmp/Patterns.csv psql:/home
	docker cp ./tmp/GeoLite2-Country-Blocks.csv psql:/home
	docker cp ./tmp/GeoLite2-Country-Locations-en.csv psql:/home

test:
	config="$(PWD)" go test -count=1 -v ./...

migrate:
	read -p "Migration name: " DESC; \
	goose -dir ./migrations create $$DESC go

swagger:
	swag init -g ./main.go --parseDependency

cover:
	go test ./... -coverprofile=testprofile.out
	go tool cover -html=testprofile.out
