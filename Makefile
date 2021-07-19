postgres:
	docker run --name postgres-12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres-12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres-12 dropdb simple_bank
migrateup:
	sql-migrate up
migratedown:
	sql-migrate down
sqlc:
	sqlc generate
.PHONY:	createdb