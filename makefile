run:
	go run main.go

build:
	go build main.go

air:
	air -c air.toml

docker-up:
	docker compose up -d --build

docker-run:
	docker compose up -d app
	
container-reset:
	docker-compose down && docker-compose up -d

sqlserver:
	docker run --platform linux/amd64 -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=YOUR_PASSWORD_SQLSERVER" -p 1433:1433 --name sqlserver -d mcr.microsoft.com/mssql/server:2022-latest
