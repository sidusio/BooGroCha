.PHONY: build_client
build_client:
	cd cmd/bgc && go build

.PHONY: build_server
build_server:
	cd cmd/bgc-server && go build

.PHONY: run_server
run_server:
	go run cmd/bgc-server/main.go

.PHONY: test
test:
	go test ./...
