# Library management system

## Local infrastructure

Prerequisite: Docker Compose.

```sh
cp .env.example .env
docker compose -f docker-compose.yaml up -d --build
docker compose -f docker-compose.yaml ps
```

The checked-in value is a development-only Kong client secret matching the
bcrypt hash in the OAuth seed SQL. Do not reuse it outside local development.

Kong listens on `http://127.0.0.1:8000`. Requests under
`/api/v1/auth` are forwarded to the auth service under `/api/v1/auth`.
For example, `/api/v1/auth/login` remains `/api/v1/auth/login` upstream.
The auth service is also available directly on `http://127.0.0.1:8081`.
Its Swagger UI is available directly at
`http://127.0.0.1:8081/api/v1/docs/auth/swagger`.
Through Kong, Swagger UIs use these public routes:

- `http://127.0.0.1:8000/api/v1/docs/auth/swagger`
- `http://127.0.0.1:8000/api/v1/docs/users/swagger`
- `http://127.0.0.1:8000/api/v1/docs/books/swagger`
- `http://127.0.0.1:8000/api/v1/docs/transactions/swagger`

Kong preserves these paths. Each service must register its matching Swagger UI
and spec routes.

`/api/v1/users`, `/api/v1/books`, and `/api/v1/transactions` proxy to their
matching services. Kong validates bearer tokens through the auth service,
requiring `library:read` for reads and `library:write` for writes. These
upstreams remain unavailable until their containers join the gateway network.
Auth, token, health, metadata, and preflight routes have no authentication
plugin.

MySQL listens on `127.0.0.1:3306`. Local credentials are:

- Database: `auth`
- User: `auth`
- Password: `auth_password`
- Root password: `root_password`

The `auth_mysql_data` volume preserves data across normal restarts and
`docker compose -f docker-compose.yaml down`.

MySQL applies credentials only when the data volume is empty. After changing
hardcoded values, recreate the local volume intentionally.

Stop services:

```sh
docker compose -f docker-compose.yaml down
```

To also delete local database data:

```sh
docker compose -f docker-compose.yaml down --volumes
```
