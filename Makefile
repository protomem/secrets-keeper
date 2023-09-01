
.DEFAULT_GOAL := all


.PHONY: all
all: run-docker


.PHONY: run-local
run-local: BIND_ADDR="localhost:8080"
run-local:
	@mkdir -p ./data
	@BIND_ADDR=${BIND_ADDR} DATABASE="./data/data.db" go run ./cmd/secrets-keeper/


.PHONY: run-web-local
run-web-local: API_URL="localhost:8080"
run-web-local:
	@cd web && VITE_API_URL=${API_URL} yarn dev


.PHONY: run-docker
run-docker: APP_PORT=8080
run-docker: WEB_PORT=80
run-docker: APP_ADDR=0.0.0.0:${APP_PORT}
run-docker:
	@APP_PORT=${APP_PORT} WEB_PORT=${WEB_PORT} APP_ADDR=${APP_ADDR} \
		docker compose -p secrets-keeper -f ./build/docker-compose.yaml up -d


.PHONY: stop-docker
stop-docker:
	@docker compose -p secrets-keeper -f ./build/docker-compose.yaml down

