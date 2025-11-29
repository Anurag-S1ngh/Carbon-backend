auth:
	@go run ./services/auth/cmd/auth/main.go

deploy:
	@go run ./services/deploy/cmd/deploy/main.go

upload:
	@go run ./services/upload/cmd/upload/main.go

request:
	@go run ./services/request/cmd/request/main.go

db:
	@docker run -e POSTGRES_PASSWORD=password -p 5432:5432 -d -v carbon:/var/lib/postgresql/data postgres

redis:
	@docker run -p 6379:6379 -d redis:7-alpine
