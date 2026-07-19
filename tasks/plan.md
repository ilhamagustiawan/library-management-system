# Authorization Plan: Roles and OAuth Scopes

## Outcome

Replace coarse `library:read`/`library:write` scopes with role-derived domain scopes. Members access only their loans and transactions. Admins act as librarians and can list every member transaction. Transaction Service uses a separate client-credentials identity for Book Service.

## Decisions

- Human roles: `member`, `admin`.
- Admin cannot borrow; admin receives no `loans:borrow:*` scope.
- Member can browse the book catalog through `books:read`.
- One role per user: `users.role_code`; no `user_roles` join table.
- Roles grant scopes through `role_scopes`.
- OAuth clients receive a separate ceiling through `oauth_client_scopes`.
- Authorization-code grant uses role scopes and client scopes.
- Client-credentials grant uses client scopes only; service tokens have no role.
- APIs authorize scopes. `role` remains identity context. `sub` enforces ownership.
- Reject unauthorized requested scopes with `invalid_scope`; never silently elevate.
- User Service owns public registration through `POST /api/v1/users`.
- Auth Service owns credentials/RBAC and creates the canonical user ID.
- User Service synchronously sends an idempotent `CreateIdentity` command to Auth, then publishes `UserRegistered` after successful local persistence.
- Admin creation is offline through bootstrap CLI only; no role-management API or `users:manage` scope.
- Preserve ongoing JWT work; extend it after its current migration/generator changes settle.

## Scope Catalog

| Scope | Audience | Principal | Use |
|---|---|---|---|
| `loans:borrow:self` | `library-api` | member | Borrow for token subject |
| `loans:return:self` | `library-api` | member | Return token subject's loan |
| `transactions:read:self` | `library-api` | member | Read token subject's history |
| `books:read` | `library-api` | member | Browse/search catalog and availability |
| `transactions:read:any` | `library-api` | admin | List/filter all member transactions |
| `loans:return:any` | `library-api` | admin | Librarian-assisted return |
| `fines:manage` | `library-api` | admin | View/update fine records |
| `books:manage` | `library-api` | admin | Manage catalog/inventory |
| `identities:create` | `auth-service` | user-service | Create member identity during registration |
| `book-stock:read` | `book-service` | transaction-service | Check stock |
| `book-stock:reserve` | `book-service` | transaction-service | Atomically reserve stock |
| `book-stock:release` | `book-service` | transaction-service | Restore stock on return/compensation |

## Data Model

```sql
CREATE TABLE roles (
    code VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    description VARCHAR(255) NOT NULL
);

CREATE TABLE scopes (
    code VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    audience VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    description VARCHAR(255) NOT NULL
);

CREATE TABLE role_scopes (
    role_code VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    scope_code VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    PRIMARY KEY (role_code, scope_code),
    FOREIGN KEY (role_code) REFERENCES roles(code),
    FOREIGN KEY (scope_code) REFERENCES scopes(code)
);

ALTER TABLE users
    ADD COLUMN role_code VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin
        NOT NULL DEFAULT 'member',
    ADD CONSTRAINT users_role_fk FOREIGN KEY (role_code) REFERENCES roles(code);

CREATE TABLE oauth_client_scopes (
    client_id VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    scope_code VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    PRIMARY KEY (client_id, scope_code),
    FOREIGN KEY (client_id) REFERENCES oauth_clients(id) ON DELETE CASCADE,
    FOREIGN KEY (scope_code) REFERENCES scopes(code)
);
```

Migration seeds roles, scopes, role mappings, and client mappings. Backfill `oauth_client_scopes` from `oauth_clients.allowed_scopes`; remove the text column only after repository cutover.

## Registration Contract

`UserRegistered` is an event: a fact published after registration commits. It must not be used as a synchronous request expecting Auth to return an ID.

```text
POST /api/v1/users
  -> User Service validates input and records idempotency key
  -> sync CreateIdentity command to Auth Service
  -> Auth creates member credentials and returns canonical user ID
  -> User Service persists registration using returned ID
  -> User Service outbox publishes UserRegistered.v1
  -> return 201
```

