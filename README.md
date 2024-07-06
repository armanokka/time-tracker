# Effective Mobile test task

### Documentation
Swagger documentation is located in `docs` folder

### Running
Edit `.env.dev` file according to your credentials. Then rename `.env.dev` to `.env`. Finally run:

    make migrate
    make init

### Requirements
- go (1.20+)
- make
- docker
- docker compose

### Swagger UI
http://localhost/swagger/index.html

### Jaeger UI
http://188.225.74.17:16686/

### Full list of what hass been used:
* [clean architecture](https://github.com/evrone/go-clean-template) - clean Architecture template for Golang services
* [gin-swagger](github.com/swaggo/gin-swagger) - Gin midddleware to generate Swagger docs
* [testcontainers-go](github.com/testcontainers/testcontainers-go) - test containers for integration testing
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [otlptracegrpc](go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc) - OTLP span exporter using gRPC
* [cleanenv](github.com/ilyakaznacheev/cleanenv) - environment configuration reader
* [go-redis](https://github.com/go-redis/redis) - type-safe Redis client for Golang
* [otel](https://go.opentelemetry.io/otel) -  Go implementation of OpenTelemetry
* [requestid](github.com/gin-contrib/requestid) - Request ID for Gin framework
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [jaeger](https://www.jaegertracing.io/) - open-source tracing platform
* [docker compose](https://docs.docker.com/compose/) - Docker compose
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [go-colorable](github.com/mattn/go-colorable) - colorful logging
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [gin](https://github.com/gin-gonic/gin) - web framework
* [docker](https://www.docker.com/) - Docker
* [swag](https://github.com/swaggo/swag) - Swagger
* [grpc](https://grpc.io/) - gRPC
* [gorm](https://gorm.io/) - ORM
* [zap](https://github.com/uber-go/zap) - logger

