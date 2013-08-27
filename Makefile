test:
	@go test ./...
	@if [ -f impact.test ]; then rm impact.test; fi

build:
	@goxc

.PHONY: build, test
