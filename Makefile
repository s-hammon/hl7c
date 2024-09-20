REPO_PATH=github.com/s-hammon
APP_NAME=hl7c

build:
	@go build -o bin/${APP_NAME}

install:
	@go install ${REPO_PATH}/${APP_NAME}

clean:
	@rm -rf bin
	@rm -rf internal/objects
	@go mod tidy

out: build install
	@echo "Built and installed ${APP_NAME} to GOPATH"

.PHONY: build install clean out