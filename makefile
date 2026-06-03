run:
	go run main.go

build:
	go build main.go

air:
	air -c air.toml

docker-up:
	docker compose up -d --build

container-reset:
	docker-compose down && docker-compose up -d