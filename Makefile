BIN_DIR := bin
TOOLS := $(notdir $(wildcard cmd/*))
GO_FILES := $(wildcard cmd/*/*.go)

ifeq ($(OS),Windows_NT)
	RM := del /Q /F
	RMDIR := rmdir /S /Q
	MKDIR := if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	BIN_EXT := .exe
else
	RM := rm -f
	RMDIR := rm -rf
	MKDIR := mkdir -p $(BIN_DIR)
	BIN_EXT :=
endif

all: $(TOOLS)

$(TOOLS):
	@$(MKDIR)
	@go build -o $(BIN_DIR)/$@$(BIN_EXT) ./cmd/$@/main.go

clean:
	@$(RMDIR) $(BIN_DIR)

test:
	@go test ./...

lint:
	@golangci-lint run