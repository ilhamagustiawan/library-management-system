# Book Service

Catalog, inventory, and atomic stock reservation service. Structure follows
`backend/auth-service` clean-layer conventions.

## Commands

```sh
cp .env.example .env
go run . migrate --action up
go run . serve
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

Local server listens on `:8083`. Compose overrides the container listener to
:8080` and publishes it on host port `8083`.

Swagger UI: `http://127.0.0.1:8083/api/v1/docs/books/swagger`.

## APIs

- `GET /api/v1/books` and `GET /api/v1/books/{id}` require `books:read` or
  `books:manage` from trusted Kong identity headers.
- Catalog mutations require `books:manage`.
- Internal stock routes require a Transaction Service bearer token with
  `aud=book-service`, matching client subject, and the matching `book-stock:*`
  scope.
- `PUT /internal/v1/books/{id}/reservations/{transactionId}` atomically checks
  availability and reserves one copy. Retrying an active reservation is safe.
- `DELETE` on the same reservation restores stock once and is retry-safe.

Transaction Service owns the maximum-three-active-loans rule. Borrow flow must
check that limit, generate the transaction ID, reserve stock, persist the loan,
then release the same reservation if persistence fails. A separate stock read
must not replace atomic reservation because check-then-reserve races.

Returned loans arrive through durable `LoanReturned.v1` RabbitMQ events. Book
Service releases each reservation exactly once, then publishes correlated
`BookStockUpdated.v1` acknowledgements through a transactional outbox.

## Catalog seed

Migrations `202607190003_seed_popular_books` and
`202607190004_seed_more_popular_books` seed 100 books from the Open Library
monthly-trending snapshot retrieved on 2026-07-19, including the top-ranked
**Atomic Habits**. Each book starts with five available copies and a cover URL. The
90 additional titles seeded by migration `004` have original multi-paragraph
library-staff descriptions; the initial ten retain their descriptive catalog copy.
Compatibility migrations `202607190005_backfill_popular_book_descriptions`,
`202607190006_correct_popular_book_descriptions`,
`202607190007_enrich_popular_book_descriptions`, and
`202607190008_correct_library_descriptions` upgrade earlier seed states without
overwriting a librarian-authored description. Cover URLs are returned as `coverUrl`
by the catalog API.

Open Library API source: https://openlibrary.org/trending/monthly.json

## Intentional v1 omissions
Archived-book restoration/purge, reservation retention cleanup, genre, and
publisher are deferred. Released reservation records remain stored to preserve retry
safety.
