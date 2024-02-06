BIN_DIR := bin
TOOLS := $(notdir $(wildcard cmd/*))
TOOL_TARGETS := $(addprefix $(BIN_DIR)/, $(TOOLS))

all: $(TOOL_TARGETS)

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
