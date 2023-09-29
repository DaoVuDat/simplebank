postgres:
	docker run --name postgres16alpine --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:16.0-alpine3.18

createdb:
	docker exec -it postgres16alpine createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16alpine dropdb simple_bank

migration_up:
	migrate -path ./db/migration/ -database "postgresql://root:co8ZZ1IinNVyXJ1jI4ri@simple-bank.ckwdemlqxoqw.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migration_up1:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migration_down:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

migration_down1:
	migrate -path ./db/migration/ -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migration_up migration_down sqlc test server mock