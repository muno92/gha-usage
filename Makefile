.PHONY: build

build:
	CGO_ENABLED=0 go build -ldflags '-s -w -buildid=' -trimpath -o ghausage
