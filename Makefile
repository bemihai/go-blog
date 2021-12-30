postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres blog

dropdb:
	docker exec -it postgres12 dropdb blog

migrateup:
	migrate -path sql/migration -database "postgresql://postgres:postgres@localhost:5432/blog?sslmode=disable" -verbose up

migratedown:
	migrate -path sql/migration -database "postgresql://postgres:postgres@localhost:5432/blog?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb migrateup migratedown dropdb test