include .env
export
abs_path:=$(CURDIR)


migrate:
	docker run --rm -v $(abs_path)/migration:/migrations --network host migrate/migrate \
						-path=/migrations/ \
						-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5433/${DB_DATABASE}?sslmode=disable" up