The Auth command carries the password over protected service-to-service transport. The event never contains password, password hash, token, or client secret.

```json
{
  "eventId": "evt-123",
  "type": "UserRegistered.v1",
  "occurredAt": "2026-07-19T10:00:00Z",
  "data": {
    "userId": "user-123",
    "role": "member"
  }
}
```

Both command and event consumers use `eventId`/idempotency key to make retries safe. Transactional outbox prevents lost events after database commit.

## Grant Rules

Authorization Code:

```text
granted = requested scopes
          intersect role scopes
          intersect OAuth client scopes
          filtered to requested resource audience
```

Client Credentials:

```text
granted = requested scopes
          intersect OAuth client scopes
          filtered to requested resource audience
```

This rule covers both internal clients: User Service receives only `identities:create`; Transaction Service receives only `book-stock:*`.

User token:

```json
{
  "sub": "user-123",
  "role": "member",
  "aud": ["library-api"],
  "scope": "loans:borrow:self loans:return:self transactions:read:self books:read"
}
```

Service token:

```json
{
  "sub": "transaction-service",
  "client_id": "transaction-service",
  "aud": ["book-service"],
  "scope": "book-stock:read book-stock:reserve book-stock:release"
}
```

## Endpoint Policies

| Endpoint | Required scope | Resource rule |
|---|---|---|
| `POST /api/v1/users` | Public, rate-limited | User Service ignores/rejects `role`; new identity is `member` |
| `POST /internal/identities` | `identities:create` | Auth Service accepts only User Service token; assigns `member` |
| `GET /books` | `books:read` | Member can browse/search; paginated |
| `POST /loans` | `loans:borrow:self` | Member ID always token `sub` |
| `POST /loans/{loanId}/return` | `loans:return:self` | Loan owner must equal token `sub` |
| `GET /transactions/me` | `transactions:read:self` | Query only `member_id = sub` |
| `GET /admin/transactions` | `transactions:read:any` | Paginated; optional member filter |
| `POST /admin/loans/{loanId}/return` | `loans:return:any` | Admin-assisted return |
| Book stock internal endpoints | matching `book-stock:*` | Require `aud=book-service`, `sub=transaction-service` |

Gateway validates human tokens with `aud=library-api`. Each resource service enforces endpoint scope and resource ownership. Internal Book endpoints require `aud=book-service`. HTTP method alone must not determine authorization.

## Implementation Tasks

### Task 1: Add RBAC schema and seed catalog

**Acceptance criteria:**

- Reversible migration creates `roles`, `scopes`, `role_scopes`, and `users.role_code`.
- Existing users become `member`; identity creation cannot choose a role.
- Seed data exactly matches catalog above.

**Verification:** migration up/down on fresh and populated test databases; repository migration tests.

**Dependencies:** None. **Size:** M.

### Task 2: Normalize OAuth client scope grants

**Acceptance criteria:**

- Add `oauth_client_scopes`; backfill current client grants.
- Repository reads normalized mappings, not space-delimited policy text.
- Unknown client/scope mappings fail closed.

**Verification:** client-store tests for allowed, denied, empty, and unknown scopes.

**Dependencies:** Task 1. **Size:** M.

### Task 3: Add typed role and scope policy domain

**Acceptance criteria:**

- Domain permits only `member` and `admin` roles.
- User repository loads role on every authentication path.
- Scope resolver returns structured `invalid_scope` failures.

**Verification:** domain and repository unit tests; `go test ./internal/domain/... ./internal/infra/db/repository/...`.

**Dependencies:** Tasks 1-2. **Size:** M.

### Task 4: Add admin bootstrap CLI

**Acceptance criteria:**

- Offline CLI creates an admin with securely entered password; no password command-line argument or logging.
- Command is idempotent by normalized email and cannot promote an existing member silently.
- No public/admin API can assign or change roles.

