up:
	docker-compose -f ./microservices/docker-compose.yml up -d

down:
	docker-compose -f ./microservices/docker-compose.yml down --remove-orphans

build:
	make service
	make docker

service:
	cd ./src/saiAuth && go mod tidy && go build -o ../../microservices/saiAuth/build/sai-auth
	cp ./src/saiAuth/config.json ./microservices/saiAuth/build/config.json


docker:
	docker-compose -f ./microservices/docker-compose.yml up -d --build

log:
	docker-compose -f ./microservices/docker-compose.yml logs -f

loga:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-auth

