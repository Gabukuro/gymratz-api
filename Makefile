init:
	make install

install:
	go install gotest.tools/gotestsum@latest && \
	go mod tidy && go mod vendor

run:
	go run cmd/main.go

test:
	gotestsum --format pkgname