DOCKER_COMPOSE=docker-compose

local:
	go run ./cmd/main.go

docker-build:
	$(DOCKER_COMPOSE) build

docker-up:
	$(DOCKER_COMPOSE) up

docker-clean:
	$(DOCKER_COMPOSE) down

docker:
	$(DOCKER_COMPOSE) build;
	$(DOCKER_COMPOSE) up

.PHONY: docker-build docker-up docker-clean

