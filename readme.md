# api

## pré-requisitos

```sh
cd apps/api
docker compose up -d
go mod tidy
```

## rodar o backend

```sh
go run ./cmd/server/main.go
```

# banco de dados

## migrations

### pré requisitos

```sh
export GOOSE_MIGRATION_DIR='./internal/db/migrations'
export GOOSE_DRIVER='postgres'
export GOOSE_DBSTRING='user=postgres password=postgres dbname=postgres host=localhost port=5432'
```

### criar uma migração

```sh
go tool goose create -s NOME_DA_MIGRATION sql
```

### subir as migrações

```sh
go tool goose up
```

## seeding

### pré-requisitos

```sh
export GOOSE_MIGRATION_DIR='./internal/db/seeds'
export GOOSE_DRIVER='postgres'
export GOOSE_DBSTRING='user=postgres password=postgres dbname=postgres host=localhost port=5432'
```

### criar uma seed

```sh
go tool goose create -s NOME_DA_SEED sql
```

### subir as seeds

```sh
go tool goose -no-versioning up
```

## sqlc

```sh
go tool sqlc generate
```
