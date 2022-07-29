#!/usr/bin/env bash

sudo docker-compose up -d --build --remove-orphans
cd user_service && make migrate && cd ../
cd message_service && make migrate && cd ../