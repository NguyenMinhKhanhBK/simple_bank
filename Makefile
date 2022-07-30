postgres:
	docker run --name bank_postgres \
           -p 5432:5432 \
           -e POSTGRES_USER=root \
           -e POSTGRES_PASSWORD=root \
           -d postgres:14-alpine

createdb:
	docker exec -it bank_postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -ti bank_postgres dropdb simple_bank

migrate-up:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrate-up migrate-down
