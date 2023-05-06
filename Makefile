
compile-deps:
	go run github.com/google/wire/cmd/wire ./cmd/
	
run-http:
	docker-compose build rest-api
	docker-compose up rest-api mongodb