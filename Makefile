# DB_URL=postgresql://postgres:changemeinprod%21@localhost:5432/simple_bank?sslmode=disable
DB_PASSWORD=password
DB_URL=postgres://postgres:password@localhost:5432/simple_bank?sslmode=disable&connect_timeout=5

# createdb:
# 	docker exec -it database-postgres_db_1 createdb --username=postgres --owner=postgres simple_bank

# dropdb: 
# 	docker exec -it database-postgres_db_1 dropdb --username=postgres simple_bank

# migrateup:
# 	migrate -path db/migration -database "postgresql://postgres:changemeinprod%21@localhost:5432/simple_bank?sslmode=disable" -verbose up

# migratedown:
# 	migrate -path db/migration -database "$(DB_URL)" -verbose down

getpostgres:
	docker pull postgres:16-alpine

runpostgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=simple_bank -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres16 dropdb --username=postgres simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

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

generate-random-32bit-hex:
	openssl rand -hex 64 | head -c 32

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

evans:
	docker run --net=host --rm -it -v "$(shell pwd):/mount:ro" ghcr.io/ktr0731/evans:latest --path ./proto/ --proto /service_simple_bank.proto --host 0.0.0.0 --port 9090 repl

.PHONY: createdb dropdb migrateup migratedown sqlc test server mock docker-build docker-run docker-stop create-network connect-db-to-network proto evans
