.PHONY:
build:
	go build -o ./.bin/fio cmd/main.go
run: build
	./.bin/fio

.PHONY: up
up:
	migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up
.PHONY: down
down:
	migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down -all
.PHONY: drop
drop:
	migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' drop -f
