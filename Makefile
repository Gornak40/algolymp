BIN_DIR := bin
TOOLS := $(notdir $(wildcard cmd/*))
TOOL_TARGETS := $(addprefix $(BIN_DIR)/, $(TOOLS))
GO_FILES := $(shell find $(SRC_DIRS) -type f -name '*.go')

all: $(TOOL_TARGETS)

$(TOOL_TARGETS): $(GO_FILES)

$(BIN_DIR)/%: cmd/%/main.go
	@mkdir -p $(BIN_DIR)
	go build -o $@ $<

clean:
	rm -f $(BIN_DIR)/*
	@rmdir $(BIN_DIR)

test:
	echo "No tests"

lint:
	golangci-lint run
