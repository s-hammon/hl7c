REPO_PATH=github.com/s-hammon
APP_NAME=hl7c

build:
	@go build ./...

install: build
	@go install ./...

clean:
	@rm -rf bin
	@rm -rf internal/objects
	@(rm -f model.go || true)
	@go mod tidy

out: clean build install
	@echo "Built and installed ${APP_NAME} to GOPATH"

test:
	@go test -v -cover ./...

.PHONY: build install clean out