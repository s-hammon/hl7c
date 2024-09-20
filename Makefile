REPO_PATH=github.com/s-hammon
APP_NAME=hl7c

build:
	@go build ./...

install:
	@go install ./...

clean:
	@rm -rf bin
	@rm -rf internal/objects
	@(rm -f model.go || true)
	@go mod tidy

out: clean build install
	@echo "Built and installed ${APP_NAME} to GOPATH"

.PHONY: build install clean out