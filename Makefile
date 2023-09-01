
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
run-docker: WEB_PORT=443
run-docker: APP_ADDR=localhost:${APP_PORT}
run-docker:
	@APP_PORT=${APP_PORT} WEB_PORT=${WEB_PORT} APP_ADDR=${APP_ADDR} \
		docker compose -p secrets-keeper -f ./build/docker-compose.yaml up -d


.PHONY: stop-docker
stop-docker:
	@docker compose -p secrets-keeper -f ./build/docker-compose.yaml down


.PHONY: gen-cert
gen-cert: HOSTNAME=localhost
gen-cert:
	@mkdir -p ./configs/certs
	@mkdir -p ./web/certs
	@openssl req -new -subj "/C=RU/ST=Msk/CN=${HOSTNAME}" -newkey rsa:2048 -nodes -keyout ./configs/certs/${HOSTNAME}.key -out ./configs/certs/${HOSTNAME}.csr
	@openssl x509 -req -days 365 -in ./configs/certs/${HOSTNAME}.csr -signkey ./configs/certs/${HOSTNAME}.key -out ./configs/certs/${HOSTNAME}.crt
	@cat ./configs/certs/${HOSTNAME}.crt > ./web/certs/fullchain.pem
	@cat ./configs/certs/${HOSTNAME}.key > ./web/certs/privkey.pem


