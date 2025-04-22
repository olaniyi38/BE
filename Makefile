migrate_up:
	migrate -path db/migration/ -database "postgresql://sodiq:password@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate_up1:
	migrate -path db/migration/ -database "postgresql://sodiq:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migrate_down:
	migrate -path db/migration/ -database "postgresql://sodiq:password@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrate_down1:
	migrate -path db/migration/ -database "postgresql://sodiq:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	@sqlc generate
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/olaniyi38/BE/db/sqlc Store
test:
	go test -v -cover ./...
.PHONY: migrate_up migrate_down sqlc server migrate_down1 migrate_up1 test