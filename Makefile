postgres:
	docker run --name postgres16alpine -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:16.0-alpine3.18

createdb:
	docker exec -it postgres16alpine createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16alpine dropdb simple_bank

migration_up:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migration_down:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migration_up migration_down sqlc test