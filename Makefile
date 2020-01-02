.PHONY: test

test:
	@cd test && go get -t
	@cd test && go test -v -short
