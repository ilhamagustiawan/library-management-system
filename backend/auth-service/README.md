# Library Management Auth Service

OAuth 2.0 authorization service for the library management system. It follows the Clean Architecture layout of [kittipat1413/ticket-reservation](https://github.com/kittipat1413/ticket-reservation), uses Fiber for HTTP delivery, and uses [go-oauth2/oauth2](https://github.com/go-oauth2/oauth2) for protocol handling.

## Security profile

- Authorization Code grant with refresh-token rotation.
- Client Credentials grants for provisioned internal services.
- Mandatory PKCE with `code_challenge_method=S256`.
- Confidential Next.js client authenticated with `client_secret_basic`.
- Exact redirect URI matching.
- Mandatory `state` on authorization requests.
- HS256 JWT access tokens and rotating opaque refresh tokens stored in MySQL.
- bcrypt cost 12 for user passwords and persisted OAuth client secrets.
- HttpOnly, SameSite=Lax session cookie.
- No implicit or password grants.

This is OAuth 2.0, not OpenID Connect. It does not issue ID tokens.

## Structure

```text
 cmd/                         Cobra commands
 db/migrations/              MySQL migrations
 db/seeds/                   Development-only idempotent seed SQL
 internal/api/http/          Handlers, request/response DTOs, helpers, middleware, and routes
 internal/domain/            Entities and repository contracts
 internal/infra/             bcrypt, MySQL, and OAuth adapters
 internal/usecase/           Auth, OAuth, and health business logic
 internal/server/            Dependency wiring and lifecycle
```

## Run locally

```shell
cp .env.example .env

docker compose -f ../../docker-compose.yaml up -d auth-db
go run . migrate --action up
go run . serve
```

Every command automatically loads `.env` from its current working directory. Already-exported environment variables take precedence.

Or run everything in Docker:

```shell
docker compose -f ../../docker-compose.yaml up --build auth-service
```

Health endpoints:

```shell
curl http://localhost:8081/health/liveness
curl http://localhost:8081/health/readiness
```

## API documentation

Start the service, then open the interactive Swagger UI:

```text
http://localhost:8081/api/v1/docs/auth/swagger
```

The generated Swagger 2.0 contract is available at:

```text
http://localhost:8081/api/v1/docs/auth/swagger.json
```

Handler annotations are the source of truth. After changing an endpoint,
regenerate and commit `docs/docs.go`, `docs/swagger.json`, and
`docs/swagger.yaml`:

```shell
make swagger
# Equivalent Cobra command:
go run . swagger
```

The Swagger UI supports bearer-token and HTTP Basic authentication. Session
endpoints use the `lms_session` HttpOnly cookie set by login.

## Member web OAuth client

In development, `migrate --action up` reapplies `db/seeds/oauth_clients.sql`
after schema migrations. It upserts local web, gateway, and service clients;
`serve` does not write seed data.

The first-party member web client is:

- Client ID: `member-nextjs-web`
- Client secret: `local-development-only-client-secret`
- Redirect URI: `http://localhost:3000/api/auth/callback/library`
- Scopes: `books:read loans:borrow:self loans:return:self transactions:read:self`

Configure the Next.js server with the matching local secret:

```dotenv
# Next.js server
AUTH_ISSUER=http://localhost:8081
AUTH_CLIENT_ID=member-nextjs-web
AUTH_CLIENT_SECRET=local-development-only-client-secret
```

Each development migration overwrites seeded rows with the SQL definitions,
including bcrypt secret hashes, while preserving `created_at`. Update the SQL
and matching consumer configuration together. Never reuse these credentials in
production or expose them through `NEXT_PUBLIC_*` variables.

## Create an additional OAuth client

The generic command remains available when another client is needed:

```shell
go run . create-client \
  --name "Admin portal" \
  --redirect-uri https://admin.example.com/api/auth/callback \
  --scopes "transactions:read:any loans:return:any fines:manage books:manage"
```

The command generates an ID when omitted and prints the new secret once. Use
`--kind client_credentials` with an empty redirect URI for an internal service,
or `--kind resource_server --scopes ""` for an introspection client.

## Create the first admin

Admin creation is offline only. It reads the password from the terminal,
assigns `admin`, and refuses to promote an existing member:

```shell
go run . create-admin --name "Grace Hopper" --email grace@example.com
```

## User login flow

Public registration belongs to User Service at `POST /api/v1/users`. It calls
Auth's `POST /internal/identities` with an idempotency key and User Service
token. Auth always assigns `member`; extra fields such as `role` are rejected.
Until User Service is available, frontend registration cannot complete.

When `/oauth/authorize` has no auth-service session, it redirects to:

```text
http://localhost:3000/login?return_to=<encoded authorization URL>
```

The Next.js login page should submit a top-level HTML form to the auth service. This lets the auth-service origin safely set its own HttpOnly cookie:

```tsx
<form action="http://localhost:8081/api/v1/auth/login" method="post">
  <input name="email" type="email" required />
  <input name="password" type="password" required />
  <input name="return_to" type="hidden" value={returnTo} />
  <button type="submit">Sign in</button>
</form>
```

`return_to` is accepted only when it points exactly to this service's `/oauth/authorize` endpoint.

## Next.js Authorization Code + PKCE flow

All steps below belong in server-only Route Handlers:

1. Generate a cryptographically random `state` and PKCE `code_verifier`.
2. Store both in an encrypted/HttpOnly Next.js session cookie.
3. Compute `code_challenge = BASE64URL(SHA256(code_verifier))`.
4. Redirect the browser to:

```text
GET http://localhost:8081/oauth/authorize
  ?response_type=code
  &client_id=member-nextjs-web
  &redirect_uri=http://localhost:3000/api/auth/callback/library
  &scope=books:read%20loans:borrow:self%20loans:return:self%20transactions:read:self
  &state=<state>
  &code_challenge=<challenge>
  &code_challenge_method=S256
```

5. In the callback, compare `state` with the one-time session value.
6. Exchange the code from the Next.js server:

```http
POST /oauth/token HTTP/1.1
Host: localhost:8081
Authorization: Basic base64(member-nextjs-web:client-secret)
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&code=<code>&redirect_uri=<exact-callback>&code_verifier=<verifier>
```

7. Keep access and refresh tokens in the server-side session, not browser-accessible storage.

Authorization-server metadata is available at:

```text
http://localhost:8081/.well-known/oauth-authorization-server
```

The current user can be loaded with:

```shell
curl http://localhost:8081/api/v1/oauth/userinfo \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## Token introspection and Kong

`POST /oauth/introspect` implements RFC 7662 for JWT access tokens. MySQL state,
not standalone signature validation, determines whether a token remains active.
Only the dedicated `resource_server` client seeded for Kong can call it. The
local secret is `local-development-only-introspection-secret`; Docker Compose
passes the matching value to Kong. Unknown, rotated, or expired tokens return
only `{"active":false}`.

Development seeds provision User Service with only `identities:create` for
`aud=auth-service`, and Transaction Service with only `book-stock:*` for
`aud=book-service`. Service tokens use client ID as `sub`, omit role, and have
no refresh token. Kong and Book Service are resource-server clients only.

## Commands

```shell
make run
make swagger
make migrate-up
make create-client NAME="Admin portal" REDIRECT_URI="https://admin.example.com/api/auth/callback"
go run . create-admin --name "Grace Hopper" --email grace@example.com
make test
make test-race
make precommit
```

## Production checklist

- Set `SERVICE_ENV=production`.
- Use HTTPS for `OAUTH_ISSUER`, `LOGIN_URL`, the Next.js origin, and callback URI.
- Set `SESSION_COOKIE_SECURE=true`.
- Set `OAUTH_JWT_SIGNING_KEY` to an environment-specific secret of at least 32 bytes.
- Replace local database credentials and require TLS to MySQL.
- Provision environment-specific Next.js and Kong clients before startup;
  development seed SQL is skipped when `SERVICE_ENV` is not `development`.
- Keep auth and Next.js origins on controlled domains.
- Add a consent screen before registering third-party clients.
- Add token revocation or OpenID Connect only as explicit follow-up features.
- Keep the introspection client secret separate from application client secrets.
- Keep the JWT signing key available only to the auth service while Kong uses introspection.
