run_env_containers:
	docker compose -f docker-compose.yaml up --build --remove-orphans -d

testing_unit:
	go test ./... -v

stop_env_containers:
	docker compose stop

install_oapi_codegen:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

install_sqlmap:
	sudo apt install sqlmap	

web_api:
	go run ./cmd/reports/* web-api

front:
	go run ./cmd/front/main.go

generate_oapi_server:
	oapi-codegen --config=./api/web/cfg.yaml ./api/web/openapi.yaml

