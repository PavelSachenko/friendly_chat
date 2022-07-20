include .env
export
abs_path:=$(CURDIR)

migrate-create:
	migrate create -ext sql -dir migration friendly_chat_message_service

migrate:
	docker run --rm -v $(abs_path)/migration:/migrations --network host migrate/migrate \
			-path=/migrations/ \
			-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5435/${DB_DATABASE}?sslmode=disable" up