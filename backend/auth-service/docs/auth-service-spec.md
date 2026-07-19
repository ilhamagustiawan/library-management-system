# Spec: Library Management OAuth 2.0 Service

## Objective

Scaffold an OAuth 2.0 authorization service for a Next.js library-management client. The service follows the Clean Architecture conventions from `kittipat1413/ticket-reservation`, uses Fiber for HTTP delivery, and delegates protocol mechanics to `github.com/go-oauth2/oauth2/v4`.

The first release supports:

- Login, logout, session lookup, and internal member identity creation.
- OAuth 2.0 Authorization Code grant with mandatory PKCE using S256.
- OAuth 2.0 Client Credentials grants for provisioned internal services.
- Confidential-client authentication for the server-side Next.js code exchange.
- Refresh-token rotation.
- HS256 JWT access tokens, opaque refresh tokens, and a protected current-user endpoint.
- MySQL-backed users, sessions, OAuth clients, authorization codes, and tokens.
- Authorization-server metadata and health endpoints.
- RFC 7662 JWT state introspection for authenticated resource servers.

Implicit, Resource Owner Password Credentials, dynamic client registration, social login, password reset, email verification, OpenID Connect ID tokens, public role management, and third-party consent screens are out of scope.

## Tech Stack

- Go 1.23
- Fiber v2.52.9
- `go-oauth2/oauth2/v4` v4.5.4
- `golang-jwt/jwt/v5` v5.3.0
- MySQL via `sqlx` v1.4.0 and `go-sql-driver/mysql` v1.8.1
- bcrypt password and client-secret hashing with cost 12
- Cobra v1.9.1 commands
- `golang-migrate` v4.18.2 migrations
- `go-playground/validator/v10` request validation

## API Contract

### User/session endpoints

- `POST /internal/identities` idempotently creates a member identity for authenticated User Service.
- `POST /api/v1/auth/login` accepts JSON or an HTML form, creates an HttpOnly session cookie, and either returns `200` or redirects to a validated local authorization URL.
- `POST /api/v1/auth/logout` revokes the current session and clears the cookie.
- `GET /api/v1/auth/me` returns the session user.

### OAuth endpoints

- `GET /oauth/authorize` requires `response_type=code`, non-empty `state`, `code_challenge`, and `code_challenge_method=S256`.
- `POST /oauth/token` supports `authorization_code`, `client_credentials`, and `refresh_token`, and authenticates confidential clients with HTTP Basic authentication.
- `POST /oauth/introspect` requires a dedicated resource-server client and returns RFC 7662 token state.
- `GET /.well-known/oauth-authorization-server` advertises authorization, token, grant, response, and PKCE metadata.
- `GET /api/v1/oauth/userinfo` requires an access token in the `Authorization: Bearer` header and returns the public user.

Access tokens are signed with HS256 and contain `iss`, `aud`, `sub`,
`client_id`, `scope`, `iat`, `exp`, and `jti`. User tokens also contain `role`;
service tokens use client ID as `sub` and omit `role`. Audience derives from
requested scopes; mixed-audience requests fail. Refresh tokens remain opaque
32-byte random values and are issued only for authorization-code flows. Tokens
stay MySQL-backed: introspection uses persisted state, and refresh rotation
immediately makes the prior access token inactive.

`OAUTH_JWT_SIGNING_KEY` is required when serving HTTP and must contain at least
32 bytes. Administrative commands remain usable without the runtime signing
key.

OAuth protocol errors use RFC 6749 response fields. Application endpoints follow the reference project response shape:

```json
{"code":"LMS-200000","data":{}}
```

## Commands

```shell
go run . serve
go run . migrate --action up
go run . migrate --action down
go run . create-client --name "Admin portal" --redirect-uri https://admin.example.com/api/auth/callback
go run . create-admin --name "Grace Hopper" --email grace@example.com
go test ./...
go test -race ./...
go vet ./...
go build ./...
docker compose -f ../../docker-compose.yaml up --build auth-service
```

## Project Structure

```text
 cmd/                              Cobra entry points
 db/migrations/                   Reversible MySQL migrations
 db/seeds/                        Development-only idempotent seed SQL
 internal/api/http/               Handlers, request/response DTOs, helpers, middleware, routes
 internal/config/                 Environment configuration
 internal/domain/                 Entities, errors, repository contracts
 internal/infra/auth/             bcrypt/session primitives
 internal/infra/db/repository/    MySQL repositories and OAuth stores
 internal/infra/oauth/            go-oauth2 server configuration
 internal/server/                 Dependency wiring and lifecycle
 internal/usecase/                Auth, OAuth, and health business rules
```

## Code Style

Use constructor-injected interfaces and context-aware operations. HTTP, business, and persistence concerns remain in separate layers. All Go code is formatted with `gofmt`.

## Testing Strategy

- Unit tests for internal identity creation, login, session authentication, and logout.
- OAuth protocol tests proving PKCE is required, only S256 is accepted, redirects match exactly, and valid codes can be exchanged.
- Token tests proving HS256 signatures, resource claims, opaque refresh tokens, and refresh invalidation.
- Fiber HTTP tests for strict validation, cookies, protected routes, and security headers.
- MySQL is not required for unit tests; database adapters are covered with focused tests where useful and runtime integration through Docker.

## Boundaries

- Always: exact redirect URI matching, mandatory state and S256 PKCE, hashed user passwords/client secrets, generic login errors, parameterized SQL, short-lived authorization codes, HS256 access tokens with a 32-byte minimum signing key, rotating refresh tokens, HttpOnly/SameSite cookies, strict origin checks, auth rate limits, and access tokens accepted only from Bearer headers.
- Ask first: add grants, OIDC, third-party clients/consent, social login, Redis, token revocation, or new PII.
- Never: commit production secrets, log credentials/codes/tokens, distribute the JWT signing key to clients, accept wildcard redirect URIs/origins, expose secrets to `NEXT_PUBLIC_*`, or enable implicit/password grants.

## Success Criteria

- Authorization Code + PKCE S256 completes end to end with a confidential client.
- Requests without PKCE, with `plain`, without state, or with a non-exact redirect URI fail.
- The OAuth server allows authorization-code, client-credentials, and refresh-token grants only.
- Client Credentials tokens can be issued only to eligible clients and never include refresh tokens.
- User scopes equal requested intersection role and client grants; service scopes equal requested intersection client grants.
- Unauthorized or mixed-audience scopes fail with `invalid_scope`.
- Public registration cannot assign roles; offline admin creation cannot promote an existing member.
- Issued access tokens are valid HS256 JWTs with the required resource claims.
- Refresh tokens remain opaque, rotate on use, and invalidate their prior access tokens.
- Introspection discloses active-token claims only to the resource-server client.
- Passwords and client secrets are bcrypt hashes and never serialized.
- Sessions, registered clients, and OAuth artifacts persist in MySQL and expire where applicable.
- Fiber applies recovery, security headers, bounded bodies, strict origin handling, and auth rate limits.
- Tests, race tests, vet, and build pass.

## Open Questions

- Development migrations upsert the first-party member web and Kong resource-server clients from SQL; service startup never seeds data. Production client provisioning remains deployment-owned.
- This is OAuth 2.0, not OpenID Connect; add OIDC deliberately if Next.js needs ID tokens or standard OIDC discovery/userinfo semantics.
- Kong intentionally keeps RFC 7662 introspection; JWKS and local JWT validation remain separate future work.
