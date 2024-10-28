BINARY = cf-backend
BIN_DIR = bin

build:
	@go build -o $(BIN_DIR)/$(BINARY)

run: build
	@clear
	@./$(BIN_DIR)/$(BINARY)

clean:
	@rm -rf $(BIN_DIR)

test:
	@go test ./... -v

deploy:
	@terraform -chdir=terraform apply

destroy:
	@terraform -chdir=terraform destroy
