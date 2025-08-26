run:
	cd cmd/api&& go run main.go -e ./../../config.local.yml

migrate:
	cd cmd/migration&& go run main.go -e ./../../config.local.yml

tidy:
	export GOPROXY=https://goproxy.io,direct&& go mod tidy

