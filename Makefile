SHELL := /bin/bash

.PHONY: help

help: ## Display this help message
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

topic: ## Create a topic with group
	docker-compose exec kafka kafka-topics --create --topic=test --bootstrap-server=localhost:9092 --partitions=3 && \
	kafka-console-consumer --bootstrap-server=localhost:9092 --topic=test --group=goapp-group

producer: ## Run producer
	@echo "Producer"
	docker-compose exec go-kakfa go run cmd/producer/main.go

consumer: ## Run consumer
	@echo "Consumer"
	docker-compose exec go-kakfa go run cmd/consumer/main.go
