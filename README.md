# xbankAPI

An Simple bank API


## Run using docker-compose

### prerequisites

1. [docker compose](https://docs.docker.com/compose/install/)


```bash
docker compose up
```

To verify service is ready.
```bash
curl http://localhost:3000/readiness
```

## Development

### prerequisites

1. [docker compose](https://docs.docker.com/compose/install/)

2. [golangci-lint](https://golangci-lint.run/welcome/install/)


### Set Enviroment Variables
```bash
# Required 
# Example: export DATABASE_URL="postgres://admin:admin@localhost:5432/bank_service?sslmode=disable"
export DATABASE_URL=

# Optional
# Example: export SERVCE_ADDRESS=":4002"
export SERVCE_ADDRESS=
```

### Setup Database
run *postgresql* db locally using docker-compose by running the following command:
```bash
docker compose up -d postgres 
```

### Migration

#### Run migrations
```bash
make migrate-up
```

#### Rollback migrations
```bash
make migrate-down
```

#### Create migration
```bash
make migrate-create name=[ANY-NAME-YOU-WANT]
```

### How to Generate SQLC and Mockery
```bash
make generate
```

## Continuous Integration

### Vet
```bash
make vet
```

### Lint
```bash
make lint
```

### test
```bash
make test
```

### build
```bash
make build
```
