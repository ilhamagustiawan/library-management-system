# User Service

Owns public member registration and member profile persistence. Auth Service remains credential and RBAC owner.

## Registration

`POST /api/v1/users` accepts `name`, `email`, and `password`. Unknown fields, including `role`, fail validation. User Service records a server-generated operation ID, creates the member identity through Auth Service, then atomically stores the profile and `UserRegistered.v1` outbox event.

Completed duplicate emails return `409`. An uncertain Auth call leaves the operation pending; retrying the same name and email resumes with the same internal idempotency key. RabbitMQ failure does not roll back registration.

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

With `go run . serve`, Swagger UI uses
`http://127.0.0.1:8082/api/v1/docs/users/swagger`. The default Compose stack
does not publish that service port; use
`http://127.0.0.1:8000/api/v1/docs/users/swagger` through Kong.

RabbitMQ publishes persistent `UserRegistered.v1` messages to `library.events` with routing key `user.registered.v1`. Delivery is at least once; consumers must deduplicate by `eventId`.
