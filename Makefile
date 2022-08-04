postgres:
	docker run --name postgres14 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=1234 -p 127.0.0.1:5432:5432 -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres lapcart -U postgres

dropdb:
	docker exec -it postgres14 dropdb lapcart -U postgres

goose_reset:
	goose -dir migrations postgres "user=postgres password=1234 dbname=lapcart sslmode=disable" reset

goose_up:
	goose -dir migrations postgres "user=postgres password=1234 dbname=lapcart sslmode=disable" up

goose_down:
	goose -dir migrations postgres "user=postgres password=1234 dbname=lapcart sslmode=disable" down

goose_status:
	goose -dir migrations postgres "user=postgres password=1234 dbname=lapcart sslmode=disable" status

.PHONY: postgres createdb dropdb goose_reset goose_up goose_status goose_down