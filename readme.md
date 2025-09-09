# coreflow

## pre-requirements

> make sure that you have golang installed in the version that is specified in the `go.mod` file
>
> https://go.dev/doc/install

```sh
docker compose up -d
go mod tidy
```

## run the server

```sh
go run ./cmd/server/main.go
```

## database

### migrations

#### pre-requirements

```sh
export GOOSE_MIGRATION_DIR='./internal/db/migrations'
export GOOSE_DRIVER='postgres'
export GOOSE_DBSTRING='user=postgres password=postgres dbname=postgres host=localhost port=5432'
```

#### create a migration

```sh
go tool goose create -s INSERT_MIGRATION_NAME sql
```

#### run the migrations

```sh
go tool goose up
```

### seeding

#### pre-requirements

```sh
export GOOSE_MIGRATION_DIR='./internal/db/seeds'
export GOOSE_DRIVER='postgres'
export GOOSE_DBSTRING='user=postgres password=postgres dbname=postgres host=localhost port=5432'
```

#### create a seed

```sh
go tool goose create -s INSERT_SEED_NAME sql
```

#### run the seeds

> `-no-versioning` is fundamental here. you need to provide it to avoid the versioning of the seeds.

```sh
go tool goose -no-versioning up
```

### jet

#### generate the jet files

```sh
go tool jet -dsn='postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable' -schema='public' -path='./gen'
```
