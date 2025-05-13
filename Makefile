postgres:
	docker run -d --name octo_quiz -p 5465:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:16-alpine

createdb:
	docker exec -it octo_quiz createdb --username=root --owner=root octo_quiz
dropdb:
	docker exec -it octo_quiz dropdb octo_quiz

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5465/octo_quiz?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5465/octo_quiz?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5465/octo_quiz?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5465/octo_quiz?sslmode=disable" -verbose down 1
migrateforce:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5465/octo_quiz?sslmode=disable" force 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...	

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ulugbek0217/octo-quiz/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock