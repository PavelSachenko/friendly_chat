include .env
export
abs_path:=$(CURDIR)

migrate:
	docker run --rm -v $(abs_path)/migration:/migrations --network host migrate/migrate \
						-path=/migrations/ \
						-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5433/${DB_DATABASE}?sslmode=disable" up

proto:
	protoc -I=. --go_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. pkg/pb/*.proto \
	&& sed -i "s/,omitempty//" pkg/pb/*.pb.go

migrate-create:
	migrate create -ext sql -dir migration friendly_chat_user