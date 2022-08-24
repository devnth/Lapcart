

goose_reset:
	goose -dir migrations postgres "user=postgres password=1234 dbname=lapcart sslmode=disable" reset

goose_up:
	goose -dir migrations postgres "user=postgres password=1234 dbname=lapcart sslmode=disable" up

goose_down:
	goose -dir migrations postgres "user=postgres password=1234 dbname=lapcart sslmode=disable" down

goose_status:
	goose -dir migrations postgres "user=postgres password=1234 dbname=lapcart sslmode=disable" status

docker_build:
	docker compose build --no-cache

docker_up:
	docker compose up

docker_down:
	docker compose down

.PHONY: goose_reset goose_up goose_status goose_down docker__build docker_up docker_down