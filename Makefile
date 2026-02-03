.PHONY: build test run docker-restart
build:
	go build -v ./cmd/app

run:
	go run -v ./cmd/app

docker-restart:
	docker-compose stop
	docker-compose build --no-cache
	docker-compose up -d

DEFAULT: build
