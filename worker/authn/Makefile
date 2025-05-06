.PHONY: dev
dev:
	@wrangler d1 execute otakakot --local --file=../../schema/schema.sql
	@wrangler dev --port 7777

.PHONY: build
build:
	@go run github.com/syumai/workers/cmd/workers-assets-gen@v0.28.1 -mode=go
	@GOOS=js GOARCH=wasm go build -o ./build/app.wasm .

.PHONY: deploy
deploy:
	@wrangler deploy
