default: build

build:
	go build  ./cmd/status-board

test:
	go test -race ./... -count=1

lint:
	golint -set_exit_status ./...