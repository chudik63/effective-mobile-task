dev:
	docker compose up postgres migrator
swagger:
	swag init -g cmd/main/main.go -o docs --parseDependency --parseInternal