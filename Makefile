createdb:
	docker exec -it database-postgres_db_1 createdb --username=postgres --owner=postgres simple_bank

dropdb: 
	docker exec -it database-postgres_db_1 dropdb --username=postgres simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:changemeinprod%21@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:changemeinprod%21@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createdb dropdb migrateup migratedown sqlc test