# Go Tools
GO  = GO111MODULE=on go
all: fmt check build test info
clean:
	rm -f coverage.txt
	rm -rf vendor
	rm -rf go.sum
deps: clean
	${GO} mod download
fmt:
	GO111MODULE=on ${GO} fmt ./...
generate:
	${GO} generate
