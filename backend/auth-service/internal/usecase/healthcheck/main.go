package healthcheck

import "context"

type Usecase interface {
	Readiness(ctx context.Context) error
}
