testing_unit:
	go test ./... -v

install_oapi_codegen:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

install_sqlmap:
	sudo apt install sqlmap	

web_api:
	go run ./cmd/reports/* web-api

generate_oapi_server:
	oapi-codegen --config=./api/web/cfg.yaml ./api/web/openapi.yaml

