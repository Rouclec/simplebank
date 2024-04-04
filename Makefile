DB_URL=postgresql://postgres:changemeinprod%21@localhost:5432/simple_bank?sslmode=disable

createdb:
	docker exec -it database-postgres_db_1 createdb --username=postgres --owner=postgres simple_bank

dropdb: 
	docker exec -it database-postgres_db_1 dropdb --username=postgres simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:changemeinprod%21@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	CompileDaemon --command="./simplebank"

mock: 
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/rouclec/simplebank/db/sqlc Store

.PHONY: createdb dropdb migrateup migratedown sqlc test server mock