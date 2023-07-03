postgres:
	docker pull postgres:latest
	docker rm -f postgres15
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:latest
createdb:
	docker exec -it postgres15 createdb --username=root --owner=root bank_simulator

dropdb:
	docker exec -it postgres15 dropdb bank_simulator

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank_simulator?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank_simulator?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
# run-in-sequence: postgres dropdb createdb migratedown migrateup

