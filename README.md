# Sports Event Calendar

This project is a solution for the backend version of Sportradar Coding Academy exercise.
It implements a simple sports event calendar.

## Tech stack
- Go
- PostgreSQL
- Docker

## Database design
The application models scheduled team-sport matches.
Reference data such as sports and countries is normalized into separate tables.

![ERD](docs/erd.png)

## Running the database

1. Copy `.env.example` to `.env`
2. Start PostgreSQL:
```bash
docker compose up -d
```


## Running migrations
Migrations are handled by a separate container.

Apply all pending migrations:

```bash
docker compose --profile utils run --rm migrations
```
This is equivalent to:
```bash
docker compose --profile utils run --rm migrations /migrate --up
```
Roll back all migrations:
```bash
docker compose --profile utils run --rm migrations /migrate --down --down-all 
```
Roll back a specific number of migrations:
```bash
docker compose --profile utils run --rm migrations /migrate --down --steps=[n]
```
If a migration fails and leaves the database in a dirty state,
fix the migration first, then force the migration version:
```bash
docker compose --profile utils run --rm migrations /migrate --force-version=[n]
```

## Seeding the database with exemplary data

Seed all files from /seeds
```bash
docker compose run --rm migrations /migrate --seed
```

Seed one file from given path
```bash
docker compose run --rm migrations /migrate --seed --seeds-path=[path] --seed-file-name=[file_name]
```