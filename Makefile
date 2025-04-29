ifeq ($(POSTGRES_SETUP_TEST),)
	# make test-migration-up POSTGRES_SETUP_TEST='user=postgres dbname=bookstore host=localhost port=5432 sslmode=disable'
	POSTGRES_SETUP_TEST := user=postgres dbname=sso host=localhost port=5432 sslmode=disable
endif

MIGRATION_FOLDER_BOOKS=$(CURDIR)/migrations

# .PHONY говорит о том, чтобы make не искал названия с такими файлами
#.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER_BOOKS)" create "$(name)" sql

#.PHONY: books-migration-up
sso-migration-up:
	goose -dir "$(MIGRATION_FOLDER_BOOKS)" postgres "$(POSTGRES_SETUP_TEST)" up

#.PHONY: books-migration-down
sso-migration-down:
	goose -dir "$(MIGRATION_FOLDER_BOOKS)" postgres "$(POSTGRES_SETUP_TEST)" down

protoc-inventory:
	protoc --go_out=. --go-grpc_out=. api/sso.proto
