.PHONY: test
test:
	@echo "Starting MongoDB..."
	docker-compose -f ./docker-compose.yml up -d
	@echo "Running tests..."
	MONGODB_URI="mongodb://username:password@localhost:27017" go test -v -p 1 -count=1 -race -cover ./...
	@echo "Stopping MongoDB..."
	docker-compose -f ./docker-compose.yml down