.PHONY: migration-down migration-up db-start db-stop

MIGRATOR=go run ./cmd/migrator/main.go -host=localhost -port=5432 -login=kirill -password=pass123 -db=ozon -path=./migrations

migration-down:
	$(MIGRATOR) -mode=down

migration-up:
	$(MIGRATOR) -mode=up

db-start:
	docker run --rm --name pgdocker -e POSTGRES_PASSWORD=pass123 -e POSTGRES_USER=kirill -e POSTGRES_DB=ozon -d -p 5432:5432 -v $$HOME/docker/volumes/postgres:/var/lib/postgresql/data postgres

db-stop:
	docker stop pgdocker

start:
	go run cmd/server/main.go