**Verification:** CLI tests for create, duplicate retry, existing-member conflict, invalid input, and secret redaction.

**Dependencies:** Tasks 1 and 3. **Size:** S.

### Task 5: Move public registration to User Service

**Acceptance criteria:**

- `POST /api/v1/users` validates registration and exposes no `role` field.
- User Service sends idempotent `CreateIdentity`; Auth assigns `member` and returns canonical user ID.
- Existing public Auth registration route is removed or made internal without breaking login.

**Verification:** contract tests for success, submitted-role rejection, retry after timeout, duplicate email, invalid service token, and rate limit; outbox test proves one `UserRegistered.v1` event without secrets.

**Dependencies:** Tasks 1-3. **Size:** M.

### Task 6: Enforce role-derived user grants

**Acceptance criteria:**

- Authorization Code tokens use requested ∩ role ∩ client ∩ audience scopes.
- Member cannot receive `transactions:read:any`, even when requested.
- Member receives `books:read`; admin cannot receive `loans:borrow:self`.

**Verification:** OAuth tests for member success, escalation rejection, and admin success.

**Dependencies:** Tasks 3 and 5. **Size:** M.

### Checkpoint: Human authorization

- Migrations reversible.
- Existing registration/login flow passes.
- Member/admin scope issuance proven.

### Task 7: Provision internal service clients

**Acceptance criteria:**

- Transaction Service receives only `book-stock:*`; User Service receives only `identities:create`.
- Service tokens use client subject, target service audience, no role, and no refresh token.
- Human and internal clients cannot request each other's scopes.

**Verification:** client-credentials protocol tests and seed tests.

**Dependencies:** Tasks 2-3. **Size:** M.

### Task 8: Extend token claims and introspection

**Acceptance criteria:**

- User tokens expose trusted `role`; service tokens omit it.
- Token/introspection returns exact granted scopes and correct audience/subject.
- Bootstrap role appears in new tokens; refresh cannot gain broader scopes.

**Verification:** JWT claim tests, introspection tests, refresh-scope regression tests.

**Dependencies:** Tasks 6-7 and current JWT work. **Size:** M.

### Task 9: Replace gateway coarse-scope policy

**Acceptance criteria:**

- Remove `library:read`/`library:write` policy and seeds.
- Gateway validates authentication; resource services enforce endpoint scopes.
- Wrong audience, missing scope, and subjectless user token fail with `401`/`403` consistently.

**Verification:** Kong plugin tests plus route-level authorization integration tests.

**Dependencies:** Task 8. **Size:** M.

### Task 10: Implement transaction access policies

**Acceptance criteria:**

- Member history query derives identity from `sub`; request cannot override it.
- Admin list returns paginated all-member transactions and supports member filter.
- Borrow/return endpoints enforce matching `self`/`any` scopes before domain logic.

**Verification:** handler/use-case tests for ownership, cross-member denial, admin access, pagination.

**Dependencies:** Tasks 8-9 and Transaction Service scaffold. **Size:** M.

### Task 11: Document and verify end to end

**Acceptance criteria:**

- OAuth metadata, Swagger, event schema, env examples, and architecture docs show final contracts.
- Registration, member catalog/history/loan, and admin history flows pass end to end.
- Internal client-credentials calls pass; human tokens fail against internal endpoints.

**Verification:** service tests, race tests, vet, build, Docker Compose smoke flow.

**Dependencies:** Tasks 1-10. **Size:** M.

## Risks

| Risk | Mitigation |
|---|---|
| Role/client mapping drift | One seeded catalog; foreign keys; startup/migration checks |
| Registration succeeds in Auth but User persistence fails | Idempotent command plus retryable pending registration |
| Database commits but event publish fails | Transactional outbox and idempotent consumers |
| Gateway-only authorization bypass | Resource services re-check scopes and ownership |
| Service token used at wrong API | Audience validation and network isolation |
| Existing JWT edits conflict | Land/rebase current JWT work before Tasks 6 and 8 |

## Unresolved Questions

None.
