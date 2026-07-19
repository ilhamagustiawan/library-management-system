# Library Management Auth Service

OAuth 2.0 authorization service for the library management system. It follows the Clean Architecture layout of [kittipat1413/ticket-reservation](https://github.com/kittipat1413/ticket-reservation), uses Fiber for HTTP delivery, and uses [go-oauth2/oauth2](https://github.com/go-oauth2/oauth2) for protocol handling.

## Security profile

- Authorization Code grant with refresh-token rotation.
- Client Credentials protocol support, with no machine clients provisioned.
- Mandatory PKCE with `code_challenge_method=S256`.
- Confidential Next.js client authenticated with `client_secret_basic`.
- Exact redirect URI matching.
- Mandatory `state` on authorization requests.
- Opaque access/refresh tokens stored in MySQL.
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
after schema migrations. It upserts both local infrastructure clients; `serve`
does not write seed data.

The first-party member web client is:

- Client ID: `member-nextjs-web`
- Client secret: `local-development-only-client-secret`
- Redirect URI: `http://localhost:3000/api/auth/callback/library`
- Scopes: `library:read library:write`

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
  --scopes "library:read"
```

The command generates an ID when omitted and prints the new secret once.

## User login flow

Register a user:

```shell
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ada Lovelace","email":"ada@example.com","password":"correct horse battery staple"}'
```

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
  &scope=library:read
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

`POST /oauth/introspect` implements RFC 7662 for opaque access tokens. Only
the dedicated `resource_server` client seeded for Kong can call it. The local
secret is `local-development-only-introspection-secret`; Docker Compose passes
the matching value to Kong. Unknown or expired tokens return only
`{"active":false}`.

Client Credentials is enabled in the authorization server but intentionally
has no provisioning command or seeded client. Add onboarding and secret
delivery as separate work before enabling machine access.

## Commands

```shell
make run
make swagger
make migrate-up
make create-client NAME="Admin portal" REDIRECT_URI="https://admin.example.com/api/auth/callback"
make test
make test-race
make precommit
```

## Production checklist

- Set `SERVICE_ENV=production`.
- Use HTTPS for `OAUTH_ISSUER`, `LOGIN_URL`, the Next.js origin, and callback URI.
- Set `SESSION_COOKIE_SECURE=true`.
- Replace local database credentials and require TLS to MySQL.
- Provision environment-specific Next.js and Kong clients before startup;
  development seed SQL is skipped when `SERVICE_ENV` is not `development`.
- Keep auth and Next.js origins on controlled domains.
- Add a consent screen before registering third-party clients.
- Add token revocation or OpenID Connect only as explicit follow-up features.
- Keep the introspection client secret separate from application client secrets.
