

.PHONY: all clean

# Define the output directory
BIN_DIR := ./bin

all: $(BIN_DIR)/ics-to-matrix $(BIN_DIR)/matrix-to-ics 

$(BIN_DIR)/ics-to-matrix:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/ics-to-matrix ./cmd/ics-to-matrix

$(BIN_DIR)/matrix-to-ics:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/matrix-to-ics ./cmd/matrix-to-ics

$(BIN_DIR)/tui:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/tui ./cmd/tui


## test: Run all Go unit tests and integration tests
test: all
	@echo "Running Go unit tests..."
	go test -v ./cmd/...
	@echo "Running integration pipe test..."
	@mkdir -p test_output
	./bin/ics-to-matrix --input testdata/jameson.m.jolley@gmail.com.ics --metadata test_output/meta.json --output test_output/matrix1.txt
	cat test_output/matrix1.txt | uiua scripts/pass.ua | ./bin/matrix-to-ics --stdin --metadata test_output/meta.json --output test_output/roundtrip.ics
	@echo "Testing round-trip ICS -> matrix -> ICS -> matrix..."
	./bin/ics-to-matrix --input test_output/roundtrip.ics --metadata test_output/meta2.json --output test_output/matrix2.txt
	@echo "Integration test complete. Output saved to test_output/"

## fuzz: Run all Go fuzz tests
fuzz:
	go test -fuzz=FuzzParseICS ./cmd/ics-to-matrix
	go test -fuzz=FuzzParseMatrix ./cmd/matrix-to-ics

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
clean:
	rm -rf $(BIN_DIR) test_output
