part1: ## Running code at local for testing
	@go run main.go part1
part2: ## Running code at local for testing
	@go run main.go part2

check: ## Running Code Dependency Check
	@go mod tidy
	@go mod download
	@go mod verify

