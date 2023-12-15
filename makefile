.PHONY: proto

dev-mac:
	env GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.Version=1.0.0'" -o bin/runpodctl .
dev:
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=1.0.0'" -o bin/runpodctl .
lint:
	golangci-lint run
