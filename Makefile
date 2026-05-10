# GO_DIR := ./go
# GO_BIN := $(GO_DIR)/bin/go
#
MAIN_GO := ./app/cmd/testum/main.go
build:
	docker compose up --build -d
	

# .PHONY: build clear

# build: check_go run
#
# check_go:
# 	@command -v go >/dev/null 2>&1 || $(MAKE) install_go
#
# install_go:
# 	@echo "Installing Go..."
# 	@mkdir -p $(GO_DIR)
# 	@wget -q https://go.dev/dl/go1.26.1.linux-amd64.tar.gz -O go.tar.gz
# 	@tar -xzf go.tar.gz -C $(GO_DIR) --strip-components=1
# 	@rm go.tar.gz
# 	@echo "Go installed locally in $(GO_DIR)"
#
# run:
# 	@echo "Running application..."
# 	@PATH=$(GO_DIR)/bin:$$PATH go run $(MAIN_GO)
#
# clear:
# 	@echo "Removing local installations..."
# 	@rm -rf $(GO_DIR)
# 	@echo "Done"
#
test:
	@echo "Running tests..."
	@go test ./... -coverprofile=coverage.out -p=1
	@go tool cover -html=coverage.out -o coverage.html
	@xdg-open coverage.html
	@echo "Done"
