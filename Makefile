init: build run

migrate_up:
	migrate -database "postgresql://user:password@localhost:5432/database?sslmode=disable" -path ./migrations/ up 1

migrate_down:
	migrate -database "postgresql://user:password@localhost:5432/database?sslmode=disable" -path ./migrations/ down 1


build:
	GOOS=linux GOARCH=amd64 go build -o effectiveMobile ./cmd/api
	docker-compose build

run:
	docker-compose up --remove-orphans --attach backend

stop:
	docker-compose stop
	