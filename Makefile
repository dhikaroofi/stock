run: ## Running code at local for testing
	@go run main.go

check: ## Running Code Dependency Check
	@go mod tidy
	@go mod download
	@go mod verify

