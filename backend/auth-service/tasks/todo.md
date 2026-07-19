# OAuth Service Tasks

- [ ] User/session behavior
  - Acceptance: register normalizes identity, login hashes no plaintext and creates an opaque session, me resolves it, logout revokes it.
  - Verify: auth use-case tests pass.
- [ ] PostgreSQL persistence
  - Acceptance: reversible migrations create users, sessions, OAuth clients, and OAuth token storage; all queries are parameterized.
  - Verify: repository packages compile and tests pass.
- [ ] OAuth protocol
  - Acceptance: only code/refresh grants work; state and S256 PKCE are mandatory; redirect matching is exact; client secrets are hashed.
  - Verify: OAuth protocol tests complete a valid code exchange and reject insecure variants.
- [ ] Fiber delivery
  - Acceptance: application, OAuth, metadata, userinfo, and health routes are wired with cookies, validation, security headers, origin checks, and limits.
  - Verify: Fiber HTTP tests pass.
- [ ] Runtime scaffold
  - Acceptance: serve/migrate/create-client commands, config, graceful shutdown, Docker, Makefile, and README work together.
  - Verify: `go test -race ./...`, `go vet ./...`, and `go build ./...` pass.
- [ ] Final review
  - Acceptance: no Gin/JWT-only draft, no secrets, no deprecated grants, and no permissive redirects.
  - Verify: inspect diff, audit dependencies, and complete security review.
