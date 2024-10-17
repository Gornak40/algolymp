BIN_DIR := bin
TOOLS := $(notdir $(wildcard cmd/*))
GO_FILES := $(wildcard cmd/*/*.go)

all: $(TOOLS)

$(TOOLS):
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$@ ./cmd/$@/main.go

clean:
	@rm -rf $(BIN_DIR)

test:
	@go test ./...

lint:
	@golangci-lint run
