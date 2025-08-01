build:
	@go build -o ./bin/app ./cmd/server/

run: build
	@./bin/app
