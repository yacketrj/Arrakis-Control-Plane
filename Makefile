.PHONY: build web go linux dev-server setup deploy-web

build: go

go:
	go build -o arrakis-control-panel .

linux:
	GOOS=linux GOARCH=amd64 go build -o arrakis-control-panel-linux .

dev-server:
	go run .

setup:
	go run . -setup

web:
	cd web && npm ci && npm run build

deploy-web:
	cd web && npm ci && npm run build && wrangler pages deploy dist --project-name arrakis-control-panel
