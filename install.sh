#!/usr/bin/env bash

sudo docker-compose up -d --build --remove-orphans

sudo mkdir -p zoo/data zoo/log broker/data
sudo chmod -R 777 zoo/ broker/

cd user_service && make migrate && cd ../
cd message_service && make migrate && cd ../