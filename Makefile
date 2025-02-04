.PHONY: init, install, docker-up, run, test

init:
	make install

install:
	go install gotest.tools/gotestsum@latest && \
	go mod tidy && go mod vendor

docker-up:
	docker compose up -d 

run:
	sql-migrate up && \
	go run cmd/main.go

test:
	gotestsum --format pkgname