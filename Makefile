LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=${LOCAL_BIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54

lint:
	./bin/golangci-lint run ./... --config .golangci.pipeline.yaml