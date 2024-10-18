# Start up
## Container operation:
1. To save time, compile binary packages locally:
Run the following file on Windows:
windows.bat
Under Linux, run the command: go build ldflags="- w - s" - o/ main ./main.go 
2. Start the container
Execution:
docker-compose up --build
Backend service port 80080, PostgreSQL database port 5432
## Local debugging:
1. Start PostgreSQL database:
docker compose up db -d
2. Start the backend:
go run ./main.go
## Document Description
The repository data layer is constructed by SQL C using SQL.
Serviccs control layer
Handler API interface
Storage of tests mock files
util       Read configuration, initialize table
## Duration: 48 hours
## Code View
[Route](.\handler\handler.go)
## Unit testing
[handler](.\handler\accounts_test.go)
[servicces](.\services\accounts_test.go)
[sqlc](.\repository\sqlc\account_test.go)
## API online documentation
https://apifox.com/apidoc/shared-09c25764-2738-405e-b9e9-bace92c4cc78
[API Interface Document](.\go-wallet.md)