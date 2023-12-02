build:
	@go build -o bin/videostream

run: build
	@./bin/videostream