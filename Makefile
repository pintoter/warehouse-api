.DEFAULT_GOAL = up

.PHONY: up
up:
	docker-compose -f docker-compose.yml up --remove-orphans

.PHONY: stop
stop:
	docker-compose down --remove-orphans

.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $$3}'

.PHONY: lint
lint:
	golangci-lint run ./...
