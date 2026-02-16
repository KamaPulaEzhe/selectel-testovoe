.PHONY: plugin

plugin:
	go build -buildmode=plugin -o mylinter.so ./plugin/golangci-lint

golangci:
	golangci-lint run