package loan

import (
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
)

func TestQuoteFineChargesStartedDays(t *testing.T) {
	due := time.Date(2026, 7, 26, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		returned time.Time
		want     int
	}{
		{name: "on time", returned: due, want: 0},
		{name: "one second late", returned: due.Add(time.Second), want: 1},
		{name: "exact day", returned: due.Add(24 * time.Hour), want: 1},
		{name: "started second day", returned: due.Add(24*time.Hour + time.Second), want: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			quote := entity.QuoteFine(due, test.returned, 5000)
			got := 0
			if quote != nil {
				got = quote.OverdueDays
			}
			if got != test.want {
				t.Fatalf("QuoteFine() days = %d, want %d", got, test.want)
			}
		})
	}
}
