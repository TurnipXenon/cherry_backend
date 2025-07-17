.PHONY: proto clean

# Variables
PROTO_DIR := proto
PROTO_OUT := pkg/api/v1
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
# Default to auto-detection
USE_UNIX_COMMANDS ?= 0

# Check if we're in Git Bash (MSYSTEM is set in Git Bash)
ifdef MSYSTEM
    USE_UNIX_COMMANDS = 1
endif

# Generate Go code from protobuf definitions
proto:
	@echo "Generating Go code from protobuf definitions..."
ifeq ($(USE_UNIX_COMMANDS),1)
	@mkdir -p $(PROTO_OUT)
else ifeq ($(OS),Windows_NT)
	@if not exist $(subst /,\,$(PROTO_OUT)) mkdir $(subst /,\,$(PROTO_OUT))
else
	@mkdir -p $(PROTO_OUT)
endif
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=. \
		--go_opt=module=cherry_backend \
		--go-grpc_out=. \
		--go-grpc_opt=module=cherry_backend \
		$(PROTO_FILES)
	@echo "Done!"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
ifeq ($(USE_UNIX_COMMANDS),1)
	@rm -rf $(PROTO_OUT)
else ifeq ($(OS),Windows_NT)
	@if exist $(subst /,\,$(PROTO_OUT)) rmdir /s /q $(subst /,\,$(PROTO_OUT))
else
	@rm -rf $(PROTO_OUT)
endif
	@echo "Done!"
