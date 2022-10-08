DB_URL=postgresql://root:root@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name bank_postgres \
		   --network bank_network \
           -p 5432:5432 \
           -e POSTGRES_USER=root \
           -e POSTGRES_PASSWORD=root \
           -d postgres:14-alpine

createdb:
	docker exec -it bank_postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -ti bank_postgres dropdb simple_bank

migrate-up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -coverprofile=coverage.out ./...

server:
	go build -o main main.go
	
mock:
	go generate -v ./...

migrate-up1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migrate-down1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt allow_merge=true --openapiv2_opt merge_file_name=simple_bank \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc test server mock migrate-up1 migrate-down1 db_docs db_schema proto evans
