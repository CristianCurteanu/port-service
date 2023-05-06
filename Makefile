
compile-deps:
	go run github.com/google/wire/cmd/wire ./cmd/
	
run-container-mongo:
	docker-compose build rest-api-mongo
	docker-compose up rest-api-mongo mongodb

run-container-inmem:
	docker-compose build rest-api-inmem
	docker-compose up rest-api-inmem

build:
	go build -o bin/server cmd/api/main.go

run:
	./bin/server --port=${PORT} --mongo-db-uri=${MONGO_DB_URI} --mongo-db-name=${MONGO_DB_NAME}

test:
	go test -cover -v ./...

test-data-race:
	go test -cover -v -race ./...