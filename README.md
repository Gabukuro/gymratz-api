# GymRatz-API 💪

## Context 📖

The GymRatz-API is a personal project designed to track and measure exercise progress. It aims to provide users with a comprehensive tool to monitor their fitness journey.

## Built with 📦

The application is primarily written in Golang (version 1.23.3) and utilizes the following libraries:

- **Fiber**: `github.com/gofiber/fiber/v2 v2.52.6`
- **JWT**: `github.com/golang-jwt/jwt v3.2.2+incompatible`
- **UUID**: `github.com/google/uuid v1.6.0`
- **Godotenv**: `github.com/joho/godotenv v1.5.1`
- **PostgreSQL Driver**: `github.com/lib/pq v1.10.7`
- **SQL Migrate**: `github.com/rubenv/sql-migrate v1.7.1`
- **Testify**: `github.com/stretchr/testify v1.10.0`
- **Testcontainers**: `github.com/testcontainers/testcontainers-go v0.35.0`
- **Bun ORM**:
  - `github.com/uptrace/bun v1.2.9`
  - `github.com/uptrace/bun/dialect/pgdialect v1.2.9`
  - `github.com/uptrace/bun/driver/pgdriver v1.2.9`
  - `github.com/uptrace/bun/extra/bundebug v1.2.9`
- **Crypto**: `golang.org/x/crypto v0.32.0`

## How to 🚀

### Install the project

To install the project dependencies, run:

```sh
make init
```

### Run the project locally

To run the project locally, use:

```sh
# build and run the database containers
make docker-up

# run server
make run
```

### Access the database

You can access the database using Adminer at: [http://localhost:8080/](http://localhost:8080/?pgsql=db&username=postgres&db=gymratz-api&ns=public)

Use the following credentials to log in:
- **System**: PostgreSQL
- **Server**: db
- **Username**: postgres
- **Password**: postgres
- **Database**: gymratz-api

### Run the tests

To execute the tests, run:

```sh
make test
```