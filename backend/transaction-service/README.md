# Library Management Transaction Service

Owns loans, returns, fines, and transaction history. Go/Fiber structure matches
`backend/auth-service`.

Borrowing atomically reserves one copy through Book Service. A member may hold
at most three unfinished loans. Loans are due after seven days; each started
24-hour overdue period creates an IDR 5,000 charge on return.

Returning commits a `LoanReturned.v1` transactional-outbox event. Book Service
updates inventory through RabbitMQ and responds with `BookStockUpdated.v1`.
The return endpoint waits five seconds, then returns `202` when acknowledgement
remains pending. Retrying is safe.

Members can `GET /api/v1/transactions/loans/{loanId}/return` for an authoritative
fine quote. New clients send that quote as `acceptedFineAmountMinor` to the
matching `POST`. If the amount changed, the service returns `LMS-409006` without
changing the loan; bodyless legacy and librarian returns remain supported.

```sh
cp .env.example .env
go run . migrate --action up
go run . serve
```

Public API uses Kong prefix `/api/v1/transactions`. For direct local
development, the service listens on `http://127.0.0.1:8084`.

## API documentation

Interactive Swagger UI through Kong:

```text
http://127.0.0.1:8000/api/v1/docs/transactions/swagger
```

Generated Swagger 2.0 contract:

```text
http://127.0.0.1:8000/api/v1/docs/transactions/swagger.json
```

Handler annotations are source of truth. After endpoint changes, regenerate
and commit `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml`:

```sh
make swagger
# Equivalent Cobra command:
go run . swagger
```

```sh
make test
make test-race
make precommit
```
