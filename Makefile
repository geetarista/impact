build:
	@goxc

test:
	@go test ./...
	@if [ -f impact.test ]; then rm impact.test; fi

.PHONY: test
