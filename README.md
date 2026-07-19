# Library management system

## Local infrastructure

Prerequisites: Docker Compose and an auth service listening on
`0.0.0.0:8080`.

```sh
docker compose -f docker-compose.yaml up -d
docker compose -f docker-compose.yaml ps
```

Kong listens on `http://127.0.0.1:8000`. Requests under
`/api/auth` are forwarded to the host auth service with that prefix removed.
For example, `/api/auth/login` becomes `/login` upstream.

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
