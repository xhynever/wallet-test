
up:

docker-compose up --build
	@docker compose up db -d
	@make up
	@go run cmd/app/*.go
	

# TESTS
test:
	@docker compose up db -d
	@make up
	@go test -v -cover  -coverpkg=./... ./...
	@docker compose stop

testCI:
	@go test -v -cover  -coverpkg=./... ./...
	
coverage:
	@docker compose up db -d
	@make up
	@go test -coverprofile=coverage.out -coverpkg=./... ./...
	@go tool cover -html=coverage.out
	@rm coverage.out
	@docker compose stop




