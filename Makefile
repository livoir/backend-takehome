connection := "$(BACKEND_TAKE_HOME_MYSQL_USER):$(BACKEND_TAKE_HOME_MYSQL_PASSWORD)@tcp($(BACKEND_TAKE_HOME_MYSQL_HOST):$(BACKEND_TAKE_HOME_MYSQL_PORT))/$(BACKEND_TAKE_HOME_MYSQL_DATABASE)?parseTime=true&multiStatements=true"
dir := ./app/db/migrations
goose := goose -dir $(dir) mysql $(connection)

migration-status:
	$(goose) status
migration-create:
	$(goose) create $(name) sql
migration-up:
	$(goose) up
migration-down:
	$(goose) down

test:
	go test -v -cover ./...

run:
	docker compose up -d --build $(svc)

create-private-key:
	openssl genrsa -out $(BACKEND_TAKE_HOME_JWT_PRIVATE_KEY_PATH) 4096
	openssl rsa -in $(BACKEND_TAKE_HOME_JWT_PRIVATE_KEY_PATH) -pubout -out $(BACKEND_TAKE_HOME_JWT_PUBLIC_KEY_PATH)

create-public-key:
	openssl rsa -in $(BACKEND_TAKE_HOME_JWT_PRIVATE_KEY_PATH) -pubout -out $(BACKEND_TAKE_HOME_JWT_PUBLIC_KEY_PATH)