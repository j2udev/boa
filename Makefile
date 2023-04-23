all: tidy fmt lint vet test

lint:
	golangci-lint run ./...

fmt:
	gofmt -w -s .

fmtd:
	gofmt -w -d .

tidy:
	go mod tidy

vet:
	go vet ./...

test:
	go test ./... -v

cover:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
