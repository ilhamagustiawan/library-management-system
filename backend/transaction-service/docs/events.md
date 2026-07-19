# Transaction events

Durable topic exchange: `library.events`.

## `LoanReturned.v1`

Routing key: `transactions.loan.returned.v1`. Book Service consumes it from
`book-service.loan-returned.v1`. `eventId` is the idempotency key.

```json
{
  "eventId": "dd8cd46e-41e1-4583-aab2-c0342884201e",
  "type": "LoanReturned.v1",
  "occurredAt": "2026-07-19T10:00:00Z",
  "data": {
    "loanId": "52a88672-a4c2-4876-be5a-65863aeb35e4",
    "bookId": "7b36fe43-f31d-4861-884f-42ed7386b1e9",
    "memberId": "31c73b2e-0640-49bd-8f06-3bb7272921fe",
    "returnedAt": "2026-07-19T10:00:00Z"
  }
}
```

Book Service must release the reservation and atomically record an outgoing
`BookStockUpdated.v1` event. Duplicate `LoanReturned.v1` messages must return
the same acknowledgement without incrementing stock twice.

## `BookStockUpdated.v1`

Routing key: `books.stock.updated.v1`. `causationId` must equal the triggering
`LoanReturned.v1.eventId`.

```json
{
  "eventId": "3618f9db-b1ba-4e05-a88f-f0f5f92a1755",
  "type": "BookStockUpdated.v1",
  "occurredAt": "2026-07-19T10:00:01Z",
  "causationId": "dd8cd46e-41e1-4583-aab2-c0342884201e",
  "data": {
    "loanId": "52a88672-a4c2-4876-be5a-65863aeb35e4",
    "bookId": "7b36fe43-f31d-4861-884f-42ed7386b1e9",
    "updatedAt": "2026-07-19T10:00:01Z"
  }
}
```
