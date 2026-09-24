postgres:
	docker run --name postgres18 -p <post>:<port> -e POSTGRES_USER=<user> -e POSTGRES_PASSWORD=<password> -d postgres:18-alpine

createdb:
	docker exec -it postgres18 createdb --username=<username> --owner=<owner> <application_name>

dropdb:
	docker exec -it postgres18 dropdb <application_name>

migrateup:
	migrate -path db/migration -database "postgresql://<user>:<password>@localhost:5433/<application_name>?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://<user>:<password>@localhost:5433/<application_name>?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test