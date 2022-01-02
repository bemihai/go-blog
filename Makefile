server:
	go run main.go 

postgres:
	docker run --name postgres12 --network blog-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres blog

dropdb:
	docker exec -it postgres12 dropdb blog

migrateinit:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/blog?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/blog?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: server postgres createdb migrateinit migrateup migratedown dropdb test