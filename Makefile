run:
	sudo docker-compose up -d

stop:
	sudo docker-compose down

migrate-rollback-user:
	cd user_service/ && make migrate-down && cd ..

migrate-rollback-message:
	cd message_service/ && make migrate-down && cd ..

migrate:
	cd user_service/ && make migrate && cd .. && \
	cd message_service/ && make migrate && cd ..

install:
	sudo ./install.sh