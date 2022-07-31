include .env
export
abs_path:=$(CURDIR)

migrate-create:
	migrate create -ext sql -dir migration friendly_chat_message_service

migrate:
	docker run --rm -v $(abs_path)/migration:/migrations --network host migrate/migrate \
			-path=/migrations/ \
			-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5435/${DB_DATABASE}?sslmode=disable" up

migrate-down:
	docker run --rm -v $(abs_path)/migration:/migrations --network host migrate/migrate \
			-path=/migrations/ \
			-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@0.0.0.0:5435/${DB_DATABASE}?sslmode=disable" down 10000


proto:
	protoc -I=. --go_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. gateway/point/user/pb/*.proto \
	&& sed -i "s/,omitempty//" gateway/point/user/pb/*.pb.go