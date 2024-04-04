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

docker-build:
	docker build -t simplebank:latest .

docker-run:
	 docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://postgres:changemeinprod%21@database-postgres_db_1:5432/simple_bank?sslmode=disable" simplebank:latest

docker-stop:
	docker rm simplebank

create-network:
	docker network create bank-network

connect-db-to-network:
	docker network connect bank-network database-postgres_db_1

.PHONY: createdb dropdb migrateup migratedown sqlc test server mock docker-build docker-run docker-stop create-network connect-db-to-network