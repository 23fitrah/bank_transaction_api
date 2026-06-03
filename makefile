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

sqlserver:
	docker run --platform linux/amd64 -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=P@ssw0rd" -p 1433:1433 --name sqlserver -d mcr.microsoft.com/mssql/server:2022-latest