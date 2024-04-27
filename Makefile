.PHONY: build
build:
	@echo "Building Docker image..."
	docker-compose build

.PHONY: run
run:
	@echo "Starting Docker container..."
	docker-compose up -d

.PHONY: stop
stop:
	@echo "Stopping Docker container..."
	docker-compose down

.PHONY: clean
clean:
	@echo "Cleaning up..."
	docker-compose down -v --rmi local

.PHONY: test
test:
	go test -v $$(go list ./...) --count=1