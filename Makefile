
.DEFAULT_GOAL := all


.PHONY: all
all: run-local


.PHONY: run-local
run-local: BIND_ADDR="localhost:8080"
run-local:
	@mkdir -p ./data
	@BIND_ADDR=${BIND_ADDR} DATABASE="./data/data.db" go run ./cmd/secrets-keeper/


.PHONY: run-web-local
run-web-local: API_URL="localhost:8080"
run-web-local:
	@cd web && VITE_API_URL=${API_URL} yarn dev

