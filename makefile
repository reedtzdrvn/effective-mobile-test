# Makefile

PROJECT_NAME ?= $(notdir $(CURDIR))

NETWORK        ?= $(PROJECT_NAME)_default

DB_DSN         = postgres://user:password@db:5432/subscriptions?sslmode=disable

.PHONY: all docker-up docker-down docker-restart migrate docker-init goose-img

all: docker-init

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-restart: docker-down docker-up

goose-img:
	docker build -t goose-run -f goose.Dockerfile .

migrate: goose-img
	docker run --rm \
	  --network $(NETWORK) \
	  -v ${CURDIR}/migrations:/migrations \
	  goose-run \
	  -dir /migrations postgres "$(DB_DSN)" up


docker-init: docker-up migrate

wait-for-db:
	@echo "Waiting for '$(NETWORK)'..."
	@sleep 5
