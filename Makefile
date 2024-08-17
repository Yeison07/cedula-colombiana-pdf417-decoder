# Variables
GOOS = windows
GOARCH = amd64
BUILD_DIR = build
OUTPUT = $(BUILD_DIR)/barcode_scanner.exe

# Phony targets
.PHONY: all clean

# Default target
all: $(OUTPUT)

# Build target
$(OUTPUT): | $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUTPUT)

# Create build directory if it doesn't exist
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Clean target
clean:
	rm -rf $(BUILD_DIR)
