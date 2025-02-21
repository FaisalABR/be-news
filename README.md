# Golang BE News

A REST API backend is developed for news management, where an admin can publish news created by users. It is built using Hexagonal Architecture, which separates business logic from external technologies.

# Techstack

- Docker
- Golang LTS v1.23
- PostgreSQL LTS
- Supabase
- CloudFlare
- Swagger

# How to Run the Project

- Clone this respository

```bash
git clone https://github.com/FaisalABR/be-news.git
```

- Change to project directory

```bash
cd be-news
```

- Install dependencies

```bash
go mod tidy
```

- Create `.env` file (check .env.example for the variables)

- Run the project

```bash
go run main.go
```

# API Documentation

For checking API Documentation that created using swagger, check url

```bash
localhost:8080/api/docs
```

# Database

Create migration:

```bash
migrate create -ext sql -dir database/migrations -seq create_users_table

#it will generate two file in database/migrations directories create_users_table.up.sql and create_users_table.down.sql
## create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);

## create_users_table.down.up
DROP TABLE IF EXISTS users;

```

Run migration:

- Up migration

```bash
migrate -database {POSTGRESQL_URL} -path database/migrations up
```

for `POSTGRES_URL` structure check the `.env` variable:

```bash
#POSTGRES_URL structure
postgress://{DATABASE_USERNAME}:{DATABASE_PASSWORD}@{DATABASE_HOST}:{DATABASE_PORT}/{DATABASE_NAME}
```

- Down migration

```bash
migrate -database {POSTGRESQL_URL} -path database/migrations down
```
