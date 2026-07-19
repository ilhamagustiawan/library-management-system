# Library management system

## Local infrastructure

Prerequisite: Docker Compose.

```sh
docker compose -f docker-compose.yaml up -d --build
docker compose -f docker-compose.yaml ps
```

The checked-in values are development-only OAuth fixtures. Do not reuse them
elsewhere.

Kong listens on `http://127.0.0.1:8000`. Requests under
`/api/v1/auth` are forwarded to the auth service under `/api/v1/auth`.
For example, `/api/v1/auth/login` remains `/api/v1/auth/login` upstream.
The auth service is also available directly on `http://127.0.0.1:8081`.
The book service is available directly on `http://127.0.0.1:8083`.
The user service is available directly on `http://127.0.0.1:8082`.
The transaction service is available directly on `http://127.0.0.1:8084`.
Direct Swagger UIs use each service port and matching `/api/v1/docs/*/swagger`
path.
Through Kong, Swagger UIs use these public routes:

- `http://127.0.0.1:8000/api/v1/docs/auth/swagger`
- `http://127.0.0.1:8000/api/v1/docs/users/swagger`
- `http://127.0.0.1:8000/api/v1/docs/books/swagger`
- `http://127.0.0.1:8000/api/v1/docs/transactions/swagger`

Kong preserves these paths. Each service must register its matching Swagger UI
and spec routes.

`/api/v1/users`, `/api/v1/books`, and transaction routes proxy to matching
services. Kong validates token audience, subject, role, and exact endpoint
scope. Book Service rechecks `books:read`/`books:manage` using gateway
identity headers. Internal stock routes are not exposed through Kong. Auth,
token, health, metadata, public user creation, and preflight routes have no
authentication plugin. Auth Service no longer exposes public registration.

Auth MySQL listens on `127.0.0.1:3306`. Book MySQL listens on `127.0.0.1:3308`.
User MySQL listens on `127.0.0.1:3307`. Transaction MySQL listens on
`127.0.0.1:3309`. RabbitMQ listens on `127.0.0.1:5672`; its management UI uses
`http://127.0.0.1:15672`.
Local Book credentials: database/user `book`, password `book_password`.
Local User credentials: database/user `users`, password `users_password`.
Local Transaction credentials: database/user `transactions`, password
`transactions_password`.
Local Auth credentials are:

- Database: `auth`
- User: `auth`
- Password: `auth_password`
- Root password: `root_password`

The `auth_mysql_data`, `book_mysql_data`, `user_mysql_data`,
`transaction_mysql_data`, and `rabbitmq_data` volumes preserve data across normal restarts and
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
