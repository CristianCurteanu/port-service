
compile-deps:
	go run github.com/google/wire/cmd/wire ./cmd/
	
run-http:
	docker-compose build rest-api
	docker-compose up rest-api mongodb

build:
	go build -o bin/server cmd/api/main.go

run:
	./bin/server --port=${PORT} --mongo-db-uri=${MONGO_DB_URI} --mongo-db-name=${MONGO_DB_NAME}
