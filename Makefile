APP_NAME := magnet2torrent
BIN_DIR := bin

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

.PHONY: test
test:
	GO111MODULE=on go test ./...

# Example cross-compilation (uncomment as needed):
# GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o $(BIN_DIR)/$(APP_NAME).exe ./cmd/$(APP_NAME)
