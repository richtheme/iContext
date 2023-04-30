.PHONY:
.SILENT:

build:
	go build -o ./.bin/api cmd/api/main.go

run: build
	./.bin/api -host=$(host) -port=$(port)

test:
	go test -v ./...

migrate:
	migrate -path schema -database "postgres://localhost:5436/icontext?sslmode=disable&user=postgres&password=qwerty" up
