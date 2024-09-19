APP_NAME=hl7c

build:
	@go build -o bin/${APP_NAME}

run:
	@./bin/${APP_NAME}

clean:
	@rm -rf bin
	@go mod tidy

out: build run