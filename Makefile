.PHONY: lint

lint:
	golangci-lint run -c .golangci.yaml