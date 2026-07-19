# OAuth Service Implementation Plan

1. Replace the JWT-login draft with user/session domain contracts and failing use-case tests.
2. Implement bcrypt, secure opaque sessions, user/session PostgreSQL repositories, and migrations.
3. Implement PostgreSQL OAuth client/token stores satisfying `go-oauth2` interfaces.
4. Configure `go-oauth2` for Authorization Code + rotating refresh tokens, mandatory S256 PKCE, exact redirects, Basic client authentication, strict scopes, and Bearer-header access tokens.
5. Add Fiber registration/login/logout/me, OAuth adapters, metadata/userinfo, health endpoints, security middleware, and rate limiting.
6. Add configuration, Cobra serve/migrate/create-client commands, graceful shutdown, Docker, and Next.js integration documentation.
7. Run formatting, tests/race tests, vet, build, dependency audit, and a security-focused review.

## Risks

- `go-oauth2` defaults enable insecure/deprecated grants and permissive redirect matching; override both explicitly.
- Fiber is fasthttp-based while `go-oauth2` uses `net/http`; use Fiber's official adaptor only at OAuth protocol endpoints.
- Browser login cookies can fail in cross-site fetches; support top-level form POST + redirect and use a dedicated auth-domain cookie.
- Opaque token tables can grow; expiry filters and cleanup are included, with scheduled maintenance left for deployment.
