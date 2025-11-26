auth:
	@go run ./services/auth/cmd/auth/main.go

deploy:
	@go run ./services/deploy/cmd/deploy/main.go

upload:
	@go run ./services/upload/cmd/upload/main.go

request:
	@go run ./services/request/cmd/request/main.go
